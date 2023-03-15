package service

import (
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
	"github.com/jmoiron/sqlx"
)

type discoverContentStatements struct {
	countContent *sqlx.Stmt
	findContent  *sqlx.Stmt
}

func initDiscoverContentStatements(db *nsql.SqlDatabase) discoverContentStatements {
	return discoverContentStatements{
		countContent: db.Prepare(`SELECT COUNT(id) FROM discover_content WHERE status_id = 1`),
		findContent:  db.Prepare(`SELECT discover_content.id,COALESCE(organization.name , '') as title,discover_content.content_body,discover_content.logo_file,discover_content.external_url,discover_content.status_id,discover_content.sort,discover_content.created_at,discover_content.updated_at,discover_content.version,discover_content.modified_by,discover_content.image_files,discover_content.tags,discover_content.headline,discover_content.image_files FROM discover_content LEFT JOIN organization on discover_content.organization_id = organization.id WHERE discover_content.status_id = 1 ORDER BY sort LIMIT $1 OFFSET $2`),
	}
}
