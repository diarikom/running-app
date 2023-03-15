package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/model"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
)

func NewInitiativeRepository(db *nsql.SqlDatabase, idGen *api.SnowflakeGen, apiErrors *api.Errors, logger nlog.Logger) api.InitiativeRepository {
	r := InitiativeRepository{
		IdGen:   idGen,
		Errors:  apiErrors,
		Db:      db,
		Stmt:    initInitiativeStatement(db),
		Differs: initInitiativeDiffer(),
		Logger:  logger,
	}

	return &r
}

type InitiativeRepository struct {
	IdGen   *api.SnowflakeGen
	Errors  *api.Errors
	Db      *nsql.SqlDatabase
	Stmt    InitiativeStatement
	Differs initiativeDiffer
	Logger  nlog.Logger
}

func (r *InitiativeRepository) FindDonationByUser(userId string, skip int64, limit int8) ([]model.Donation, error) {
	rows := make([]model.Donation, 0)
	err := r.Stmt.findDonationByUser.Select(&rows, userId, limit, skip)
	return rows, err
}

func (r *InitiativeRepository) UpdateDonation(oldDonation, newDonation model.Donation, changelog []string) error {
	// Get differ
	differ := r.Differs.donation

	// Compare instance
	diff, err := differ.Compare(oldDonation, newDonation, changelog)
	if err != nil {
		return err
	}

	// If no changes, return
	if diff.Count == 0 {
		r.Logger.Debug("no donation changes detected")
		return nil
	}

	// Generate insert donation query
	updateQuery, updateArgs, err := differ.UpdateQuerySafe(diff, oldDonation.Version)
	if err != nil {
		r.Logger.Error("unable to generate update donation query", err)
		return err
	}

	// Generate insert donation log query
	insertLogQuery, insertLogArgs, err := differ.InsertLogQuery(diff, r.IdGen.New(), changelog)
	if err != nil {
		r.Logger.Error("unable to generate insert donation_log query", err)
		return err
	}

	// Rebind all query
	updateQuery = r.Db.Conn.Rebind(updateQuery)
	insertLogQuery = r.Db.Conn.Rebind(insertLogQuery)

	// Begin transaction
	trx, err := r.Db.Conn.Beginx()
	if err != nil {
		return err
	}
	defer nsql.ReleaseTx(trx, &err, r.Logger)

	// Update Donation
	result, err := trx.Exec(updateQuery, updateArgs...)
	if err != nil {
		r.Logger.Error("failed to insert donation", err)
		return err
	}

	// Check for affected rows
	count, err := result.RowsAffected()
	if err != nil {
		r.Logger.Error("cannot get affected rows", err)
		return err
	}

	if count == 0 {
		r.Logger.Errorf("no donation update affected. Rolling back")
		err = r.Errors.New("INT002")
		return err
	}

	// Insert Log
	_, err = trx.Exec(insertLogQuery, insertLogArgs...)
	if err != nil {
		r.Logger.Error("failed to insert donation_log", err)
		return err
	}

	// Clean error
	err = nil
	return nil
}

func (r *InitiativeRepository) Insert(donation model.Donation, donationLog model.DonationLog) error {
	// Begin transaction
	trx, err := r.Db.Conn.Beginx()
	if err != nil {
		return err
	}
	defer nsql.ReleaseTx(trx, &err, r.Logger)

	// Insert donation
	_, err = trx.NamedStmt(r.Stmt.insertDonation).Exec(&donation)
	if err != nil {
		r.Logger.Error("failed to insert donation", err)
		return err
	}

	// Insert donation log
	_, err = trx.NamedStmt(r.Stmt.insertDonationLog).Exec(&donationLog)
	if err != nil {
		r.Logger.Error("failed to insert donation_log", err)
		return err
	}

	return nil
}

func (r *InitiativeRepository) FindById(id string) (*model.Initiative, error) {
	var result model.Initiative
	err := r.Stmt.findById.Get(&result, id)
	return &result, err
}

func (r *InitiativeRepository) FindActive(skip int64, limit int8) ([]model.Initiative, error) {
	result := make([]model.Initiative, 0)
	err := r.Stmt.findActive.Select(&result, limit, skip)
	return result, err
}
