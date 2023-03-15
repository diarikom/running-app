package service

import (
	"database/sql"
	"encoding/json"
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/dto"
	"github.com/diarikom/running-app/running-app-api/internal/api/model"
	"github.com/diarikom/running-app/running-app-api/pkg/nhttp"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
	validate "github.com/go-playground/validator/v10"
	"strings"
	"time"
)

type Initiative struct {
	IdGen         *api.SnowflakeGen
	Errors        *api.Errors
	Logger        nlog.Logger
	Repository    api.InitiativeRepository
	AssetService  api.AssetService
	CreditService api.CreditService
	UserService   api.UserService
	Validator     *validate.Validate
}

func (s *Initiative) Init(app *api.Api) error {
	// Init run service
	s.IdGen = app.Components.Id
	s.Errors = app.Components.Errors
	s.Logger = app.Logger
	s.Repository = NewInitiativeRepository(app.Datasources.Db, app.Components.Id, app.Components.Errors, app.Logger)
	s.AssetService = app.Services.Asset
	s.CreditService = app.Services.Credit
	s.UserService = app.Services.User
	s.Validator = validate.New()
	return nil
}

func (s *Initiative) ListUserDonation(opt dto.UserResourcesReq) ([]dto.DonationHistoryResp, error) {
	// Get donation history
	rows, err := s.Repository.FindDonationByUser(opt.UserId, opt.Skip, opt.Limit)
	if err != nil {
		return nil, err
	}

	// Compose response
	count := len(rows)
	resp := make([]dto.DonationHistoryResp, count)
	for k, v := range rows {
		// Get thumbnail url
		thumbnailUrl := s.AssetService.GetPublicUrl(api.AssetInitiative, v.InitiativeSnapshot.ImageFiles.Thumbnail)

		// Compose response item
		resp[k] = dto.DonationHistoryResp{
			Id:                     v.Id,
			InitiativeId:           v.InitiativeId,
			InitiativeName:         v.InitiativeSnapshot.Name,
			InitiativeDescription:  v.InitiativeSnapshot.Description,
			InitiativeThumbnailURL: thumbnailUrl,
			ItemPrice:              v.InitiativeSnapshot.Price,
			Quantity:               v.Quantity,
			TotalDonation:          v.TotalPrice,
			StatusId:               v.StatusId,
			CreatedAt:              v.CreatedAt.Unix(),
			UpdatedAt:              v.UpdatedAt.Unix(),
		}
	}

	return resp, nil
}

func (s *Initiative) Donate(opt dto.DonateReq) (*dto.DonateResp, error) {
	// Validate request
	err := s.Validator.Struct(&opt)
	if err != nil {
		s.Logger.Error("failed to validate", err)
		return nil, nhttp.ErrBadRequest
	}

	// Get initiative
	initiative, err := s.Repository.FindById(opt.InitiativeId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, s.Errors.New("INT003")
		}
		s.Logger.Error("failed to FindById initiative", err)
		return nil, err
	}

	// Check if initiative is active
	if initiative.StatusId != api.ActiveInitiative {
		return nil, s.Errors.New("INT001")
	}

	// Calculate total donation
	totalDonation := initiative.Price * float64(opt.Quantity)

	// Check if balance is enough
	walletVersion, err := s.CreditService.CheckChargeAmount(dto.CreditChargeOpt{
		UserId: opt.UserId,
		Amount: totalDonation,
	})
	if err != nil {
		return nil, err
	}

	// Insert donation
	donation, err := s.createDonation(initiative, totalDonation, opt)
	if err != nil {
		return nil, err
	}

	// Charge donation
	err = s.chargeDonation(donation, walletVersion)
	if err != nil {
		return nil, err
	}

	// Settle Pending Transaction
	err = s.CreditService.SettlePendingTrx(dto.CreditSettleOpt{
		TrxId: donation.PaymentTrxRef.String,
	})

	// Update donation status to success
	err = s.updateDonationSuccess(donation)
	if err != nil {
		return nil, err
	}

	balance, err := s.CreditService.GetUserBalance(opt.UserId)
	if err != nil {
		return nil, err
	}

	resp := dto.DonateResp{
		Balance:       balance.Balance,
		TotalDonation: totalDonation,
	}
	return &resp, nil
}

func (s *Initiative) createDonation(initiative *model.Initiative, totalDonation float64, opt dto.DonateReq) (
	*model.Donation, error) {
	// Get user snapshot
	userSnapshot, err := s.UserService.GetProfileSnapshot(opt.UserId)
	if err != nil {
		return nil, err
	}

	// Create initiative snapshot
	initiativeSnapshot := model.NewInitiativeSnapshot(initiative)

	// Init timestamp
	timestamp := time.Now()

	// Create modifier meta
	modifiedBy := model.ModifierMeta{
		Id:       userSnapshot.Id,
		Role:     api.ModifierUser,
		FullName: userSnapshot.FullName,
	}

	// Create Donation
	donation := model.Donation{
		Id:                 s.IdGen.New(),
		InitiativeId:       initiative.Id,
		InitiativeSnapshot: &initiativeSnapshot,
		UserId:             userSnapshot.Id,
		UserSnapshot:       userSnapshot,
		PaymentMethodId:    api.PaymentMethodCredit,
		PaymentSnapshot:    json.RawMessage("{}"),
		PaymentTrxRef:      sql.NullString{},
		Quantity:           opt.Quantity,
		TotalPrice:         totalDonation,
		CurrencyId:         api.CurrencyDecibel,
		StatusId:           api.DonationCreated,
		Notes:              nsql.NullString(opt.Notes),
		CreatedAt:          timestamp,
		UpdatedAt:          timestamp,
		ModifiedBy:         &modifiedBy,
		Version:            1,
	}

	donationLog := model.DonationLog{
		LogId:           s.IdGen.New(),
		Changelog:       []string{},
		Id:              donation.Id,
		PaymentMethodId: sql.NullInt32{Int32: int32(donation.PaymentMethodId), Valid: true},
		PaymentSnapshot: donation.PaymentSnapshot,
		PaymentTrxRef:   sql.NullString{},
		StatusId:        sql.NullInt32{Int32: int32(donation.StatusId), Valid: true},
		Notes:           donation.Notes,
		UpdatedAt:       timestamp,
		ModifiedBy:      &modifiedBy,
		Version:         1,
	}

	err = s.Repository.Insert(donation, donationLog)
	if err != nil {
		return nil, err
	}

	return &donation, nil
}

func (s *Initiative) chargeDonation(donation *model.Donation, walletVersion int) error {
	// Charge transaction
	trxId, err := s.CreditService.Charge(dto.CreditChargeOpt{
		UserId:        donation.UserId,
		Amount:        donation.TotalPrice,
		WalletVersion: walletVersion,
	})
	if err != nil {
		return err
	}

	// Duplicate donation
	oldDonation, err := model.CopyDonation(donation)
	if err != nil {
		s.Logger.Error("failed to duplicate donation for changelog", err)
		return err
	}

	// Update donation
	timestamp := time.Now()
	donation.StatusId = api.DonationPaymentPending
	donation.PaymentTrxRef = nsql.NullString(trxId)
	donation.UpdatedAt = timestamp
	donation.ModifiedBy = &model.ModifierMeta{
		Id:       "SYSTEM",
		Role:     "SYSTEM",
		FullName: "SYSTEM",
	}
	donation.Version = oldDonation.Version + 1

	// Create changelog
	changelog := []string{
		"status_id",
		"payment_trx_ref",
	}

	// Update donation
	err = s.Repository.UpdateDonation(*oldDonation, *donation, changelog)
	if err != nil {
		return err
	}

	return nil
}

func (s *Initiative) updateDonationSuccess(donation *model.Donation) error {
	// Duplicate donation
	oldDonation, err := model.CopyDonation(donation)
	if err != nil {
		s.Logger.Error("failed to duplicate donation for changelog", err)
		return err
	}

	// Update donation
	timestamp := time.Now()
	donation.StatusId = api.DonationPaymentOK
	donation.UpdatedAt = timestamp
	donation.ModifiedBy = &model.ModifierMeta{
		Id:       "SYSTEM",
		Role:     "SYSTEM",
		FullName: "SYSTEM",
	}
	donation.Version = oldDonation.Version + 1

	// Create changelog
	changelog := []string{
		"status_id",
	}

	// Update donation
	err = s.Repository.UpdateDonation(*oldDonation, *donation, changelog)
	if err != nil {
		return err
	}

	return nil
}

func (s *Initiative) List(opt dto.PageReq) ([]dto.InitiativeResp, error) {
	// Get initiative list
	initiativeList, err := s.Repository.FindActive(opt.Skip, opt.Limit)
	if err != nil {
		s.Logger.Error("unable to retrieve initiative list", err)
		return nil, err
	}

	// Compose response
	resp := make([]dto.InitiativeResp, len(initiativeList))
	for k, v := range initiativeList {
		// Get full image url
		imageUrls := dto.InitiativeImageFileResp{
			Thumbnail:  s.AssetService.GetPublicUrl(api.AssetInitiative, v.ImageFiles.Thumbnail),
			DetailPage: s.AssetService.GetPublicUrl(api.AssetInitiative, v.ImageFiles.DetailPage),
		}

		// Convert tags to array
		rawTags := strings.Trim(v.Tags, ",")
		tags := strings.Split(rawTags, ",")

		// Set response
		resp[k] = dto.InitiativeResp{
			Id:                 v.Id,
			Name:               v.Name,
			Description:        v.Description,
			ImageFiles:         imageUrls,
			Tags:               tags,
			ExternalUrls:       v.ExternalUrls,
			Price:              v.Price,
			CurrencyId:         v.CurrencyId,
			DonationConversion: v.DonationConversion,
			CurrencyName:       api.CurrencyName[v.CurrencyId],
			UpdatedAt:          v.UpdatedAt.Unix(),
			Version:            v.Version,
			Headline:           v.Headline.String,
		}
	}
	return resp, nil
}
