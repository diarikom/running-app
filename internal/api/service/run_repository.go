package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/model"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
	"time"
)

func NewRunRepository(db *nsql.SqlDatabase, logger nlog.Logger) api.RunRepository {
	r := runRepository{
		Db:     db,
		Stmt:   initRunStatements(db),
		Logger: logger,
	}

	return &r
}

type runRepository struct {
	Db     *nsql.SqlDatabase
	Stmt   runStatements
	Logger nlog.Logger
}

func (r runRepository) CountRunSessionHistory(userId string) (total int, err error) {
	total = 0
	err = r.Stmt.countRunSessionHistory.Get(&total, userId)
	if err != nil {
		return
	}
	return
}

func (r runRepository) FindRunSessionHistory(userId string, limit int8, skip int64) (result []model.RunSession, err error) {
	err = r.Stmt.findRunSessionHistory.Select(&result, userId, limit, skip)
	if err != nil {
		return
	}
	if len(result) == 0 {
		result = []model.RunSession{}
	}
	return
}

func (r runRepository) InsertRunSession(session model.RunSession) error {
	_, err := r.Stmt.insertRunSession.Exec(&session)
	return err
}

func (r runRepository) UpdateRunSyncStatus(id, userId string, syncStatus int) error {
	_, err := r.Stmt.updateRunSyncStatus.Exec(syncStatus, id, userId)
	return err
}

func (r runRepository) SumRunSessionDistance(userID string, start time.Time, end time.Time) (result int, err error) {
	err = r.Stmt.sumRunSessionDistance.Get(&result, userID, start, end)

	return result, err
}
