package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/pkg/nhttp"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"net/http"
)

type SubscriptionPlanHandler struct {
	SubscriptionPlanService api.SubscriptionPlanService
	Logger                  nlog.Logger
}

func NewSubscriptionPlanHandler(app *api.Api) SubscriptionPlanHandler {
	return SubscriptionPlanHandler{
		SubscriptionPlanService: app.Services.SubscriptionPlan,
		Logger:                  app.Logger,
	}
}

func (h *SubscriptionPlanHandler) List(req *http.Request) (*nhttp.Success, error) {
	// Get list of subscription plan;
	r, err := h.SubscriptionPlanService.List()

	// Handler if subscription plan doesnt exist;
	if err != nil {
		return nil, err
	}

	// Set data return;
	resp := &nhttp.Success{
		Result: r,
	}

	return resp, nil
}
