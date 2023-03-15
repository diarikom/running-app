package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/pkg/nhttp"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/gorilla/mux"
	"net/http"
)

type SiteSettingHandler struct {
	SiteSettingService api.SiteSettingService
	Logger             nlog.Logger
}

func NewSiteSettingHandler(app *api.Api) SiteSettingHandler {
	return SiteSettingHandler{
		SiteSettingService: app.Services.SiteSetting,
		Logger:             app.Logger,
	}
}

func (h *SiteSettingHandler) StaticContent(r *http.Request) (success *nhttp.Success, err error) {
	contentType := mux.Vars(r)["type"]
	feature, err := h.SiteSettingService.StaticContent(contentType)

	if err != nil {
		return
	}

	resp := &nhttp.Success{
		Result: feature,
	}

	return resp, nil
}
