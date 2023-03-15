package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/pkg/nhttp"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"net/http"
)

type AdTagHandler struct {
	AdTagService api.AdTagService
	Logger       nlog.Logger
}

func NewTagHandler(app *api.Api) AdTagHandler {
	return AdTagHandler{
		AdTagService: app.Services.Tag,
		Logger:       app.Logger,
	}
}

func (h *AdTagHandler) GetTags(r *http.Request) (success *nhttp.Success, err error) {
	sessions, err := h.AdTagService.GetAdTags()
	if err != nil {
		return
	}

	resp := &nhttp.Success{
		Result: sessions,
	}

	return resp, nil
}
