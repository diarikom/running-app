package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/dto"
	"github.com/diarikom/running-app/running-app-api/pkg/nhttp"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"net/http"
)

type Asset struct {
	Datasource api.S3Provider
	IdGen      *api.SnowflakeGen
	Errors     *api.Errors
	Logger     nlog.Logger
	BaseUrl    string
}

func (s *Asset) Init(app *api.Api) error {
	// Init asset service
	s.Datasource = app.Datasources.Asset
	s.IdGen = app.Components.Id
	s.Errors = app.Components.Errors
	s.Logger = app.Logger
	s.BaseUrl = app.Config.GetString(api.ConfAssetBaseUrl)
	return nil
}

func (s *Asset) GetUploadRule(assetType int) (*nhttp.UploadRule, error) {
	var rule nhttp.UploadRule

	switch assetType {
	case api.AssetAvatarProfile:
		rule = nhttp.UploadRule{
			Key:       nhttp.DefaultKeyFile,
			MaxSize:   nhttp.MaxFileSizeImage,
			MimeTypes: nhttp.MimeTypesImage,
		}

	default:
		return nil, s.Errors.New("AST001")
	}

	return &rule, nil
}

func (s *Asset) UploadFile(req dto.UploadReq) (*dto.UploadResp, error) {
	// Determine dir
	dir, err := s.GetDir(req.AssetType)
	if err != nil {
		return nil, err
	}

	// Generate filename
	fileName := req.File.Rename(s.IdGen.New())

	// Set destination
	dest := dir + fileName

	// Upload file
	err = s.Datasource.Upload(req.File.File, req.File.MimeType, dest, api.AssetsPublicScope)
	if err != nil {
		s.Logger.Error("unable to upload file", err)
		return nil, err
	}

	// Resolve file url
	fileUrl := s.buildUrl(dest)

	// Compose response
	resp := dto.UploadResp{
		FileName: fileName,
		FileUrl:  fileUrl,
	}

	return &resp, nil
}

func (s *Asset) buildUrl(filePath string) string {
	return s.BaseUrl + "/" + filePath
}

func (s *Asset) GetPublicUrl(assetType int, fileName string) string {
	// If file name is empty, return empty
	if fileName == "" {
		return ""
	}

	// Determine sub dir
	dir, err := s.GetDir(assetType)
	if err != nil {
		return ""
	}

	// Set file path
	filePath := dir + fileName

	return s.buildUrl(filePath)
}

func (s *Asset) GetDir(assetType int) (string, error) {
	dir, ok := api.AssetDirs[assetType]
	if !ok {
		err := s.Errors.New("AST001")
		s.Logger.Errorf("unknown asset type: %d", assetType)
		return "", err
	}
	dir += "/"
	return dir, nil
}

func (s *Asset) NewImageError(err error) error {
	switch err {
	case nhttp.ErrFileTooLarge:
		return s.Errors.New("AST002")
	case nhttp.ErrMimeTypeNotAccepted:
		return s.Errors.New("AST003")
	case http.ErrMissingFile:
		return s.Errors.New("AST004")
	default:
		return nhttp.ErrBadRequest
	}
}
