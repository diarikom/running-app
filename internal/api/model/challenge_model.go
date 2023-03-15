package model

import (
	"github.com/diarikom/running-app/running-app-api/internal/pkg/ngrule"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"time"
)

type Challenge struct {
	Id          string      `db:"id" json:"id"`
	MilestoneId string      `db:"milestone_id" json:"milestone_id"`
	Title       string      `db:"title" json:"title"`
	Description string      `db:"description" json:"description"`
	Level       int         `db:"level" json:"level"`
	Status      int         `db:"status" json:"status"`
	Rules       ngrule.Rule `db:"rules" json:"rules"`
	Sort        int         `db:"sort" json:"sort"`
	CreatedAt   time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time   `db:"updated_at" json:"updated_at"`
	Version     int         `db:"version" json:"version"`
}

type ChallengeRule struct {
	RuleBuilder *builder.RuleBuilder
	Params      []ngrule.FactParam
	Targets     []ngrule.FactParam
}
