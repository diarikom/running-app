package nsql

import (
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/jmoiron/sqlx"
)

// NamedStmtTx determines whether to use transaction or not
func NamedStmtTx(s *sqlx.NamedStmt, tx *sqlx.Tx) *sqlx.NamedStmt {
	if tx != nil {
		return tx.NamedStmt(s)
	}
	return s
}

// StmtTx determines whether to use transaction or not
func StmtTx(s *sqlx.Stmt, tx *sqlx.Tx) *sqlx.Stmt {
	if tx != nil {
		return tx.Stmtx(s)
	}
	return s
}

// ReleaseTx clean db transaction by commit if no error, or rollback if an error occurred
func ReleaseTx(tx *sqlx.Tx, err *error, log nlog.Logger) {
	if *err != nil {
		// If an error occurred, rollback transaction
		errRollback := tx.Rollback()
		if errRollback != nil {
			log.Error("Unable to rollback transaction", errRollback)
		} else {
			log.Debug("Transaction rolled back")
		}
		return
	}
	// Else, commit transaction
	errCommit := tx.Commit()
	if errCommit != nil {
		log.Error("Unable to commit transaction", errCommit)
	}
}
