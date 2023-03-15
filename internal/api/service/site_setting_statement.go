package service

import (
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
	"github.com/jmoiron/sqlx"
)

type SiteSettingStatement struct {
	staticContent *sqlx.Stmt
}

func initSiteSettingStatement(db *nsql.SqlDatabase) SiteSettingStatement {
	return SiteSettingStatement{
		staticContent: db.Prepare(`SELECT value FROM site_setting WHERE key = $1`),
	}
}
