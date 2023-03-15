package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/model"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
)

type SubscriptionPlanRepository struct {
	IdGen  *api.SnowflakeGen
	Errors *api.Errors
	Db     *nsql.SqlDatabase
	Stmt   SubscriptionPlanStatement
	Logger nlog.Logger
}

func NewSubscriptionPlanRepository(db *nsql.SqlDatabase, idGen *api.SnowflakeGen, errors *api.Errors, logger nlog.Logger) api.SubscriptionPlanRepository {
	r := SubscriptionPlanRepository{
		IdGen:  idGen,
		Errors: errors,
		Db:     db,
		Stmt:   initSubscriptionPlanStatement(db),
		Logger: logger,
	}

	return &r
}

func (r *SubscriptionPlanRepository) SubscriptionPlans() (result []model.SubscriptionPlan, err error) {
	err = r.Stmt.list.Select(&result)

	return result, err
}
