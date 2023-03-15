package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/dto"
	"github.com/diarikom/running-app/running-app-api/pkg/nhttp"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/diarikom/running-app/running-app-api/pkg/nstr"
	"net/http"
)

func NewAssetHandler(app *api.Api) AssetHandler {
	return AssetHandler{
		AssetService: app.Services.Asset,
		Logger:       app.Logger,
	}
}

type AssetHandler struct {
	AssetService api.AssetService
	Logger       nlog.Logger
}

func (h *AssetHandler) PostUploadFile(r *http.Request) (*nhttp.Success, error) {
	q := r.URL.Query()

	// Get asset type
	assetType := nstr.ParseInt(q.Get("asset_type"), 0)

	// Get rules by asset type
	rule, err := h.AssetService.GetUploadRule(assetType)
	if err != nil {
		return nil, err
	}

	// Parse multipart file
	file, err := nhttp.GetFile(r, rule.Key, rule.MaxSize, rule.MimeTypes)
	if err != nil {
		return nil, err
	}

	// Upload file
	resp, err := h.AssetService.UploadFile(dto.UploadReq{
		AssetType: assetType,
		File:      file,
	})

	return &nhttp.Success{
		Result: resp,
	}, nil
}
