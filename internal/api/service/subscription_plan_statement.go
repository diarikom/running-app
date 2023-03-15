package service

import (
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
	"github.com/jmoiron/sqlx"
)

type SubscriptionPlanStatement struct {
	list *sqlx.Stmt
}

func initSubscriptionPlanStatement(db *nsql.SqlDatabase) SubscriptionPlanStatement {
	return SubscriptionPlanStatement{
		list: db.Prepare(`SELECT provider_subscription_plan.id,provider_subscription_plan.provider_id,m_provider.name as provider,provider_subscription_plan.plan_type_id,m_subscription_plan_type.name as plan_type,provider_subscription_plan.provider_trx_ref,provider_subscription_plan.options,provider_subscription_plan.description,provider_subscription_plan.created_at,provider_subscription_plan.updated_at FROM provider_subscription_plan JOIN m_provider on provider_subscription_plan.provider_id = m_provider.id JOIN m_subscription_plan_type on provider_subscription_plan.plan_type_id = m_subscription_plan_type.id`),
	}
}
