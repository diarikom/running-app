package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/dto"
	"github.com/diarikom/running-app/running-app-api/pkg/nhttp"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRunHandler(app *api.Api) RunHandler {
	return RunHandler{
		RunService: app.Services.Run,
		Logger:     app.Logger,
	}
}

type RunHandler struct {
	RunService api.RunService
	Logger     nlog.Logger
}

func (h *RunHandler) GetRunSessionHistory(r *http.Request) (success *nhttp.Success, err error) {
	// Get user id
	userId := r.Header.Get(nhttp.KeyUserId)
	// Get query
	q := r.URL.Query()
	// Get pagination
	skip, limit := api.Pagination(q)

	// Call service
	sessions, err := h.RunService.GetRunSessionHistory(userId, skip, limit)
	if err != nil {
		return
	}

	// Compose response
	resp := &nhttp.Success{
		Result: sessions,
	}

	return resp, nil
}

func (h *RunHandler) PostRunSummary(r *http.Request) (*nhttp.Success, error) {
	// Get Request
	var reqBody dto.RunSessionReq
	err := nhttp.ParseJSON(&reqBody, r)
	if err != nil {
		return nil, nhttp.ErrBadRequest
	}

	userId := r.Header.Get(nhttp.KeyUserId)

	// Call service
	err = h.RunService.NewRunSession(userId, &reqBody)
	if err != nil {
		return nil, err
	}

	// Compose response
	resp := nhttp.OK()

	// Return response
	return resp, nil
}

func (h *RunHandler) PutRunStatusSync(r *http.Request) (*nhttp.Success, error) {
	// Get Request
	var reqBody dto.RunStatusSyncReq
	err := nhttp.ParseJSON(&reqBody, r)
	if err != nil {
		return nil, nhttp.ErrBadRequest
	}

	id := mux.Vars(r)["id"]

	userId := r.Header.Get(nhttp.KeyUserId)

	// Call service
	err = h.RunService.UpdateRunSyncStatus(id, userId, reqBody.Status)
	if err != nil {
		return nil, err
	}

	// Compose response
	resp := nhttp.OK()

	// Return response
	return resp, nil
}
