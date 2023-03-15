package service

import (
	"database/sql"
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
)

type SiteSettingService struct {
	IdGen      *api.SnowflakeGen
	Errors     *api.Errors
	Logger     nlog.Logger
	Repository api.SiteSettingRepository
}

func (s *SiteSettingService) Init(app *api.Api) error {
	s.IdGen = app.Components.Id
	s.Errors = app.Components.Errors
	s.Logger = app.Logger
	s.Repository = NewSiteSettingRepository(app.Datasources.Db, app.Components.Id, app.Components.Errors, app.Logger)

	return nil
}

func (s *SiteSettingService) StaticContent(contentType string) (string, error) {
	// Get static content;
	q, err := s.Repository.StaticContent(contentType)

	// Get exception value;
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		s.Logger.Error("unable to retrieve static content", err)

		return "", err
	}

	return q.Value, nil
}
