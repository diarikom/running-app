package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql/pqx"
	"github.com/lib/pq"
	"github.com/stripe/stripe-go/v71"
	"time"
)

type UserProfile struct {
	Id            string         `db:"id" diff:"id" json:"id"`
	FullName      string         `db:"full_name" json:"full_name"`
	AvatarFile    sql.NullString `db:"avatar_file" json:"-"`
	GenderId      int            `db:"gender_id" diff:"gender" json:"gender_id"`
	DateOfBirth   pqx.Date       `db:"date_of_birth" diff:"dob" json:"date_of_birth"`
	Email         string         `db:"email" diff:"-" json:"email"`
	EmailVerified bool           `db:"email_verified" diff:"-" json:"email_verified"`
	CreatedAt     time.Time      `db:"created_at" diff:"-" json:"-"`
	UpdatedAt     time.Time      `db:"updated_at" diff:"required" json:"-"`
}

type UserAuth struct {
	Id        string    `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	StatusId  int       `db:"status_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserAuthThirdParty struct {
	Id             string    `db:"id"`
	UserId         string    `db:"user_id"`
	AuthProviderId int       `db:"auth_provider_id"`
	AccessKey      string    `db:"access_key"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type UserSession struct {
	Id                    string    `db:"id"`
	UserId                string    `db:"user_id"`
	AuthProviderId        int       `db:"auth_provider_id"`
	DevicePlatformId      int       `db:"device_platform_id"`
	DeviceId              string    `db:"device_id"`
	DeviceManufacturer    string    `db:"device_manufacturer"`
	DeviceModel           string    `db:"device_model"`
	NotificationChannelId int       `db:"notification_channel_id"`
	NotificationToken     string    `db:"notification_token"`
	Signature             string    `db:"signature"`
	ExpiredAt             time.Time `db:"expired_at"`
	CreatedAt             time.Time `db:"created_at"`
	UpdatedAt             time.Time `db:"updated_at"`
}

type UserChallenge struct {
	Id                      string          `db:"id" diff:"id"`
	UserId                  string          `db:"user_id" diff:"-"`
	MilestoneId             string          `db:"milestone_id" diff:"-"`
	MilestoneSnapshot       json.RawMessage `db:"milestone_snapshot" diff:"-"`
	MilestoneVersion        int             `db:"milestone_version" diff:"-"`
	ChallengeId             string          `db:"challenge_id" diff:"-"`
	ChallengeSnapshot       json.RawMessage `db:"challenge_snapshot" diff:"-"`
	ChallengeVersion        int             `db:"challenge_version" diff:"-"`
	ChallengeResultSnapshot json.RawMessage `db:"challenge_result_snapshot" diff:"-"`
	RewardSnapshot          json.RawMessage `db:"reward_snapshot" diff:"-"`
	RewardTypeId            int8            `db:"reward_type_id" diff:"-"`
	RewardRefId             string          `db:"reward_ref_id"`
	RewardValue             float64         `db:"reward_value" diff:"-"`
	Status                  int8            `db:"status"`
	UpdatedAt               time.Time       `db:"updated_at" diff:"required"`
}

type UserCreditWallet struct {
	Id                  string      `db:"id"`
	UserId              string      `db:"user_id"`
	Balance             float64     `db:"balance"`
	BalancePending      float64     `db:"balance_pending"`
	BalanceExpiring     float64     `db:"balance_expiring"`
	BalanceExpiringDate pq.NullTime `db:"balance_expiring_date"`
	CreatedAt           time.Time   `db:"created_at"`
	UpdatedAt           time.Time   `db:"updated_at"`
	Version             int         `db:"version"`
	CurrentVersion      int         `db:"current_version"`
}

type UserCreditWalletTrx struct {
	Id                 string         `db:"id"`
	UserCreditWalletId string         `db:"user_credit_wallet_id"`
	Balance            float64        `db:"balance"`
	BalancePending     float64        `db:"balance_pending"`
	Amount             float64        `db:"amount"`
	TrxEntryTypeId     int8           `db:"trx_entry_type_id"`
	TrxRefId           sql.NullString `db:"trx_ref_id"`
	Notes              sql.NullString `db:"notes"`
	Status             int8           `db:"status"`
	CreatedAt          time.Time      `db:"created_at"`
	ExpiredAt          pq.NullTime    `db:"expired_at"`
	Version            int            `db:"version"`
}

type UserSnapshot struct {
	*UserProfile
	AvatarFile string `json:"avatar_file"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

func (u UserSnapshot) Value() (driver.Value, error) {
	return json.Marshal(u)
}

func (u *UserSnapshot) Scan(src interface{}) error {
	return nsql.ScanJSON(src, u)
}

func NewUserSnapshot(user *UserProfile) UserSnapshot {
	return UserSnapshot{
		UserProfile: user,
		AvatarFile:  user.AvatarFile.String,
		CreatedAt:   user.CreatedAt.Unix(),
		UpdatedAt:   user.UpdatedAt.Unix(),
	}
}

type ProviderUserMapping struct {
	Id          string    `db:"id"`
	UserId      string    `db:"user_id"`
	ProviderId  int8      `db:"provider_id"`
	ProviderRef string    `db:"provider_ref"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type UserSubscription struct {
	Id                      string          `db:"id"`
	UserId                  string          `db:"user_id"`
	PlanTypeId              int8            `db:"plan_type_id"`
	ProviderId              int8            `db:"provider_id"`
	ProviderSubscriptionRef string          `db:"provider_subscription_ref"`
	ProviderOptions         json.RawMessage `db:"provider_options"`
	PeriodStart             time.Time       `db:"period_start"`
	PeriodEnd               time.Time       `db:"period_end"`
	StatusId                int8            `db:"status_id"`
	Metadata                json.RawMessage `db:"metadata"`
	CreatedAt               time.Time       `db:"created_at"`
	UpdatedAt               time.Time       `db:"updated_at"`
	ModifiedBy              ModifierMeta    `db:"modified_by"`
}

type SubscriptionSnapshot struct {
	Stripe *StripeSubscriptionSnapshot `json:"stripe"`
}

type StripeSubscriptionSnapshot struct {
	Subscriptions *stripe.Subscription `json:"subscriptions"`
}

type AdminInvitation struct {
	Id        string    `db:"id"`
	UserId    string    `db:"user_id"`
	Email     string    `db:"email"`
	Token     string    `db:"token"`
	ExpiredAt time.Time `db:"expired_at"`
	CreatedAt time.Time `db:"created_at"`
}
