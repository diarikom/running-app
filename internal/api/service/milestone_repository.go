package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/model"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
	"time"
)

func NewMilestoneRepository(db *nsql.SqlDatabase, logger nlog.Logger) api.MilestoneRepository {
	r := MilestoneRepository{
		Db:      db,
		Differs: initMilestoneDiffer(),
		Stmt:    initMilestoneStatement(db),
		Logger:  logger,
	}

	return &r
}

type MilestoneRepository struct {
	Db      *nsql.SqlDatabase
	Differs milestoneDiffer
	Stmt    MileStatement
	Logger  nlog.Logger
}

func (r *MilestoneRepository) FindUnaccomplishedChallengeByUser(userID string, milestoneID string) ([]model.Challenge, error) {
	// Create query
	rows := make([]model.Challenge, 0)
	err := r.Stmt.findUnaccomplishedChallengeByUser.Select(&rows, userID, milestoneID, api.ChallengeStart)
	return rows, err
}

func (r MilestoneRepository) FindById(id string) (*model.Milestone, error) {
	var m model.Milestone
	err := r.Stmt.findById.Get(&m, id)

	return &m, err
}

func (r MilestoneRepository) InsertUserChallenge(uc model.UserChallenge) error {
	_, err := r.Stmt.insertUserChallenge.Exec(uc)
	return err
}

func (r MilestoneRepository) FindChallengeById(id string) (*model.Challenge, error) {
	var c model.Challenge
	err := r.Stmt.findChallengeById.Get(&c, id)
	return &c, err
}

func (r MilestoneRepository) FindCurrentChallenges() ([]model.Challenge, error) {
	var c []model.Challenge
	err := r.Stmt.findCurrentChallenges.Select(&c, api.ChallengeStart, api.MilestoneStart)
	return c, err
}

func (r MilestoneRepository) Current(now time.Time, status int) (result model.Milestone, err error) {
	err = r.Stmt.current.Get(&result, now, status)

	return result, err
}

func (r *MilestoneRepository) FindUserChallenge(userID string, challengeID string) (*model.UserChallenge, error) {
	var uc model.UserChallenge
	err := r.Stmt.findUserChallenge.Get(&uc, userID, challengeID)
	return &uc, err
}

func (r *MilestoneRepository) UpdateUserChallenge(o, n model.UserChallenge, changes []string) error {
	// Get differ
	differ := r.Differs.userChallenge

	// Compare instance
	diff, err := differ.Compare(o, n, changes)
	if err != nil {
		return err
	}

	// If no changes, return
	if diff.Count == 0 {
		r.Logger.Debug("no user challenge changes detected")
		return nil
	}

	// Generate query
	q, args, err := differ.UpdateQuery(diff)
	if err != nil {
		r.Logger.Error("unable to generate update query", err)
		return err
	}

	// Rebind query
	q = r.Db.Conn.Rebind(q)

	// Execute
	_, err = r.Db.Conn.Exec(q, args...)
	return err
}

func (r MilestoneRepository) GetChallengesByStatus(userID string, milestoneID string, status int) ([]model.Challenge, error) {
	var result []model.Challenge
	err := r.Stmt.getChallengesByStatus.Select(&result, userID, milestoneID, status)

	return result, err
}

func (r MilestoneRepository) GetChallengesIn(challengesID string) ([]model.Challenge, error) {
	var result []model.Challenge
	err := r.Stmt.getChallengesIn.Select(&result, challengesID)

	return result, err
}

func (r *MilestoneRepository) GetUserChallengeByMilestoneId(userID string, milestoneID string, status int8) ([]model.UserChallenge, error) {
	var uc []model.UserChallenge
	err := r.Stmt.getUserChallengeByMilestoneId.Select(&uc, userID, milestoneID, status)

	return uc, err
}
