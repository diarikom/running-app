package service_test

import (
	"fmt"
	"github.com/diarikom/running-app/running-app-api/cmd/apitest"
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/dto"
	"github.com/diarikom/running-app/running-app-api/internal/api/mocks"
	"github.com/diarikom/running-app/running-app-api/internal/api/service"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestMilestoneTestSuite(t *testing.T) {
	suite.Run(t, new(MilestoneTestSuite))
}

type MilestoneTestSuite struct {
	suite.Suite
	App apitest.Api
}

func (s *MilestoneTestSuite) SetupTest() {
	// Init app
	s.App = apitest.InitApi()

	// setup data
	err := s.SetupData()
	if err != nil {
		panic(fmt.Errorf("failed to set-up data"))
	}

	// init service
	s.InitService()
}

func (s *MilestoneTestSuite) TearDownTest() {
	// Drop credit
	s.App.IgnoreDbExec(`DELETE FROM user_credit_wallet_trx`)
	s.App.IgnoreDbExec(`DELETE FROM user_credit_wallet`)
	// Drop user challenges
	s.App.IgnoreDbExec(`DELETE FROM user_challenge`)
	// Drop challenges
	s.App.IgnoreDbExec(`DELETE FROM challenge`)
	// Drop milestone
	s.App.IgnoreDbExec(`DELETE FROM milestone`)
	// Drop run session
	s.App.IgnoreDbExec(`DELETE FROM run_session_log`)
	s.App.IgnoreDbExec(`DELETE FROM run_session`)
	// Drop users
	s.App.IgnoreDbExec(`DELETE FROM user_profile`)
}

func (s *MilestoneTestSuite) SetupData() error {
	// Get instances
	db := s.App.Datasources.Db.Conn
	logger := s.App.Logger

	// Begin Transaction
	tx := db.MustBegin()
	var err error
	defer nsql.ReleaseTx(tx, &err, logger)

	// Insert users data
	_, err = db.Exec(`INSERT INTO user_profile (id, full_name, avatar_file, gender_id, date_of_birth, email, created_at, updated_at, email_verified) VALUES (1267772569398808576, 'John Doe', null, 1, '1999-12-31', 'johndoe@email.com', '2020-06-02 17:58:29.277934', '2020-06-02 17:58:29.277934', false);`)
	if err != nil {
		logger.Error("failed to insert user_profile", err)
		return err
	}
	logger.Debug("user_profile inserted")

	// Insert sample milestone
	_, err = db.Exec(`INSERT INTO milestone (id, name, period_start, period_end, period_tz, status, created_at, updated_at, version) VALUES (6667296413904404480, '2020 Milestone', '2020-01-01 00:00:00.000000', '2021-01-01 00:00:00.000000', 420, 2, '2020-05-01 00:00:00.000000', '2020-05-01 00:00:00.000000', 1);`)
	if err != nil {
		logger.Error("failed to insert milestone", err)
		return err
	}
	logger.Debug("milestone inserted")

	// Insert sample challenges
	_, err = db.Exec(`INSERT INTO challenge (id, milestone_id, title, description, level, status, rules, sort, created_at, updated_at, version) VALUES (6667296413828907008, 6667296413904404480, '10K', 'Run in 10 kilometers', 1, 2, '{ "description": "User reach 10km of accumulated run distance during milestone will get 1 credit", "priority": 0, "params": [], "conditions": [ { "type": "int", "options": { "param": "user_acc_run_distance", "operator": ">=", "ref_value": 10000 } } ], "actions": [ { "type": "add", "options": { "target": "credit", "value_type": "int", "value": 1 } } ] }', 0, '2020-05-01 00:00:00.000000', '2020-05-01 00:00:00.000000', 1), (6667296413854072833, 6667296413904404480, '50K', 'Run in 50 kilometers', 5, 2, '{ "description": "User reach 50km of accumulated run distance during milestone will get 1 credit", "priority": 0, "params": [], "conditions": [ { "type": "int", "options": { "param": "user_acc_run_distance", "operator": ">=", "ref_value": 50000 } } ], "actions": [ { "type": "add", "options": { "target": "credit", "value_type": "int", "value": 1 } } ] }', 0, '2020-05-01 00:00:00.000000', '2020-05-01 00:00:00.000000', 1), (6667296413849878528, 6667296413904404480, '20K', 'Run in 20 kilometers', 2, 2, '{ "description": "User reach 20km of accumulated run distance during milestone will get 1 credit", "priority": 0, "params": [], "conditions": [ { "type": "int", "options": { "param": "user_acc_run_distance", "operator": ">=", "ref_value": 20000 } } ], "actions": [ { "type": "add", "options": { "target": "credit", "value_type": "int", "value": 1 } } ] }', 0, '2020-05-01 00:00:00.000000', '2020-05-01 00:00:00.000000', 1), (6667296413849878529, 6667296413904404480, '30K', 'Run in 30 kilometers', 3, 2, '{ "description": "User reach 30km of accumulated run distance during milestone will get 1 credit", "priority": 0, "params": [], "conditions": [ { "type": "int", "options": { "param": "user_acc_run_distance", "operator": ">=", "ref_value": 30000 } } ], "actions": [ { "type": "add", "options": { "target": "credit", "value_type": "int", "value": 1 } } ] }', 0, '2020-05-01 00:00:00.000000', '2020-05-01 00:00:00.000000', 1), (6667296413854072832, 6667296413904404480, '40K', 'Run in 40 kilometers', 4, 2, '{ "description": "User reach 40km of accumulated run distance during milestone will get 1 credit", "priority": 0, "params": [], "conditions": [ { "type": "int", "options": { "param": "user_acc_run_distance", "operator": ">=", "ref_value": 40000 } } ], "actions": [ { "type": "add", "options": { "target": "credit", "value_type": "int", "value": 1 } } ] }', 0, '2020-05-01 00:00:00.000000', '2020-05-01 00:00:00.000000', 1);`)
	if err != nil {
		logger.Error("failed to insert challenge", err)
		return err
	}
	logger.Debug("challenge inserted")

	// Insert run sample data
	_, err = db.Exec(`INSERT INTO run_session (id, user_id, session_started, session_ended, time_elapsed, distance, speed, step_count, created_at, sync_status_id, updated_at, version) VALUES (1267003410025025536, 1267772569398808576, '2019-05-31 07:19:29.000000', '2019-05-31 08:02:04.000000', 2555, 9451, 2.86, 7321, '2019-05-31 08:02:07.000000', 1, '2019-05-31 08:02:07.000000', 1), (1267003410025025537, 1267772569398808576, '2020-01-31 07:19:29.000000', '2020-01-31 08:02:04.000000', 2555, 7289, 2.86, 7321, '2020-01-31 08:02:04.000000', 1, '2020-01-31 08:02:04.000000', 1), (1267003410025025538, 1267772569398808576, '2020-02-01 07:19:29.000000', '2020-02-01 08:02:04.000000', 2555, 6349, 2.86, 7321, '2020-02-01 08:02:04.000000', 1, '2020-02-01 08:02:04.000000', 1), (1267003410025025539, 1267772569398808576, '2020-05-31 07:19:29.000000', '2020-05-31 08:02:04.000000', 2555, 8542, 2.86, 7321, '2020-05-31 08:02:04.000000', 1, '2020-05-31 08:02:04.000000', 1);`)
	if err != nil {
		logger.Error("failed to insert run_session", err)
		return err
	}
	logger.Debug("run_session inserted")

	err = nil
	return nil
}

func (s *MilestoneTestSuite) InitService() {
	apiTest := s.App

	// Init services
	apiTest.Services = &api.Services{
		Asset:            &mocks.AssetService{},
		Auth:             &mocks.AuthenticatorService{},
		User:             &mocks.UserService{},
		Run:              &mocks.RunService{},
		DiscoverContent:  &mocks.DiscoverContentService{},
		Tag:              &mocks.AdTagService{},
		MilestoneService: &service.MilestoneService{},
		Credit:           &service.CreditService{},
	}

	// Init services
	apiTest.MustInitService("CreditService", apiTest.Services.Credit)
	apiTest.MustInitService("MilestoneService", apiTest.Services.MilestoneService)
}

func (s *MilestoneTestSuite) TestChallenge() {
	svc := s.App.Services.MilestoneService
	logger := s.App.Logger

	resp, err := svc.CheckChallengeAchieve(dto.UserChallengeReq{
		UserId:    "1267772569398808576",
		Timestamp: time.Unix(1591748100, 0),
	})
	if err != nil {
		logger.Error("error while check challenge achieve", err)
		s.T().FailNow()
	}

	logger.Debugf("Result = %+v", resp)
}
