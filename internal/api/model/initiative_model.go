package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/diarikom/running-app/running-app-api/internal/api/entity"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
	"github.com/jinzhu/copier"
	"time"
)

type Initiative struct {
	Id                 string              `db:"id" json:"id"`
	OrganizationId     string              `db:"organization_id" json:"organization_id"`
	Name               string              `db:"name" json:"name"`
	Description        string              `db:"description" json:"description"`
	ImageFiles         InitiativeImageFile `db:"image_files" json:"image_files"`
	ExternalUrls       ExternalURLArray    `db:"external_urls" json:"external_urls"`
	Price              float64             `db:"price" json:"price"`
	CurrencyId         int8                `db:"currency_id" json:"currency_id"`
	DonationConversion string              `db:"donation_conversion" json:"donation_conversion"`
	StatusId           int8                `db:"status_id" json:"status_id"`
	Tags               string              `db:"tags" json:"tags"`
	StatDonationCount  int64               `db:"stat_donation_count" json:"stat_donation_count"`
	CreatedAt          time.Time           `db:"created_at" json:"-"`
	UpdatedAt          time.Time           `db:"updated_at" json:"-"`
	Version            int64               `db:"version" json:"version"`
	Headline           sql.NullString      `db:"headline" json:"-"`
}

type ExternalURLArray []entity.ExternalURL

func (e *ExternalURLArray) Scan(src interface{}) error {
	return nsql.ScanJSON(src, e)
}

func (e ExternalURLArray) Value() (driver.Value, error) {
	return json.Marshal(e)
}

type InitiativeImageFile struct {
	Thumbnail  string `json:"thumbnail"`
	DetailPage string `json:"detail_page"`
}

func (i *InitiativeImageFile) Scan(src interface{}) error {
	return nsql.ScanJSON(src, i)
}

func (i InitiativeImageFile) Value() (driver.Value, error) {
	return json.Marshal(i)
}

type InitiativeSnapshot struct {
	*Initiative
	Headline  string `json:"headline"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func NewInitiativeSnapshot(initiative *Initiative) InitiativeSnapshot {
	return InitiativeSnapshot{
		Initiative: initiative,
		Headline:   initiative.Headline.String,
		CreatedAt:  initiative.CreatedAt.Unix(),
		UpdatedAt:  initiative.UpdatedAt.Unix(),
	}
}

func (i *InitiativeSnapshot) Scan(src interface{}) error {
	return nsql.ScanJSON(src, i)
}

func (i InitiativeSnapshot) Value() (driver.Value, error) {
	return json.Marshal(i)
}

type Donation struct {
	Id                 string              `db:"id" diff:"id"`
	InitiativeId       string              `db:"initiative_id" diff:"-"`
	InitiativeSnapshot *InitiativeSnapshot `db:"initiative_snapshot"  diff:"-"`
	UserId             string              `db:"user_id"  diff:"-"`
	UserSnapshot       *UserSnapshot       `db:"user_snapshot"  diff:"-"`
	PaymentMethodId    int8                `db:"payment_method_id"`
	PaymentSnapshot    json.RawMessage     `db:"payment_snapshot"`
	PaymentTrxRef      sql.NullString      `db:"payment_trx_ref"`
	Quantity           int                 `db:"qty"  diff:"-"`
	TotalPrice         float64             `db:"total_price"  diff:"-"`
	CurrencyId         int8                `db:"currency_id"  diff:"-"`
	StatusId           int8                `db:"status_id"`
	Notes              sql.NullString      `db:"notes"`
	CreatedAt          time.Time           `db:"created_at"  diff:"-"`
	UpdatedAt          time.Time           `db:"updated_at" diff:"required"`
	ModifiedBy         *ModifierMeta       `db:"modified_by"  diff:"required"`
	Version            int64               `db:"version"  diff:"required,cc"`
}

func CopyDonation(d *Donation) (*Donation, error) {
	// Duplicate donation
	oldDonation := Donation{
		InitiativeSnapshot: &InitiativeSnapshot{
			Initiative: &Initiative{},
		},
		UserSnapshot: &UserSnapshot{
			UserProfile: &UserProfile{},
		},
		ModifiedBy: &ModifierMeta{},
	}
	err := copier.Copy(&oldDonation, d)
	if err != nil {
		return nil, err
	}
	return &oldDonation, nil
}

type DonationLog struct {
	LogId           string          `db:"log_id"`
	Changelog       nsql.Changelog  `db:"changelog"`
	Id              string          `db:"id"`
	PaymentMethodId sql.NullInt32   `db:"payment_method_id"`
	PaymentSnapshot json.RawMessage `db:"payment_snapshot"`
	PaymentTrxRef   sql.NullString  `db:"payment_trx_ref"`
	StatusId        sql.NullInt32   `db:"status_id"`
	Notes           sql.NullString  `db:"notes"`
	UpdatedAt       time.Time       `db:"updated_at"`
	ModifiedBy      *ModifierMeta   `db:"modified_by"`
	Version         int64           `db:"version"`
}
