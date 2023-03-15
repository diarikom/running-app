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
)

func TestInitiativeTestSuite(t *testing.T) {
	suite.Run(t, new(InitiativeTestSuite))
}

type InitiativeTestSuite struct {
	suite.Suite
	App apitest.Api
}

func (s *InitiativeTestSuite) SetupTest() {
	// Init app
	s.App = apitest.InitApi()

	// Setup data
	err := s.SetupData()
	if err != nil {
		panic(fmt.Errorf("failed to set-up data"))
	}

	// init service
	s.InitService()
}

func (s *InitiativeTestSuite) TearDownTest() {
	// Drop donation
	s.App.IgnoreDbExec(`DELETE FROM donation_log`)
	s.App.IgnoreDbExec(`DELETE FROM donation`)
	// Drop initiatives
	s.App.IgnoreDbExec(`DELETE FROM initiative`)
	s.App.IgnoreDbExec(`DELETE FROM organization_member`)
	s.App.IgnoreDbExec(`DELETE FROM organization`)
	// Drop credit
	s.App.IgnoreDbExec(`DELETE FROM user_credit_wallet_trx`)
	s.App.IgnoreDbExec(`DELETE FROM user_credit_wallet`)
	// Drop users
	s.App.IgnoreDbExec(`DELETE FROM user_profile`)
}

func (s *InitiativeTestSuite) SetupData() error {
	// Get instances
	db := s.App.Datasources.Db.Conn
	logger := s.App.Logger

	// Begin Transaction
	tx := db.MustBegin()
	var err error
	defer nsql.ReleaseTx(tx, &err, logger)

	// Insert users data
	_, err = db.Exec(`INSERT INTO user_profile (id, full_name, avatar_file, gender_id, date_of_birth, email, created_at, updated_at, email_verified) VALUES (1267772569398808560, 'John Doe', null, 1, '1999-12-31', 'johndoe@email.com', '2020-06-02 17:58:29.277934', '2020-06-02 17:58:29.277934', false), (1267772569398808561, 'Jane Doe', null, 1, '1999-12-31', 'janedoe@email.com', '2020-06-02 17:58:29.277934', '2020-06-02 17:58:29.277934', false);`)
	if err != nil {
		logger.Error("failed to insert users", err)
		return err
	}
	logger.Debug("user_profile inserted")

	// Insert runner credit
	_, err = db.Exec(` INSERT INTO public.user_credit_wallet (id, user_id, balance, balance_pending, balance_expiring, balance_expiring_date, created_at, updated_at, version) VALUES (1263038349258526720, 1267772569398808561, 4.00, 0.00, 0.00, null, '2020-05-20 16:26:23.238245', '2020-05-20 12:20:41.544957', 9); INSERT INTO public.user_credit_wallet_trx (id, user_credit_wallet_id, balance, balance_pending, amount, trx_entry_type_id, trx_ref_id, notes, status, created_at, expired_at, version) VALUES (1263038349258526721, 1263038349258526720, 0.00, 0.00, 0.00, 1, null, 'Init wallet', 2, '2020-05-20 16:26:23.238245', null, 1), (1263038348725850112, 1263038349258526720, 0.00, 1.00, 1.00, 2, null, null, 1, '2020-05-20 16:26:23.353251', '2020-05-27 09:26:23.353251', 2), (1263038803480678400, 1263038349258526720, 0.00, 2.00, 1.00, 2, null, null, 1, '2020-05-20 16:28:11.717725', '2020-05-27 09:28:11.717725', 3), (1263038875563986944, 1263038349258526720, 0.00, 3.00, 1.00, 2, null, null, 1, '2020-05-20 16:28:28.888209', '2020-05-27 09:28:28.888209', 4), (1263039035165642752, 1263038349258526720, 0.00, 4.00, 1.00, 2, null, null, 1, '2020-05-20 16:29:06.933349', '2020-05-27 09:29:06.933349', 5), (1263055783482888192, 1263038349258526720, 1.00, 3.00, 1.00, 2, 1263038803480678400, 'Claimed Credit from Challenge 1263038803845582848', 2, '2020-05-20 10:35:39.877308', '2020-06-19 10:35:39.877308', 6), (1263057018558615552, 1263038349258526720, 2.00, 2.00, 1.00, 2, 1263038348725850112, 'Claimed Credit from Challenge 1263038349011062784', 2, '2020-05-20 10:40:34.341756', '2020-06-19 10:40:34.341756', 7), (1263064084123750400, 1263038349258526720, 3.00, 1.00, 1.00, 2, 1263038875563986944, 'Claimed Credit from Challenge 1263038875840811008', 2, '2020-05-20 11:08:38.903478', '2020-06-19 11:08:38.903478', 8), (1263082214589992960, 1263038349258526720, 4.00, 0.00, 1.00, 2, 1263039035165642752, 'Claimed Credit from Challenge 1263039035387940864', 2, '2020-05-20 12:20:41.544957', '2020-06-19 12:20:41.544957', 9);`)
	if err != nil {
		logger.Error("failed to insert user credit", err)
		return err
	}
	logger.Debug("user credit inserted")

	// Insert initiative
	_, err = db.Exec(`INSERT INTO public.organization(id, name, org_type_id, description, logo_file, status_id, created_at, updated_at) VALUES (6675921466795622400, 'Sample Initiative', 3, NULL, NULL, 1, '2020-01-01 00:00:00', '2020-01-01 00:00:00'); INSERT INTO public.organization_member(id, organization_id, user_profile_id, member_role_id, created_at, updated_at) VALUES (6675925228188729344, 6675921466795622400, 1267772569398808560, 1, '2020-01-01 00:00:00', '2020-01-01 00:00:00'); INSERT INTO public.initiative(id, organization_id, name, description, image_files, external_urls, price, currency_id, status_id, tags, stat_donation_count, created_at, updated_at, version) VALUES (6675927595583930367, 6675921466795622400, 'Active Initiative Sample', 'Active Initiative Sample', '{ "thumbnail": "thumb_initiative.jpg", "detail_page": "detail_page_initiative.jpg" }', '[ { "text": "Sample Content", "url": "https://nbs.co.id" } ]', 1, 1, 2, 'sample', 0, '2020-01-01 00:00:00', '2020-01-01 00:00:00', 1), (6675927595583930368, 6675921466795622400, 'Draft Initiative Sample', 'Draft Initiative Sample', '{ "thumbnail": "thumb_initiative.jpg", "detail_page": "detail_page_initiative.jpg" }', '[ { "text": "Sample Content", "url": "https://nbs.co.id" } ]', 1, 1, 1, 'sample,draft', 0, '2020-01-01 00:00:00', '2020-01-01 00:00:00', 1), (6675927595583930369, 6675921466795622400, 'Inactive Initiative Sample', 'Draft Initiative Sample', '{ "thumbnail": "thumb_initiative.jpg", "detail_page": "detail_page_initiative.jpg" }', '[ { "text": "Sample Content", "url": "https://nbs.co.id" } ]', 1, 1, 3, 'sample,inactive', 0, '2020-01-01 00:00:00', '2020-01-01 00:00:00', 1);`)
	if err != nil {
		logger.Error("failed to insert initiative", err)
		return err
	}
	logger.Debug("initiative inserted")

	err = nil
	return nil
}

func (s *InitiativeTestSuite) InitService() {
	apiTest := s.App

	// Init services
	apiTest.Services = &api.Services{
		Asset:            &mocks.AssetService{},
		Auth:             &mocks.AuthenticatorService{},
		User:             &service.User{},
		Run:              &mocks.RunService{},
		DiscoverContent:  &mocks.DiscoverContentService{},
		Tag:              &mocks.AdTagService{},
		MilestoneService: &service.MilestoneService{},
		Credit:           &service.CreditService{},
		Initiative:       &service.Initiative{},
	}

	// Init services
	apiTest.MustInitService("UserService", apiTest.Services.User)
	apiTest.MustInitService("CreditService", apiTest.Services.Credit)
	apiTest.MustInitService("InitiativeService", apiTest.Services.Initiative)
}

func (s *InitiativeTestSuite) TestDonateSuccess() {
	svc := s.App.Services.Initiative
	logger := s.App.Logger

	resp, err := svc.Donate(dto.DonateReq{
		UserId:       "1267772569398808561",
		InitiativeId: "6675927595583930367",
		Quantity:     2,
	})
	if err != nil {
		logger.Error("error while donate", err)
		s.T().Fail()
		return
	}

	if resp.Balance != 2 {
		s.T().Errorf("expected remaining balance is %d, got %f", 2, resp.Balance)
		return
	}
}
