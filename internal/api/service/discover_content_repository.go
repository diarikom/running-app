package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/model"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
)

func NewDiscoverContentRepository(db *nsql.SqlDatabase, logger nlog.Logger) api.DiscoverContentRepository {
	r := discoverContentRepository{
		Db:     db,
		Stmt:   initDiscoverContentStatements(db),
		Logger: logger,
	}

	return &r
}

type discoverContentRepository struct {
	Db     *nsql.SqlDatabase
	Stmt   discoverContentStatements
	Logger nlog.Logger
}

func (r discoverContentRepository) CountContents() (total int, err error) {
	total = 0
	err = r.Stmt.countContent.Get(&total)
	if err != nil {
		return
	}
	return
}

func (r discoverContentRepository) FindContents(limit int8, skip int64) (result []model.DiscoverContent, err error) {
	err = r.Stmt.findContent.Select(&result, limit, skip)
	if err != nil {
		return
	}
	if len(result) == 0 {
		result = []model.DiscoverContent{}
	}
	return
}
