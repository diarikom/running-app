package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/dto"
	"github.com/diarikom/running-app/running-app-api/internal/api/model"
	"github.com/diarikom/running-app/running-app-api/internal/pkg/ngrule"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"strconv"
	"time"
)

const TopicCheckAchievedChallenge = "check_achieved_challenge"

type MilestoneService struct {
	IdGen               *api.SnowflakeGen
	Error               *api.Errors
	Logger              nlog.Logger
	MilestoneRepository api.MilestoneRepository
	RunRepository       api.RunRepository
	CreditService       api.CreditService
	RuleEngineMemory    *ast.WorkingMemory
	RuleEngine          *engine.GruleEngine
	RuleMap             map[string]*model.ChallengeRule
	FactFinder          *ngrule.FactFinderMap
	PubSub              *gochannel.GoChannel
	UserService         api.UserService
}

func (s *MilestoneService) Init(app *api.Api) error {
	mRepo := NewMilestoneRepository(app.Datasources.Db, app.Logger)
	rRepo := NewRunRepository(app.Datasources.Db, app.Logger)

	s.IdGen = app.Components.Id
	s.Error = app.Components.Errors
	s.Logger = app.Logger
	s.MilestoneRepository = mRepo
	s.RunRepository = rRepo
	s.CreditService = app.Services.Credit
	s.RuleEngineMemory = ast.NewWorkingMemory()
	s.RuleEngine = engine.NewGruleEngine()
	s.UserService = app.Services.User

	// Init current active challenge rules
	s.LoadMilestone()

	// Register Fact Finder function
	s.initFactFinder()

	// Init PubSub
	s.PubSub = gochannel.NewGoChannel(gochannel.Config{}, nil)

	// Subscribe to check challenge achieved
	msg, err := s.PubSub.Subscribe(context.Background(), TopicCheckAchievedChallenge)
	if err != nil {
		return err
	}
	go s.handleCheckChallengeAchieved(msg)

	return nil
}

func (s *MilestoneService) TriggerCheckChallengeAchieved(req dto.UserChallengeReq) error {
	// Encode to gob
	var w bytes.Buffer
	enc := gob.NewEncoder(&w)
	err := enc.Encode(req)
	if err != nil {
		s.Logger.Error("failed to encode dto.UserChallengeReq payload", err)
		return err
	}

	// Create message
	msg := message.NewMessage(watermill.NewUUID(), w.Bytes())

	// Publish
	err = s.PubSub.Publish(TopicCheckAchievedChallenge, msg)
	if err != nil {
		s.Logger.Error("failed to publish to "+TopicCheckAchievedChallenge, err)
		return err
	}

	return nil
}

func (s *MilestoneService) handleCheckChallengeAchieved(messages <-chan *message.Message) {
	for msg := range messages {
		s.Logger.Debugf("Received message. Id = %s", msg.UUID)

		// Parse payload
		w := bytes.NewBuffer(msg.Payload)
		dec := gob.NewDecoder(w)
		var payload dto.UserChallengeReq
		err := dec.Decode(&payload)
		if err != nil {
			s.Logger.Error("failed to parse payload", err)
			msg.Ack()
			return
		}

		// Check challenge achieved
		_, err = s.CheckChallengeAchieve(payload)
		if err != nil {
			s.Logger.Error("failed to check achieved challenge", err)
			msg.Ack()
			return
		}

		s.Logger.Debug("Done handleCheckChallengeAchieved")
		msg.Ack()
	}
}

func (s MilestoneService) Current(userID string) (resp *dto.MilestoneResp, err error) {

	now := time.Now().UTC()
	milestone, err := s.MilestoneRepository.Current(now, api.MilestoneStart)

	if err != nil {

		if err == sql.ErrNoRows {
			resp = &dto.MilestoneResp{
				Milestone: nil,
			}
			return resp, nil
		}

		s.Logger.Error("unable to get milestone", err)
		return resp, err
	}

	items, err := s.GetMilestoneChallenges(userID, milestone.Id)
	if err != nil {
		s.Logger.Error("unable's to get challenge", err)
		return resp, err
	}

	m := dto.CurrentMilestoneResp{
		Id:             milestone.Id,
		PeriodStart:    milestone.PeriodStart.Unix(),
		PeriodEnd:      milestone.PeriodEnd.Unix(),
		PeriodTzOffset: milestone.PeriodTZ,
		Status:         milestone.Status,
		Challenges:     items,
	}

	resp = &dto.MilestoneResp{
		Milestone: &m,
	}

	return resp, nil
}

func (s MilestoneService) GetMilestoneChallenges(userID string, milestoneID string) (resp []dto.MilestoneChallengesResp, err error) {

	challenges, err := s.MilestoneRepository.GetChallengesByStatus(userID, milestoneID, api.MilestoneStart)
	if err != nil {
		s.Logger.Error("unable to get challenges", err)
		return resp, err
	}

	resp = make([]dto.MilestoneChallengesResp, len(challenges))
	for k, v := range challenges {
		// Get rewards value
		value := 0
		target := v.Rules.Actions.FindByTargetName("credit")
		if target != nil {
			value = ngrule.ParseInt(target.GetValue(), 0)
		}

		resp[k] = dto.MilestoneChallengesResp{
			Id:   v.Id,
			Name: v.Title,
			Reward: dto.ChallengeRewardResp{
				Type:     "Credit",
				Value:    value,
				Currency: "dB",
			},
			Status:    v.Status,
			UpdatedAt: v.UpdatedAt.Unix(),
		}
	}

	return resp, err
}

func (s MilestoneService) CheckChallengeAchieve(req dto.UserChallengeReq) (*dto.MilestoneAchievementResp, error) {
	// Get active milestone
	m, err := s.MilestoneRepository.Current(req.Timestamp, api.MilestoneStart)
	if err != nil {
		s.Logger.Error("unable to get current active milestone", err)
		return nil, err
	}

	// Get unaccomplished challenges
	challenges, err := s.MilestoneRepository.FindUnaccomplishedChallengeByUser(req.UserId, m.Id)
	if err != nil {
		s.Logger.Error("unable to get list user challenge", err)
		return nil, err
	}

	// If there's user challenge
	if len(challenges) == 0 {
		return nil, s.Error.New("CLG001")
	}

	// Get premium status
	premium, err := s.UserService.IsPremiumRunner(req.UserId)
	if err != nil {
		return nil, err
	}

	// Determine multiplier by premium status
	var rewardMultiplier int64 = 1
	if premium {
		rewardMultiplier = 2
	}

	// Get required facts keys
	factParams := s.getRequiredFactParams(challenges)

	// Init fact
	facts := ngrule.NewFactMap(factParams)

	// find facts
	err = s.FactFinder.FindFacts(&facts, factParams, req.UserId)
	if err != nil {
		return nil, err
	}

	// Init data context
	data := ast.NewDataContext()
	err = data.Add("Var", &facts)
	if err != nil {
		s.Logger.Error("failed to add facts", err)
	}

	// Init reward
	var rewardCredit int64

	for _, v := range challenges {
		// Get loaded challenges
		challengeId := v.Id
		rule, ok := s.RuleMap[challengeId]

		// If rule is not loaded, load rule
		if !ok {
			s.Logger.Debugf("Rule not loaded. ChallengeId: %s", challengeId)
			rule, err = s.PrepareChallengeRule(v, "Var")
			if err != nil {
				s.Logger.Error("failed to load rule", err)
				s.Logger.Errorf("ChallengeId: %s", challengeId)
				return nil, err
			}

			s.RuleMap[challengeId] = rule
		}

		// Set target
		for _, vt := range rule.Targets {
			facts.Assign(vt.GetName(), vt)
		}

		err := s.RuleEngine.Execute(data, rule.RuleBuilder.KnowledgeBase, s.RuleEngineMemory)
		if err != nil {
			s.Logger.Error("failed to execute rule", err)
		}

		// Get credits
		reward := facts.GetInt("credit")

		// If no reward, then break
		if reward == 0 {
			s.Logger.Debug("no reward. ending challenge check")
			break
		}

		// Multiply reward
		reward *= rewardMultiplier

		// Add to reward credit
		err = s.UpdateUserChallengeAchievement(dto.UserChallengeReq{
			UserId:      req.UserId,
			ChallengeId: challengeId,
			RewardValue: float64(reward),
			ChallengeResultSnapshot: map[string]interface{}{
				api.UserAccumulatedRunDistanceFact: facts.GetInt(api.UserAccumulatedRunDistanceFact),
			},
		})
		if err != nil {
			s.Logger.Error("failed to update user challenge achievement", err)
			return nil, err
		}

		rewardCredit += reward
	}

	s.Logger.Debugf("RewardCredit = %d", rewardCredit)

	return &dto.MilestoneAchievementResp{
		CreditReward: rewardCredit,
	}, nil
}

func (s *MilestoneService) getRequiredFactParams(challenges []model.Challenge) []ngrule.FactParam {
	// Init parameters
	factParams := make([]ngrule.FactParam, 0)

	for _, v := range challenges {
		// Get parameters
		p := v.Rules.GetParams()

		// Merge parameters from challenges
		factParams = ngrule.MergeParams(factParams, p)
	}

	return factParams
}

func (s *MilestoneService) ClaimCredit(opt dto.UserChallengeReq) (*dto.ChallengeRewardClaimResp, error) {
	// Validate options
	if opt.ChallengeId == "" {
		return nil, errors.New("ChallengeId is required")
	}

	if opt.UserId == "" {
		return nil, errors.New("UserId is required")
	}

	// Get available user challenge based on user id, challenge id
	userChallenge, err := s.MilestoneRepository.FindUserChallenge(opt.UserId, opt.ChallengeId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, s.Error.New("CLG001")
		}
		s.Logger.Error("unable to find user challenge", err)
		return nil, err
	}

	// Check status
	switch userChallenge.Status {
	case api.ChallengeAchieved:
		break
	case api.ChallengeRewardClaimed:
		return nil, s.Error.New("CLG003")
	default:
		return nil, s.Error.New("CLG002")
	}

	// Create timestamp
	timestamp := time.Now()

	// Update user challenge
	newUserChallenge := model.UserChallenge{
		Id:        userChallenge.Id,
		Status:    api.ChallengeRewardClaimed,
		UpdatedAt: timestamp,
	}

	// Set expired at to 1 month
	// TODO: Get expired at from config
	expiredAt := timestamp.Add(time.Duration(720) * time.Hour)

	// Settle pending transaction
	err = s.CreditService.SettlePendingTrx(dto.CreditSettleOpt{
		TrxId:     userChallenge.RewardRefId,
		Notes:     "Claimed Credit from Challenge " + userChallenge.Id,
		ExpiredAt: &expiredAt,
		Timestamp: &timestamp,
	})
	if err != nil {
		return nil, err
	}

	// Update user challenge status
	err = s.MilestoneRepository.UpdateUserChallenge(*userChallenge, newUserChallenge, []string{"status"})
	if err != nil {
		s.Logger.Error("unable to persist user challenge update", err)
		return nil, err
	}

	// Compose response
	resp := &dto.ChallengeRewardClaimResp{
		Status: newUserChallenge.Status,
	}

	return resp, nil
}

func (s *MilestoneService) AddUserChallenge(opt dto.UserChallengeReq) error {
	// Validate options
	if opt.ChallengeId == "" {
		return errors.New("ChallengeId is required")
	}

	if opt.UserId == "" {
		return errors.New("UserId is required")
	}

	// Get challenge
	c, err := s.MilestoneRepository.FindChallengeById(opt.ChallengeId)
	if err != nil {
		s.Logger.Error("unable to find challenge by id", err)
		return err
	}

	// Get milestones
	m, err := s.MilestoneRepository.FindById(c.MilestoneId)
	if err != nil {
		s.Logger.Error("unable to find challenge by id", err)
		return err
	}

	// Create snapshot
	// TODO: Implement scanner and valuer interface
	ms, err := json.Marshal(m)
	if err != nil {
		s.Logger.Error("unable to capture milestone snapshot", err)
		return err
	}
	s.Logger.Debugf("MilestoneSnapshot: %s", ms)

	cs, err := json.Marshal(c)
	if err != nil {
		s.Logger.Error("unable to capture challenge snapshot", err)
		return err
	}
	s.Logger.Debugf("ChallengeSnapshot: %s", cs)

	// Create challenge result snapshot
	var crs []byte
	if opt.ChallengeResultSnapshot != nil {
		crs, err = json.Marshal(opt.ChallengeResultSnapshot)
		if err != nil {
			s.Logger.Error("unable to capture challenge result snapshot", err)
			return err
		}
	} else {
		crs = []byte("{}")
	}

	// Create user challenge
	uc := model.UserChallenge{
		Id:                      s.IdGen.New(),
		UserId:                  opt.UserId,
		MilestoneId:             m.Id,
		MilestoneSnapshot:       ms,
		MilestoneVersion:        m.Version,
		ChallengeId:             c.Id,
		ChallengeSnapshot:       cs,
		ChallengeVersion:        c.Version,
		ChallengeResultSnapshot: crs,
		RewardSnapshot:          []byte("{}"),
		RewardTypeId:            api.CreditReward,
		RewardRefId:             opt.RewardRefId,
		RewardValue:             opt.RewardValue,
		Status:                  api.ChallengeAchieved,
		UpdatedAt:               opt.Timestamp,
	}

	// Insert user challenge
	err = s.MilestoneRepository.InsertUserChallenge(uc)
	if err != nil {
		s.Logger.Error("unable to insert user challenge", err)
		return err
	}

	return nil
}

func (s *MilestoneService) UpdateUserChallengeAchievement(opt dto.UserChallengeReq) error {
	// Generate trx reference id
	opt.RewardRefId = s.IdGen.New()
	opt.Timestamp = time.Now()

	// Insert user challenge
	err := s.AddUserChallenge(opt)
	if err != nil {
		return err
	}

	// Get wallet
	wallet, err := s.CreditService.GetUserWallet(opt.UserId)
	if err != nil {
		return err
	}

	// Init timestamp
	timestamp := time.Now()

	// TODO: Set expire to from config
	// Set expire to 7 days
	claimExpire := timestamp.Add(time.Duration(168) * time.Hour)

	_, err = s.CreditService.InsertPendingTrx(dto.CreditTrxOpt{
		Id:        opt.RewardRefId,
		WalletId:  wallet.Id,
		Amount:    opt.RewardValue,
		EntryType: api.Debit,
		ExpiredAt: &claimExpire,
		Timestamp: &timestamp,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *MilestoneService) LoadMilestone() {
	// Get current challenges
	challenges, err := s.MilestoneRepository.FindCurrentChallenges()
	if err != nil {
		panic(fmt.Errorf("running-app-api: error while retrieving challenge rule (%w)", err))
	}

	// Prepare challenge rule
	challengeRules := make(map[string]*model.ChallengeRule)
	for _, v := range challenges {
		cr, err := s.PrepareChallengeRule(v, "Var")
		if err != nil {
			panic(fmt.Errorf("running-app-api: error while preparing challenge rule (%w)", err))
		}

		// Set challenge rules
		challengeRules[v.Id] = cr
	}

	// Set challenge rules
	s.RuleMap = challengeRules
}

func (s *MilestoneService) PrepareChallengeRule(challenge model.Challenge, varName string) (*model.ChallengeRule, error) {
	// Get rule
	r := challenge.Rules

	// if code is not set, get from user id
	if r.Code == "" {
		r.Code = "Rule" + challenge.Id
	}

	// Set variable name
	r.VariableName = varName

	// Render grule syntax
	ruleSyntax, err := r.Render()
	if err != nil {
		s.Logger.Error("unable to render rule", err)
		return nil, err
	}

	// Create knowledge base
	kb := ast.NewKnowledgeBase(r.Code, strconv.Itoa(challenge.Version))

	// Create rule builder
	rb := builder.NewRuleBuilder(kb, s.RuleEngineMemory)

	// Parse rule syntax
	byteArr := pkg.NewBytesResource([]byte(ruleSyntax))
	err = rb.BuildRuleFromResource(byteArr)
	if err != nil {
		s.Logger.Error("unable to build grule syntax", err)
		return nil, err
	}

	cr := model.ChallengeRule{
		RuleBuilder: rb,
		Params:      r.GetParams(),
		Targets:     r.GetTargets(),
	}

	return &cr, nil
}

func (s *MilestoneService) initFactFinder() {
	factFinder := ngrule.NewFactFinderMap()

	factFinder.RegisterParamFn(api.UserAccumulatedRunDistanceFact, s.CalcUserAccRunDistance)

	s.FactFinder = factFinder
}

func (s *MilestoneService) CalcUserAccRunDistance(m *ngrule.FactMap, userId string) error {
	now := time.Now().UTC()
	milestone, err := s.MilestoneRepository.Current(now, api.MilestoneStart)
	if err != nil {
		return err
	}

	total, err := s.RunRepository.SumRunSessionDistance(userId, milestone.PeriodStart, milestone.PeriodEnd)
	if err != nil {
		total = 0
	}

	m.Set(api.UserAccumulatedRunDistanceFact, total)
	return nil
}
