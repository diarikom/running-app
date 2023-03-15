package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/model"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
)

type AdTagRepository struct {
	Db     *nsql.SqlDatabase
	Stmt   AdTagStatements
	Logger nlog.Logger
}

func NewAdTagRepository(db *nsql.SqlDatabase, log nlog.Logger) api.AdTagRepository {
	r := AdTagRepository{
		Db:     db,
		Stmt:   initAdTagStatements(db),
		Logger: log,
	}

	return &r
}

func (r AdTagRepository) GetAdTags() (result []model.AdTag, err error) {
	err = r.Stmt.getAdTags.Select(&result)

	return result, err
}
