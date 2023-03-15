package model

import (
	"time"
)

type RunSession struct {
	Id             string    `db:"id" diff:"id"`
	UserId         string    `db:"user_id" diff:"user_id"`
	SessionStarted time.Time `db:"session_started"`
	SessionEnded   time.Time `db:"session_ended"`
	TimeElapsed    int       `db:"time_elapsed"`
	Distance       int       `db:"distance"`
	Speed          float64   `db:"speed"`
	StepCount      int       `db:"step_count"`
	SyncStatusId   int       `db:"sync_status_id"`
	CreatedAt      time.Time `db:"created_at" diff:"-"`
	UpdatedAt      time.Time `db:"updated_at" diff:"required"`
	Version        int       `db:"version" diff:"required"`
}
