package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
	"time"
)

type DiscoverContent struct {
	Id          string             `db:"id" diff:"id"`
	Title       string             `db:"title"`
	ContentBody string             `db:"content_body"`
	LogoFile    sql.NullString     `db:"logo_file"`
	ExternalUrl json.RawMessage    `db:"external_url"`
	StatusId    int                `db:"status_id"`
	Sort        int                `db:"sort"`
	CreatedAt   time.Time          `db:"created_at" diff:"-"`
	UpdatedAt   time.Time          `db:"updated_at" diff:"required"`
	Version     int                `db:"version" diff:"required"`
	ModifiedBy  json.RawMessage    `db:"modified_by"`
	ImageFiles  DiscoverImageFiles `db:"image_files"`
	Headline    sql.NullString     `db:"headline"`
	Tags        sql.NullString     `db:"tags"`
}

type DiscoverImageFiles struct {
	Thumbnail  string `json:"thumbnail"`
	DetailPage string `json:"detail_page"`
}

func (i *DiscoverImageFiles) Scan(src interface{}) error {
	return nsql.ScanJSON(src, i)
}

func (i DiscoverImageFiles) Value() (driver.Value, error) {
	return json.Marshal(i)
}
