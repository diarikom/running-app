package dto

import (
	"encoding/json"
	"github.com/stripe/stripe-go/v71"
)

type UserProfileResp struct {
	Id            string `json:"id"`
	FullName      string `json:"full_name"`
	AvatarUrl     string `json:"avatar_url"`
	GenderId      int    `json:"gender_id"`
	DateOfBirth   string `json:"date_of_birth"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	PremiumRunner bool   `json:"premium_runner"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
}

type UserSubscriptionRequestResp struct {
	Stripe json.RawMessage `json:"stripe"`
}

type UserSubscribeResp struct {
	Status                 int8                    `json:"status"`
	SubscriptionPlanTypeId int8                    `json:"subscription_plan_type_id"`
	PeriodStart            int64                   `json:"period_start"`
	PeriodEnd              int64                   `json:"period_end"`
	Stripe                 UserSubscribeStripeResp `json:"stripe"`
}

type UserSubscribeStripeResp struct {
	SubscriptionId            string                          `json:"subscription_id"`
	InvoiceId                 string                          `json:"invoice_id"`
	InvoiceURL                string                          `json:"invoice_url"`
	PaymentIntentStatus       string                          `json:"payment_intent_status"`
	PaymentIntentClientSecret string                          `json:"payment_intent_client_secret"`
	PaymentIntentNextAction   *stripe.PaymentIntentNextAction `json:"payment_intent_next_action"`
}

type SubscriptionVoucherResp struct {
	ProviderRefId  string  `json:"provider_ref_id"`
	ProviderId     int8    `json:"provider_id"`
	Valid          bool    `json:"valid"`
	Value          float64 `json:"value"`
	ValueType      string  `json:"value_type"`
	RecurringMonth int64   `json:"recurring_month"`
}
