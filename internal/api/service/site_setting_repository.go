package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/model"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
)

type SiteSettingRepository struct {
	IdGen  *api.SnowflakeGen
	Errors *api.Errors
	Db     *nsql.SqlDatabase
	Stmt   SiteSettingStatement
	Logger nlog.Logger
}

func NewSiteSettingRepository(db *nsql.SqlDatabase, idGen *api.SnowflakeGen, errors *api.Errors, logger nlog.Logger) api.SiteSettingRepository {
	r := SiteSettingRepository{
		IdGen:  idGen,
		Errors: errors,
		Db:     db,
		Stmt:   initSiteSettingStatement(db),
		Logger: logger,
	}

	return &r
}

func (r *SiteSettingRepository) StaticContent(contentType string) (result model.SiteSetting, err error) {
	err = r.Stmt.staticContent.Get(&result, contentType)

	return result, err
}
