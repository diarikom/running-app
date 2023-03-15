package model

import (
	"time"
)

type AdTag struct {
	Id        int       `db:"id" diff:"id"`
	Name      string    `db:"name"`
	UpdatedAt time.Time `db:"updated_at" diff:"required"`
}
