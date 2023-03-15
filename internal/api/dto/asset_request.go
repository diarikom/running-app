package dto

import "github.com/diarikom/running-app/running-app-api/pkg/nhttp"

type UploadReq struct {
	AssetType int
	File      nhttp.MultipartFile
	DestDir   string
}
