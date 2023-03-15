package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/dto"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
)

type AdTagService struct {
	IdGen           *api.SnowflakeGen
	Errors          *api.Errors
	Logger          nlog.Logger
	AdTagRepository api.AdTagRepository
}

func (r *AdTagService) Init(app *api.Api) error {
	r.IdGen = app.Components.Id
	r.Errors = app.Components.Errors
	r.Logger = app.Logger
	r.AdTagRepository = NewAdTagRepository(app.Datasources.Db, app.Logger)
	return nil
}

func (r AdTagService) GetAdTags() (resp *dto.AdTagResp, err error) {
	result, err := r.AdTagRepository.GetAdTags()

	tags := make([]string, len(result))
	for k := range result {
		tags[k] = result[k].Name
	}
	resp = &dto.AdTagResp{
		Tags: tags,
	}

	return resp, nil
}
