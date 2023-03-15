package service

import (
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
	"github.com/jmoiron/sqlx"
)

type creditStatements struct {
	findTrxById         *sqlx.Stmt
	findWalletById      *sqlx.Stmt
	findWalletByTrx     *sqlx.Stmt
	findWalletByUser    *sqlx.Stmt
	insertTrx           *sqlx.NamedStmt
	insertWallet        *sqlx.NamedStmt
	isExistTrxRef       *sqlx.Stmt
	isExistWalletByUser *sqlx.Stmt
	updateWalletBalance *sqlx.NamedStmt
}

func initCreditStatement(db *nsql.SqlDatabase) creditStatements {
	return creditStatements{
		findTrxById:         db.Prepare(`SELECT id, user_credit_wallet_id, balance, balance_pending, amount, trx_entry_type_id, trx_ref_id, notes, status, created_at, expired_at, version FROM user_credit_wallet_trx WHERE id = $1`),
		findWalletById:      db.Prepare(`SELECT w.id, w.user_id, w.balance, w.balance_pending, w.balance_expiring, w.balance_expiring_date, w.created_at, w.updated_at, w.version FROM user_credit_wallet AS w WHERE w.id = $1`),
		findWalletByTrx:     db.Prepare(`SELECT w.id, w.user_id, w.balance, w.balance_pending, w.balance_expiring, w.balance_expiring_date, w.created_at, w.updated_at, w.version FROM user_credit_wallet AS w INNER JOIN user_credit_wallet_trx t on w.id = t.user_credit_wallet_id WHERE t.id = $1`),
		findWalletByUser:    db.Prepare(`SELECT w.id, w.user_id, w.balance, w.balance_pending, w.balance_expiring, w.balance_expiring_date, w.created_at, w.updated_at, w.version FROM user_credit_wallet AS w WHERE w.user_id = $1`),
		insertTrx:           db.PrepareNamed(`INSERT INTO user_credit_wallet_trx(id, user_credit_wallet_id, balance, balance_pending, amount, trx_entry_type_id, trx_ref_id, notes, status, created_at, expired_at, version) VALUES (:id, :user_credit_wallet_id, :balance, :balance_pending, :amount, :trx_entry_type_id, :trx_ref_id, :notes, :status, :created_at, :expired_at, :version)`),
		insertWallet:        db.PrepareNamed(`INSERT INTO user_credit_wallet(id, user_id, balance, balance_pending, balance_expiring, balance_expiring_date, created_at, updated_at, version) VALUES (:id, :user_id, :balance, :balance_pending, :balance_expiring, :balance_expiring_date, :created_at, :updated_at, :version)`),
		isExistTrxRef:       db.Prepare(`SELECT COUNT(*) > 0 as "isExist" FROM user_credit_wallet_trx WHERE user_credit_wallet_id = $1 AND trx_ref_id = $2`),
		isExistWalletByUser: db.Prepare(`SELECT COUNT(*) > 0 as "isExist" FROM user_credit_wallet WHERE user_id = $1`),
		updateWalletBalance: db.PrepareNamed(`UPDATE user_credit_wallet SET balance = :balance, balance_pending = :balance_pending, updated_at = :updated_at, version = :version WHERE id = :id AND version = :current_version`),
	}
}
