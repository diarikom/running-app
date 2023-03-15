package model

import (
	"encoding/json"
	"time"
)

type SiteSetting struct {
	Id          string          `db:"id"`
	Key         string          `db:"key"`
	Value       string          `db:"value"`
	Description string          `db:"description"`
	Field       json.RawMessage `db:"field"`
	Version     int8            `db:"version"`
	UpdatedAt   time.Time       `db:"updated_at"`
	ModifiedBy  json.RawMessage `db:"modified_by"`
}
