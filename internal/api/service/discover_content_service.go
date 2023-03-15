package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/dto"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"strings"
)

type DiscoverContent struct {
	IdGen                     *api.SnowflakeGen
	Errors                    *api.Errors
	Logger                    nlog.Logger
	DiscoverContentRepository api.DiscoverContentRepository
	AssetService              api.AssetService
}

func (r *DiscoverContent) Init(app *api.Api) error {
	r.IdGen = app.Components.Id
	r.Errors = app.Components.Errors
	r.Logger = app.Logger
	r.DiscoverContentRepository = NewDiscoverContentRepository(app.Datasources.Db, app.Logger)
	r.AssetService = app.Services.Asset
	return nil
}

func (r DiscoverContent) GetContents(skip int64, limit int8) (resp *dto.DiscoverContentResp, err error) {
	// Find contents
	contents, err := r.DiscoverContentRepository.FindContents(limit, skip)
	if err != nil {
		r.Logger.Error("unable to find contents", err)
		return nil, err
	}

	// Count total active contents
	count, err := r.DiscoverContentRepository.CountContents()
	if err != nil {
		r.Logger.Error("unable to count contents", err)
		return nil, err
	}

	// Copy model to response
	items := make([]dto.DiscoverContentItem, len(contents))
	for k, v := range contents {
		// Get full image url
		imageUrls := dto.DiscoverImageFileResp{
			Thumbnail:  r.AssetService.GetPublicUrl(api.AssetDiscoverContent, v.ImageFiles.Thumbnail),
			DetailPage: r.AssetService.GetPublicUrl(api.AssetDiscoverContent, v.ImageFiles.DetailPage),
		}

		// Convert tags to array
		var tags []string
		if v.Tags.String == "" {
			tags = []string{}
		} else {
			rawTags := strings.Trim(v.Tags.String, ",")
			tags = strings.Split(rawTags, ",")
		}

		items[k] = dto.DiscoverContentItem{
			Id:          v.Id,
			Title:       v.Title,
			ContentBody: v.ContentBody,
			ExternalUrl: v.ExternalUrl,
			StatusId:    v.StatusId,
			Sort:        v.Sort,
			CreatedAt:   v.CreatedAt.Unix(),
			UpdatedAt:   v.UpdatedAt.Unix(),
			Version:     v.Version,
			ModifiedBy:  v.ModifiedBy,
			ImageFiles:  imageUrls,
			Headline:    v.Headline.String,
			Tags:        tags,
		}
	}

	// Create response result
	resp = &dto.DiscoverContentResp{
		Contents: items,
		Count:    count,
	}

	return resp, nil
}
