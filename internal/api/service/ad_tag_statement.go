package service

import (
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
	"github.com/jmoiron/sqlx"
)

type AdTagStatements struct {
	getAdTags *sqlx.Stmt
}

func initAdTagStatements(db *nsql.SqlDatabase) AdTagStatements {
	return AdTagStatements{
		getAdTags: db.Prepare(`SELECT id, name, updated_at FROM ad_tag`),
	}
}
