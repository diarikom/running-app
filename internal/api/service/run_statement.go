package service

import (
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
	"github.com/jmoiron/sqlx"
)

type runStatements struct {
	countRunSessionHistory *sqlx.Stmt
	findRunSessionHistory  *sqlx.Stmt
	insertRunSession       *sqlx.NamedStmt
	updateRunSyncStatus    *sqlx.Stmt
	sumRunSessionDistance  *sqlx.Stmt
}

func initRunStatements(db *nsql.SqlDatabase) runStatements {
	return runStatements{
		countRunSessionHistory: db.Prepare(`SELECT COUNT(id) FROM run_session WHERE user_id = $1`),
		findRunSessionHistory:  db.Prepare(`SELECT id, session_started, session_ended, time_elapsed, distance, speed, step_count, sync_status_id, created_at, updated_at, version FROM run_session WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`),
		insertRunSession:       db.PrepareNamed(`INSERT INTO run_session(id, user_id, session_started, session_ended, time_elapsed, distance, speed, step_count, sync_status_id, created_at) VALUES (:id, :user_id, :session_started, :session_ended, :time_elapsed, :distance, :speed, :step_count, :sync_status_id, :created_at)`),
		updateRunSyncStatus:    db.Prepare(`UPDATE run_session SET sync_status_id = $1 WHERE id = $2 AND user_id = $3`),
		sumRunSessionDistance:  db.Prepare(`SELECT SUM(distance) FROM run_session WHERE date(session_started) >= $2 AND DATE(session_ended) <= $3 AND user_id = $1`),
	}
}
