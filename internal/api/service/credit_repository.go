package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/model"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
)

func NewCreditRepository(db *nsql.SqlDatabase, logger nlog.Logger, errComponent *api.Errors) api.CreditRepository {
	r := creditRepository{
		Db:     db,
		Stmt:   initCreditStatement(db),
		Logger: logger,
		Errors: errComponent,
	}

	return &r
}

type creditRepository struct {
	Db     *nsql.SqlDatabase
	Stmt   creditStatements
	Logger nlog.Logger
	Errors *api.Errors
}

func (c *creditRepository) FindWalletByUser(userId string) (*model.UserCreditWallet, error) {
	var w model.UserCreditWallet
	err := c.Stmt.findWalletByUser.Get(&w, userId)
	return &w, err
}

func (c *creditRepository) FindWalletById(walletId string) (*model.UserCreditWallet, error) {
	var w model.UserCreditWallet
	err := c.Stmt.findWalletById.Get(&w, walletId)
	return &w, err
}

func (c *creditRepository) InsertWallet(wallet *model.UserCreditWallet, initTrx *model.UserCreditWalletTrx) error {
	// Begin transaction
	trx, err := c.Db.Conn.Beginx()
	if err != nil {
		return err
	}
	defer nsql.ReleaseTx(trx, &err, c.Logger)

	// Insert wallet
	_, err = trx.NamedStmt(c.Stmt.insertWallet).Exec(&wallet)
	if err != nil {
		c.Logger.Error("update credit wallet", err)
		return err
	}

	// Add new credit transaction
	_, err = trx.NamedStmt(c.Stmt.insertTrx).Exec(&initTrx)
	if err != nil {
		c.Logger.Error("insert credit trx", err)
		return err
	}

	return nil
}

func (c *creditRepository) IsExistWalletByUser(userId string) (bool, error) {
	var isExist bool
	err := c.Stmt.isExistWalletByUser.Get(&isExist, userId)
	return isExist, err
}

func (c *creditRepository) IsExistTrxRef(walletId, trxId string) (bool, error) {
	var isExist bool
	err := c.Stmt.isExistTrxRef.Get(&isExist, walletId, trxId)
	return isExist, err
}

func (c *creditRepository) InsertTrx(wallet *model.UserCreditWallet, newTrx *model.UserCreditWalletTrx) error {
	// Begin transaction
	trx, err := c.Db.Conn.Beginx()
	if err != nil {
		return err
	}
	defer nsql.ReleaseTx(trx, &err, c.Logger)

	// Add new credit transaction
	_, err = trx.NamedStmt(c.Stmt.insertTrx).Exec(&newTrx)
	if err != nil {
		c.Logger.Error("insert credit trx", err)
		return err
	}

	// Update user wallet
	result, err := trx.NamedStmt(c.Stmt.updateWalletBalance).Exec(&wallet)
	if err != nil {
		c.Logger.Error("update credit wallet", err)
		return err
	}

	// Check for affected rows
	count, err := result.RowsAffected()
	if err != nil {
		c.Logger.Error("cannot get affected rows", err)
		return err
	}

	if count == 0 {
		c.Logger.Errorf("no wallet update affected. Rolling back")
		return c.Errors.New("")
	}

	return nil
}

func (c *creditRepository) FindTrxById(trxId string) (*model.UserCreditWalletTrx, error) {
	var t model.UserCreditWalletTrx
	err := c.Stmt.findTrxById.Get(&t, trxId)
	return &t, err
}

func (c *creditRepository) FindWalletByTrx(trxId string) (*model.UserCreditWallet, error) {
	var w model.UserCreditWallet
	err := c.Stmt.findWalletByTrx.Get(&w, trxId)
	return &w, err
}
