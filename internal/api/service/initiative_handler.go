package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/dto"
	"github.com/diarikom/running-app/running-app-api/pkg/nhttp"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/gorilla/mux"
	"net/http"
)

func NewInitiativeHandler(app *api.Api) InitiativeHandler {
	return InitiativeHandler{
		InitiativeService: app.Services.Initiative,
		Logger:            app.Logger,
	}
}

type InitiativeHandler struct {
	InitiativeService api.InitiativeService
	Logger            nlog.Logger
}

func (h *InitiativeHandler) List(req *http.Request) (*nhttp.Success, error) {
	// Get skip and limit
	query := req.URL.Query()
	skip, limit := api.Pagination(query)

	// Call service
	respBody, err := h.InitiativeService.List(dto.PageReq{
		Skip:  skip,
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}

	resp := nhttp.Success{
		Result: respBody,
	}
	return &resp, nil
}

func (h *InitiativeHandler) Donate(req *http.Request) (*nhttp.Success, error) {
	// Parse request body
	var reqBody dto.DonateReq
	err := nhttp.ParseJSON(&reqBody, req)
	if err != nil {
		return nil, nhttp.ErrBadRequest
	}

	// Set user id and initiative id
	reqBody.UserId = req.Header.Get(nhttp.KeyUserId)
	reqBody.InitiativeId = mux.Vars(req)["id"]

	// Call service
	respBody, err := h.InitiativeService.Donate(reqBody)
	if err != nil {
		return nil, err
	}

	resp := nhttp.Success{
		Result: respBody,
	}
	return &resp, nil
}

func (h *InitiativeHandler) ListUserDonation(req *http.Request) (*nhttp.Success, error) {
	// Get skip and limit
	query := req.URL.Query()
	skip, limit := api.Pagination(query)

	// Call service
	respBody, err := h.InitiativeService.ListUserDonation(dto.UserResourcesReq{
		PageReq: dto.PageReq{
			Skip:  skip,
			Limit: limit,
		},
		UserId: req.Header.Get(nhttp.KeyUserId),
	})
	if err != nil {
		return nil, err
	}

	resp := nhttp.Success{
		Result: respBody,
	}
	return &resp, nil
}
