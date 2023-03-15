package model

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
)

type ModifierMeta struct {
	Id       string `json:"id"`
	Role     string `json:"role"`
	FullName string `json:"full_name"`
}

func (m *ModifierMeta) Scan(src interface{}) error {
	return nsql.ScanJSON(src, m)
}

func (m ModifierMeta) Value() (driver.Value, error) {
	return json.Marshal(m)
}
