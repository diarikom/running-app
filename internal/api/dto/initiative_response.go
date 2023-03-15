package dto

import "github.com/diarikom/running-app/running-app-api/internal/api/entity"

type InitiativeResp struct {
	Id                 string                  `json:"id"`
	Name               string                  `json:"name"`
	Description        string                  `json:"description"`
	ImageFiles         InitiativeImageFileResp `json:"image_files"`
	Tags               []string                `json:"tags"`
	ExternalUrls       []entity.ExternalURL    `json:"external_urls"`
	Price              float64                 `json:"price"`
	CurrencyId         int8                    `json:"currency_id"`
	DonationConversion string                  `json:"donation_conversion"`
	CurrencyName       string                  `json:"currency_name"`
	UpdatedAt          int64                   `json:"updated_at"`
	Version            int64                   `json:"version"`
	Headline           string                  `json:"headline"`
}

type InitiativeImageFileResp struct {
	Thumbnail  string `json:"thumbnail"`
	DetailPage string `json:"detail_page"`
}

type DonateReq struct {
	UserId       string `json:"-" validate:"required"`
	InitiativeId string `json:"-" validate:"required"`
	Quantity     int    `json:"qty" validate:"gte=1"`
	Notes        string `json:"notes"`
}

type DonateResp struct {
	Balance       float64 `json:"balance"`
	TotalDonation float64 `json:"total_donation"`
}

type DonationHistoryResp struct {
	Id                     string  `json:"id"`
	InitiativeId           string  `json:"initiative_id"`
	InitiativeName         string  `json:"initiative_name"`
	InitiativeDescription  string  `json:"initiative_description"`
	InitiativeThumbnailURL string  `json:"initiative_thumbnail_url"`
	ItemPrice              float64 `json:"item_price"`
	Quantity               int     `json:"quantity"`
	TotalDonation          float64 `json:"total_donation"`
	StatusId               int8    `json:"status_id"`
	CreatedAt              int64   `json:"created_at"`
	UpdatedAt              int64   `json:"updated_at"`
}
