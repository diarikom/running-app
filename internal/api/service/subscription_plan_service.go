package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/dto"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
)

type SubscriptionService struct {
	IdGen      *api.SnowflakeGen
	Errors     *api.Errors
	Logger     nlog.Logger
	Repository api.SubscriptionPlanRepository
}

func (s *SubscriptionService) Init(app *api.Api) error {
	s.IdGen = app.Components.Id
	s.Errors = app.Components.Errors
	s.Logger = app.Logger
	s.Repository = NewSubscriptionPlanRepository(app.Datasources.Db, app.Components.Id, app.Components.Errors, app.Logger)

	return nil
}

func (s *SubscriptionService) List() ([]dto.SubscriptionPlanResp, error) {
	// Get list of subscription plan;
	q, err := s.Repository.SubscriptionPlans()

	// If subscription plan doesnt exist;
	if err != nil {
		s.Logger.Error("unable to retreive subscription plan list", err)

		return nil, err
	}

	// Set response object;
	resp := make([]dto.SubscriptionPlanResp, len(q))
	for k, v := range q {
		resp[k] = dto.SubscriptionPlanResp{
			Id:          v.PlanTypeId,
			Name:        v.PlanType,
			Description: v.Description.String,
		}
	}

	return resp, nil
}
