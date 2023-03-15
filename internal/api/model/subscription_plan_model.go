package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

type SubscriptionPlan struct {
	Id             string          `db:"id"`
	ProviderId     string          `db:"provider_id"`
	Provider       string          `db:"provider"`
	PlanTypeId     string          `db:"plan_type_id"`
	PlanType       string          `db:"plan_type"`
	ProviderTrxRef string          `db:"provider_trx_ref"`
	Options        json.RawMessage `db:"options"`
	Description    sql.NullString  `db:"description"`
	CreatedAt      time.Time       `db:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at"`
}
