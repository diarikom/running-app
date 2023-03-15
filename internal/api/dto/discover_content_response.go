package dto

import (
	"encoding/json"
)

type DiscoverContentResp struct {
	Contents []DiscoverContentItem `json:"contents"`
	Count    int                   `json:"count"`
}

type DiscoverContentItem struct {
	Id          string                `json:"id"`
	Title       string                `json:"title"`
	ContentBody string                `json:"content_body"`
	ExternalUrl json.RawMessage       `json:"external_url"`
	StatusId    int                   `json:"status_id"`
	Sort        int                   `json:"sort"`
	CreatedAt   int64                 `json:"created_at"`
	UpdatedAt   int64                 `json:"updated_at"`
	Version     int                   `json:"version"`
	ModifiedBy  json.RawMessage       `json:"modified_by"`
	ImageFiles  DiscoverImageFileResp `json:"image_files"`
	Headline    string                `json:"headline"`
	Tags        []string              `json:"tags"`
}

type DiscoverImageFileResp struct {
	Thumbnail  string `json:"thumbnail"`
	DetailPage string `json:"detail_page"`
}
