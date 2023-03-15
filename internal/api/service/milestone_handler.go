package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/dto"
	"github.com/diarikom/running-app/running-app-api/pkg/nhttp"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"net/http"
	"time"
)

type MilestoneHandler struct {
	MilestoneService api.MilestoneService
	Logger           nlog.Logger
}

func NewMilestoneHandler(app *api.Api) MilestoneHandler {
	return MilestoneHandler{
		MilestoneService: app.Services.MilestoneService,
		Logger:           app.Logger,
	}
}

func (h *MilestoneHandler) Current(r *http.Request) (success *nhttp.Success, err error) {
	userID := r.Header.Get(nhttp.KeyUserId)
	feature, err := h.MilestoneService.Current(userID)

	if err != nil {
		return
	}

	resp := &nhttp.Success{
		Result: feature,
	}

	return resp, nil
}

func (h *MilestoneHandler) CheckChallengeAchieve(r *http.Request) (success *nhttp.Success, err error) {
	// Get user id from body
	var reqBody dto.UserChallengeReq
	err = nhttp.ParseJSON(&reqBody, r)
	if err != nil {
		return nil, nhttp.ErrBadRequest
	}

	// Validate body
	if reqBody.UserId == "" {
		return nil, nhttp.ErrBadRequest
	}

	// Set timestamp
	reqBody.Timestamp = time.Now()

	// Call service
	err = h.MilestoneService.TriggerCheckChallengeAchieved(reqBody)

	if err != nil {
		return
	}

	return nhttp.OK(), nil
}

func (h *MilestoneHandler) ReloadMilestone(r *http.Request) (success *nhttp.Success, err error) {
	h.MilestoneService.LoadMilestone()
	return nhttp.OK(), nil
}
