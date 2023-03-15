package service

import (
	"database/sql"
	"errors"
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/dto"
	"github.com/diarikom/running-app/running-app-api/internal/api/model"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/lib/pq"
	"time"
)

type CreditService struct {
	IdGen      *api.SnowflakeGen
	Error      *api.Errors
	Logger     nlog.Logger
	Repository api.CreditRepository
}

func (s *CreditService) Init(app *api.Api) error {
	s.IdGen = app.Components.Id
	s.Error = app.Components.Errors
	s.Logger = app.Logger
	s.Repository = NewCreditRepository(app.Datasources.Db, app.Logger, app.Components.Errors)
	return nil
}

func (s *CreditService) CheckChargeAmount(opt dto.CreditChargeOpt) (int, error) {
	// Get wallet
	wallet, err := s.GetUserWallet(opt.UserId)
	if err != nil {
		return 0, err
	}

	// Check against balance
	if wallet.Balance < opt.Amount {
		return 0, s.Error.New("CRD010")
	}

	// Else, return wallet version for concurrency control
	return wallet.Version, nil
}

func (s *CreditService) Charge(opt dto.CreditChargeOpt) (string, error) {
	// Get user wallet
	wallet, err := s.GetUserWallet(opt.UserId)
	if err != nil {
		return "", err
	}

	// If version is set and it is different with current wallet, then re-check charged amount
	if opt.WalletVersion > 0 &&
		(wallet.Version != opt.WalletVersion) &&
		(wallet.Balance < opt.Amount) {
		return "", s.Error.New("CRD010")
	}

	// Insert pending transaction
	trxId, err := s.InsertPendingTrx(dto.CreditTrxOpt{
		WalletId:  wallet.Id,
		Amount:    opt.Amount,
		EntryType: api.Credit,
	})

	return trxId, nil
}

func (s *CreditService) GetUserBalance(userId string) (*dto.UserCreditBalanceResp, error) {
	// Get user wallet
	wallet, err := s.GetUserWallet(userId)
	if err != nil {
		return nil, err
	}

	// Determine expire time null
	var expireTime int64
	if !wallet.BalanceExpiringDate.Valid {
		expireTime = 0
	} else {
		expireTime = wallet.BalanceExpiringDate.Time.Unix()
	}

	// Compose response
	resp := dto.UserCreditBalanceResp{
		Balance:         wallet.Balance,
		PendingBalance:  wallet.BalancePending,
		ExpiringBalance: wallet.BalanceExpiring,
		ExpireTime:      expireTime,
	}

	return &resp, nil
}

func (s *CreditService) InsertPendingTrx(opt dto.CreditTrxOpt) (string, error) {
	// If amount is less
	if opt.Amount <= 0 {
		return "", errors.New("invalid amount")
	}

	// Get wallet
	wallet, err := s.Repository.FindWalletById(opt.WalletId)
	if err != nil {
		s.Logger.Error("unable to retrieve wallet by id", err)
		return "", err
	}

	// Update balance
	pendingBalance := wallet.BalancePending + opt.Amount

	// Update version
	currentVersion := wallet.Version
	version := wallet.Version + 1

	// Create timestamp
	var timestamp time.Time
	if opt.Timestamp == nil {
		timestamp = time.Now()
	} else {
		timestamp = *opt.Timestamp
	}

	// Set expired at
	var expiredAt pq.NullTime
	if opt.ExpiredAt != nil {
		expiredAt = pq.NullTime{Valid: true, Time: *opt.ExpiredAt}
	}

	// Set transaction id
	var trxId string
	if opt.Id == "" {
		trxId = s.IdGen.New()
	} else {
		trxId = opt.Id
	}

	// Create trx
	pendingTrx := model.UserCreditWalletTrx{
		Id:                 trxId,
		UserCreditWalletId: wallet.Id,
		Balance:            wallet.Balance,
		BalancePending:     pendingBalance,
		Amount:             opt.Amount,
		TrxEntryTypeId:     opt.EntryType,
		TrxRefId:           sql.NullString{Valid: false},
		Notes:              sql.NullString{Valid: opt.Notes != "", String: opt.Notes},
		Status:             api.TrxPending,
		CreatedAt:          timestamp,
		ExpiredAt:          expiredAt,
		Version:            version,
	}

	// Update wallet
	wallet.BalancePending = pendingBalance
	wallet.UpdatedAt = timestamp
	wallet.Version = version
	wallet.CurrentVersion = currentVersion

	// Insert trx
	err = s.Repository.InsertTrx(wallet, &pendingTrx)
	if err != nil {
		s.Logger.Error("unable to persist wallet insert", err)
		return "", err
	}

	return pendingTrx.Id, nil
}

func (s *CreditService) GetUserWallet(userId string) (*model.UserCreditWallet, error) {
	// Check if user has a wallet
	wallet, err := s.Repository.FindWalletByUser(userId)
	if err != nil && err != sql.ErrNoRows {
		s.Logger.Error("unable to check user wallet existence", err)
		return nil, err
	}

	// If exist, return wallet
	if err == nil {
		return wallet, nil
	}

	// Reset error
	err = nil

	// Create timestamp
	t := time.Now()

	// Init wallet
	wallet = &model.UserCreditWallet{
		Id:        s.IdGen.New(),
		UserId:    userId,
		CreatedAt: t,
		UpdatedAt: t,
		Version:   1,
	}

	// Init transaction
	trx := model.UserCreditWalletTrx{
		Id:                 s.IdGen.New(),
		UserCreditWalletId: wallet.Id,
		TrxEntryTypeId:     api.InitEntryType,
		Notes:              sql.NullString{Valid: true, String: "Init wallet"},
		Status:             api.TrxSuccess,
		CreatedAt:          t,
		Version:            wallet.Version,
	}

	// Persist wallet
	err = s.Repository.InsertWallet(wallet, &trx)
	if err != nil {
		s.Logger.Error("unable to persist wallet insert", err)
		return nil, err
	}

	return wallet, err
}

func (s *CreditService) SettlePendingTrx(opt dto.CreditSettleOpt) error {
	// Validate trx id
	if opt.TrxId == "" {
		return errors.New("TrxId is required")
	}

	// Get pending transaction
	pendingTrx, err := s.Repository.FindTrxById(opt.TrxId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = s.Error.New("CRD001")
			return err
		}

		s.Logger.Error("unable to retrieve credit transaction", err)
		return err
	}

	// Check status
	if pendingTrx.Status != api.TrxPending {
		return s.Error.New("CRD002")
	}

	// Check referenced pending transaction
	isExist, err := s.Repository.IsExistTrxRef(pendingTrx.UserCreditWalletId, opt.TrxId)
	if err != nil {
		s.Logger.Error("unable to check referenced transaction existence", err)
		return err
	}
	if isExist {
		return s.Error.New("CRD006")
	}

	// Check expired time
	if pendingTrx.ExpiredAt.Valid {
		now := time.Now().Unix()
		exp := pendingTrx.ExpiredAt.Time.Unix()
		if now > exp {
			s.Logger.Errorf("Pending transaction has been expired")
			return s.Error.New("CRD009")
		}
		// TODO: Settle pending failed pending transaction
	}

	// Check transaction type
	if pendingTrx.TrxEntryTypeId != api.Debit &&
		pendingTrx.TrxEntryTypeId != api.Credit {
		return s.Error.New("CRD003")
	}

	// Get wallet by transaction id
	wallet, err := s.Repository.FindWalletByTrx(pendingTrx.Id)
	if err != nil {
		s.Logger.Error("unable to retrieve wallet", err)
		return err
	}

	// Check pending balance in wallet
	if wallet.BalancePending <= 0 {
		return s.Error.New("CRD004")
	}

	// Check pending balance with pending transaction amount
	if wallet.BalancePending < pendingTrx.Amount {
		return s.Error.New("CRD005")
	}

	// Update balance
	balance := wallet.Balance
	if pendingTrx.TrxEntryTypeId == api.Debit {
		balance += pendingTrx.Amount
	} else if pendingTrx.TrxEntryTypeId == api.Credit {
		balance -= pendingTrx.Amount
	}

	// Update pending balance
	pendingBalance := wallet.BalancePending - pendingTrx.Amount

	// Update version
	currentVersion := wallet.Version
	version := wallet.Version + 1

	// Create timestamp
	var timestamp time.Time
	if opt.Timestamp == nil {
		timestamp = time.Now()
	} else {
		timestamp = *opt.Timestamp
	}

	// Create Transaction
	newTrx := model.UserCreditWalletTrx{
		Id:                 s.IdGen.New(),
		UserCreditWalletId: wallet.Id,
		Balance:            balance,
		BalancePending:     pendingBalance,
		Amount:             pendingTrx.Amount,
		TrxEntryTypeId:     pendingTrx.TrxEntryTypeId,
		TrxRefId:           sql.NullString{Valid: true, String: pendingTrx.Id},
		Notes:              sql.NullString{Valid: opt.Notes != "", String: opt.Notes},
		Status:             api.TrxSuccess,
		CreatedAt:          timestamp,
		ExpiredAt:          newNullTime(opt.ExpiredAt),
		Version:            version,
	}

	// Update wallet
	wallet.Balance = balance
	wallet.BalancePending = pendingBalance
	wallet.UpdatedAt = timestamp
	wallet.Version = version
	wallet.CurrentVersion = currentVersion

	// Persist updates
	err = s.Repository.InsertTrx(wallet, &newTrx)
	if err != nil {
		s.Logger.Error("unable to persist transaction insert", err)
		return err
	}

	return nil
}

func newNullTime(t *time.Time) pq.NullTime {
	if t == nil {
		return pq.NullTime{}
	}

	return pq.NullTime{
		Valid: true,
		Time:  *t,
	}
}
