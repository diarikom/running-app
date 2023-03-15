package model

import "time"

type Milestone struct {
	Id          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	PeriodStart time.Time `db:"period_start" json:"period_start"`
	PeriodEnd   time.Time `db:"period_end" json:"period_end"`
	PeriodTZ    int       `db:"period_tz" json:"period_tz"`
	Status      int       `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	Version     int       `db:"version" json:"version"`
}
