package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/pkg/nhttp"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"net/http"
)

func NewDiscoverContentHandler(app *api.Api) DiscoverContentHandler {
	return DiscoverContentHandler{
		DiscoverContentService: app.Services.DiscoverContent,
		Logger:                 app.Logger,
	}
}

type DiscoverContentHandler struct {
	DiscoverContentService api.DiscoverContentService
	Logger                 nlog.Logger
}

func (h *DiscoverContentHandler) GetContents(r *http.Request) (success *nhttp.Success, err error) {
	// Get query
	q := r.URL.Query()
	// Get pagination
	skip, limit := api.Pagination(q)

	// Call service
	sessions, err := h.DiscoverContentService.GetContents(skip, limit)
	if err != nil {
		return
	}

	// Compose response
	resp := &nhttp.Success{
		Result: sessions,
	}

	return resp, nil
}
