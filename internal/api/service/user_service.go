package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/dto"
	"github.com/diarikom/running-app/running-app-api/internal/api/entity"
	"github.com/diarikom/running-app/running-app-api/internal/api/model"
	"github.com/diarikom/running-app/running-app-api/pkg/nhttp"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/diarikom/running-app/running-app-api/pkg/nmailgun"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql/pqx"
	"github.com/jinzhu/copier"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v71"
	stripeCustomer "github.com/stripe/stripe-go/v71/customer"
	stripeEphemeralKey "github.com/stripe/stripe-go/v71/ephemeralkey"
	stripePaymentMethod "github.com/stripe/stripe-go/v71/paymentmethod"
	stripePromotionCode "github.com/stripe/stripe-go/v71/promotioncode"
	stripeSubscription "github.com/stripe/stripe-go/v71/sub"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

const (
	BcryptCost = 14

	/// RawFormatResetPasswordSubjectSignature represents formatting for raw signature hashing
	/// Passing argument must be in order:
	/// ResetPasswordRequestId, UserId, Email, UpdatedAt Epoch, Salt
	RawFormatResetPasswordSubjectSignature = "ResetPassword-%s-%s-%s-%d-%s"

	/// RawFormatVerifyEmailUserSignature represents formatting for raw signature hashing
	/// Passing argument must be in order:
	/// RequestId, UserId, Email, EmailVerifiedStatus, Salt
	RawFormatVerifyEmailUserSignature = "VerifyEmail-%s-%s-%s-%d-%s"

	// Default unset password
	UnsetPassword = "-"

	// PubSub Topic
	TopicSendAdvertiserActivationEmail = "send_advertiser_activation_email"
)

type User struct {
	IdGen                             *api.SnowflakeGen
	Errors                            *api.Errors
	Facebook                          api.FacebookProviderComponent
	Mailer                            api.MailerComponent
	Logger                            nlog.Logger
	Config                            *viper.Viper
	BaseUrl                           string
	UserAccessLifetime                int
	ResetPasswordTokenLifetime        int
	VerifyEmailTokenLifetime          int
	SignatureSaltResetPasswordSubject string
	SignatureSaltVerifyEmailSubject   string
	AuthService                       api.AuthenticatorService
	AssetService                      api.AssetService
	UserRepository                    api.UserRepository
	PubSub                            *gochannel.GoChannel
}

func (s *User) Init(app *api.Api) error {
	s.IdGen = app.Components.Id
	s.Errors = app.Components.Errors
	s.Facebook = app.Components.Facebook
	s.Mailer = app.Components.Mailer
	s.Logger = app.Logger
	s.Config = app.Config
	s.BaseUrl = app.BaseUrl.String()
	s.UserAccessLifetime = app.Config.GetInt(api.ConfUserAccessLifetime)
	s.ResetPasswordTokenLifetime = app.Config.GetInt(api.ConfResetPasswordTokenLifetime)
	s.VerifyEmailTokenLifetime = app.Config.GetInt(api.ConfVerifyEmailTokenLifetime)
	s.SignatureSaltResetPasswordSubject = app.Config.GetString(api.ConfSignatureSaltResetPasswordSubject)
	s.SignatureSaltVerifyEmailSubject = app.Config.GetString(api.ConfSignatureSaltEmailVerifySubject)
	s.AuthService = app.Services.Auth
	s.AssetService = app.Services.Asset
	s.UserRepository = NewUserRepository(app.Datasources.Db, app.Logger)

	// Set stripe secret key
	stripe.Key = app.Config.GetString(api.ConfStripeSecretKey)

	// Init PubSub
	s.PubSub = gochannel.NewGoChannel(gochannel.Config{}, nil)

	// Subscribe to check challenge achieved
	msg, err := s.PubSub.Subscribe(context.Background(), TopicSendAdvertiserActivationEmail)
	if err != nil {
		return err
	}
	go s.handleSendAdvertiserActivationEmail(msg)

	return nil
}

func (s *User) GetVoucher(voucherCode string) (*stripe.PromotionCode, error) {
	params := &stripe.PromotionCodeListParams{}
	params.Filters.AddFilter("code", "", voucherCode)
	params.Filters.AddFilter("limit", "", "1")
	params.Filters.AddFilter("active", "", "true")

	i := stripePromotionCode.List(params)

	var pc *stripe.PromotionCode

	if i.Next() {
		pc = i.PromotionCode()
	}

	return pc, nil
}

func (s *User) ValidateVoucher(args dto.UserSubscriptionReq) (*dto.SubscriptionVoucherResp, error) {
	// Get voucher from stripe
	voucher, err := s.GetVoucher(args.VoucherCode)
	if err != nil {
		return nil, err
	}

	// If voucher is not
	if voucher == nil {
		s.Logger.Debugf("Voucher not found. Code = %s", args.VoucherCode)
		return nil, s.Errors.New("USR017")
	}

	// Check if user has already subscribe
	hasSubscribed, err := s.UserRepository.IsUserHasSubscribed(args.UserId, args.ProviderId)
	if err != nil {
		s.Logger.Error("failed on retrieving is user has ever been subscribed", err)
		return nil, err
	}
	if hasSubscribed {
		s.Logger.Debugf("User has ever been a subscription. UserId = %s", args.UserId)
		return nil, s.Errors.New("USR017")
	}

	// Validate metadata
	if voucher.Metadata == nil {
		s.Logger.Debugf("Voucher metadata is empty. Voucher Code = %s", args.VoucherCode)
		return nil, s.Errors.New("USR017")
	}

	// Get metadata
	actualMetadata, ok := voucher.Metadata["subscription_plan_type_id"]
	if !ok {
		s.Logger.Debugf("Metadata subscription_plan_type_id is not set. Voucher Code = %s", args.VoucherCode)
	}

	expectedMetadata := fmt.Sprintf("%d", args.SubscriptionPlanTypeId)
	if actualMetadata != expectedMetadata {
		s.Logger.Debugf("Actual metadata is different with expected metadata. Actual = %s, Expected = %s, VoucherCode = %s",
			actualMetadata, expectedMetadata, args.VoucherCode)
		return nil, s.Errors.New("USR017")
	}

	// Compose response
	resp := dto.SubscriptionVoucherResp{
		ProviderId:    api.ProviderStripe,
		ProviderRefId: voucher.ID,
	}

	if c := voucher.Coupon; c != nil {
		if c.PercentOff > 0 {
			resp.Value = c.PercentOff
			resp.ValueType = "percent"
			resp.Valid = true
		} else if c.AmountOff > 0 {
			resp.Value = float64(c.AmountOff)
			resp.ValueType = "amount"
			resp.Valid = true
		}

		if c.Duration == stripe.CouponDurationOnce {
			resp.RecurringMonth = 1
		} else if c.Duration == stripe.CouponDurationRepeating {
			resp.RecurringMonth = c.DurationInMonths
		}
	}

	return &resp, nil
}

func (s *User) handleSendAdvertiserActivationEmail(messages <-chan *message.Message) {
	for msg := range messages {
		s.Logger.Debugf("Received message. Id = %s", msg.UUID)

		// Parse payload
		w := bytes.NewBuffer(msg.Payload)
		dec := gob.NewDecoder(w)
		var payload dto.AdvertiserActivationReq
		err := dec.Decode(&payload)
		if err != nil {
			s.Logger.Error("failed to parse payload", err)
			msg.Ack()
			return
		}

		// Send Advertiser Activation Email
		err = s.SendAdvertiserActivationEmail(payload)
		if err != nil {
			msg.Ack()
			return
		}

		s.Logger.Debug("Done handleSendAdvertiserActivationEmail")
		msg.Ack()
	}
}

func (s *User) SendAdvertiserActivationEmail(args dto.AdvertiserActivationReq) error {
	// Get Profile
	profile, err := s.UserRepository.FindProfileById(args.UserId)
	if err != nil {
		s.Logger.Error("unable to retrieve user profile", err)
		return err
	}

	// Get Active subscription
	subscription, err := s.UserRepository.FindActiveSubscription(args.UserId, time.Now())
	if err != nil {
		s.Logger.Error("unable find active subscription", err)
		return err
	}

	// Check subscription type
	if subscription.PlanTypeId != api.AdvertiserSubscriptionPlanType {
		err = errors.New("user is not an advertiser")
		s.Logger.Errorf("unable to send Advertiser invitation, user is not an advertiser")
		return err
	}

	// Determine activation lifetime
	activationLifetime := s.Config.GetInt64(api.ConfAdvertiserActivationLifetime)
	activationLifetimeStr := ""
	if activationLifetime < 300 {
		// If activation lifetime is less than 5 minutes, then set to 2 weeks
		s.Logger.Debugf("Activation Lifetime is less than 5 minutes. Set to default 2 weeks. activationLifetime = %d", activationLifetime)
		activationLifetime = 2880
		activationLifetimeStr = "2 weeks"
	}

	// Set token expire
	createdAt := time.Now()
	expiredAt := createdAt.Add(time.Duration(activationLifetime) * time.Minute)

	// Generate token
	token, err := gonanoid.Nanoid(64)
	if err != nil {
		s.Logger.Error("failed to generate token", err)
		return err
	}

	// Persist token
	invitation := model.AdminInvitation{
		Id:        s.IdGen.New(),
		UserId:    profile.Id,
		Email:     profile.Email,
		Token:     token,
		ExpiredAt: expiredAt,
		CreatedAt: createdAt,
	}

	err = s.UserRepository.InsertAdminInvitation(invitation)
	if err != nil {
		s.Logger.Error("failed to persist Admin Invitation for Advertiser", err)
		return err
	}

	// Create reset password url
	activationUrl := fmt.Sprintf("%s/advertiser-registration/%s", s.Config.GetString(api.ConfDashboardUrl), token)

	// Send reset password email
	err = s.Mailer.Send(nmailgun.SendOpt{
		Sender:       s.Mailer.GetDefaultSender(),
		Recipients:   []string{profile.Email},
		Subject:      "Running App - Activate Your Advertiser Account",
		TemplateFile: "advertiser_activation.html",
		TemplateData: struct {
			URL           string
			TokenLifetime string
		}{
			URL:           activationUrl,
			TokenLifetime: activationLifetimeStr,
		},
	})
	if err != nil {
		s.Logger.Error("unable to send advertiser activation email", err)
		return err
	}

	return nil
}

func (s *User) GetSubscriptionDetail(args dto.UserSubscriptionReq) (*dto.UserSubscribeResp, error) {
	// Get latest subscription
	userSubscription, err := s.UserRepository.FindActiveSubscription(args.UserId, time.Now())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, s.Errors.New("USR016")
		}
		s.Logger.Error("unable to find latest subscription by user", err)
		return nil, err
	}

	// Check status
	params := stripe.SubscriptionParams{}
	params.AddExpand("latest_invoice.payment_intent")
	subscription, err := stripeSubscription.Get(userSubscription.ProviderSubscriptionRef, &params)
	if err != nil {
		s.Logger.Error("unable to retrieve subscription from provider Stripe", err)
		return nil, err
	}

	// Compose response
	stripeResp := dto.UserSubscribeStripeResp{}
	if invoice := subscription.LatestInvoice; invoice != nil && invoice.PaymentIntent != nil {
		stripeResp.PaymentIntentStatus = string(invoice.PaymentIntent.Status)
		stripeResp.InvoiceId = invoice.ID
		stripeResp.InvoiceURL = invoice.HostedInvoiceURL
		stripeResp.SubscriptionId = subscription.ID

		// If payment intent is not succeed, then show client secret
		if invoice.PaymentIntent.Status != stripe.PaymentIntentStatusSucceeded {
			stripeResp.PaymentIntentClientSecret = invoice.PaymentIntent.ClientSecret
			stripeResp.PaymentIntentNextAction = invoice.PaymentIntent.NextAction
		}
	}

	respBody := dto.UserSubscribeResp{
		Status:                 userSubscription.StatusId,
		SubscriptionPlanTypeId: userSubscription.PlanTypeId,
		PeriodStart:            userSubscription.PeriodStart.Unix(),
		PeriodEnd:              userSubscription.PeriodEnd.Unix(),
		Stripe:                 stripeResp,
	}

	return &respBody, nil
}

func (s *User) TriggerSendAdvertiserActivation(req dto.AdvertiserActivationReq) error {
	// Encode to gob
	var w bytes.Buffer
	enc := gob.NewEncoder(&w)
	err := enc.Encode(req)
	if err != nil {
		s.Logger.Error("failed to encode dto.AdvertiserActivationReq payload", err)
		return err
	}

	// Create message
	msg := message.NewMessage(watermill.NewUUID(), w.Bytes())

	// Publish
	err = s.PubSub.Publish(TopicSendAdvertiserActivationEmail, msg)
	if err != nil {
		s.Logger.Error("failed to publish to "+TopicSendAdvertiserActivationEmail, err)
		return err
	}

	return nil
}

func (s *User) GetStripeCustomer(userId string) (*stripe.Customer, error) {
	// Get user email
	email, err := s.UserRepository.FindEmailById(userId)
	if err != nil {
		s.Logger.Error("unable to retrieve user email", err)
		return nil, err
	}

	// Get user subscription from db
	refId, err := s.UserRepository.FindProviderRefId(api.ProviderStripe, userId)
	if err != nil && err != sql.ErrNoRows {
		s.Logger.Error("unable to find user reference id in provider", err)
		return nil, err
	}

	// If customer id not found, create a new customer id
	var customer *stripe.Customer
	if refId == "" {
		customer, err = stripeCustomer.New(&stripe.CustomerParams{Email: stripe.String(email)})
		if err != nil {
			s.Logger.Error("unable to create Stripe Customer", err)
			return nil, err
		}
		s.Logger.Debugf("Stripe Customer created. ID = %s", customer.ID)

		// Store mapping
		timestamp := time.Now()
		err = s.UserRepository.InsertProviderRefId(model.ProviderUserMapping{
			Id:          s.IdGen.New(),
			UserId:      userId,
			ProviderId:  api.ProviderStripe,
			ProviderRef: customer.ID,
			CreatedAt:   timestamp,
			UpdatedAt:   timestamp,
		})
		if err != nil {
			s.Logger.Error("unable to store user provider ref id", err)
			return nil, err
		}
		s.Logger.Debug("User mapping to Stripe Customer stored.")
	} else {
		customer, err = stripeCustomer.Get(refId, nil)
		if err != nil {
			s.Logger.Error("unable to retrieve Stripe Customer", err)
			return nil, err
		}
		s.Logger.Debug("Stripe Customer mapping is valid")
	}

	return customer, nil
}

func (s *User) GetUserProviderRefId(args dto.UserSubscriptionReq) (*dto.UserSubscriptionRequestResp, error) {
	// Normalize arguments
	if args.ProviderId != api.ProviderStripe {
		args.ProviderId = api.ProviderStripe
	}

	// Get stripe customer
	sc, err := s.GetStripeCustomer(args.UserId)
	if err != nil {
		return nil, err
	}

	// Get ephemeral key
	params := &stripe.EphemeralKeyParams{
		Customer:      stripe.String(sc.ID),
		StripeVersion: stripe.String("2020-03-02"),
	}
	ek, err := stripeEphemeralKey.New(params)
	if err != nil {
		s.Logger.Error("failed to create ephemeral key", err)
		return nil, err
	}

	// Compose response
	resp := dto.UserSubscriptionRequestResp{
		Stripe: ek.RawJSON,
	}
	return &resp, nil
}

func (s *User) getSubscriptionStatus(subscription *stripe.Subscription) int8 {
	// Determine by subscription status
	switch subscription.Status {
	case stripe.SubscriptionStatusActive:
		return api.SubscriptionActive
	case stripe.SubscriptionStatusIncomplete:
		return api.SubscriptionPending
	case stripe.SubscriptionStatusIncompleteExpired:
		return api.SubscriptionFailed
	case stripe.SubscriptionStatusCanceled:
		return api.SubscriptionCanceled
	case stripe.SubscriptionStatusPastDue, stripe.SubscriptionStatusUnpaid:
		return api.SubscriptionInactive
	}

	return api.SubscriptionFailed
}

func (s *User) validateNewSubscription(userId string, providerId int8) (bool, error) {
	// Get user active subscription
	userSubscription, err := s.UserRepository.FindLatestSubscriptionByUser(userId, providerId, time.Now())
	if err != nil && err != sql.ErrNoRows {
		s.Logger.Error("unable to find latest subscription by user", err)
		return false, err
	}

	// User subscription is empty, return false
	if userSubscription == nil {
		return false, nil
	}

	// If user status is not active, return false
	if userSubscription.StatusId != api.SubscriptionActive && userSubscription.StatusId != api.SubscriptionPending {
		return false, nil
	}

	// Validate against provider
	if providerId == api.ProviderStripe {
		subscription, err := stripeSubscription.Get(userSubscription.ProviderSubscriptionRef, nil)
		if err != nil {
			s.Logger.Error("unable to retrieve subscription from provider Stripe", err)
			return false, err
		}

		if subscription.Status != stripe.SubscriptionStatusActive {
			s.Logger.Debugf("Stripe.SubscriptionStatus = %s", subscription.Status)
			return false, err
		}
	}

	return true, nil
}

func (s *User) Subscribe(args dto.UserSubscriptionReq) (*dto.UserSubscribeResp, error) {
	// Validate plan type id
	switch args.SubscriptionPlanTypeId {
	case api.PremiumRunnerSubscriptionPlanType, api.AdvertiserSubscriptionPlanType:
	default:
		s.Logger.Errorf("invalid subscription_plan_type_id = %d", args.SubscriptionPlanTypeId)
		return nil, nhttp.ErrBadRequest
	}

	if args.Stripe == nil {
		s.Logger.Errorf("stripe object cannot be empty")
		return nil, nhttp.ErrBadRequest
	}

	if args.Stripe.PaymentMethodId == "" {
		s.Logger.Errorf("stripe payment method cannot be empty")
		return nil, nhttp.ErrBadRequest
	}

	// Normalize provider id
	if args.ProviderId != api.ProviderStripe {
		args.ProviderId = api.ProviderStripe
	}

	// Get payment method and subscription
	tmp, err := s.GetStripeCustomer(args.UserId)
	if err != nil {
		return nil, err
	}

	// Get all parameters
	customerId := tmp.ID
	paymentMethodId := args.Stripe.PaymentMethodId

	// Validate latest active subscription
	valid, err := s.validateNewSubscription(args.UserId, args.ProviderId)
	if err != nil {
		return nil, err
	}
	if valid {
		return nil, s.Errors.New("USR015")
	}

	// Validate payment method
	pm, err := stripePaymentMethod.Get(paymentMethodId, nil)
	if err != nil {
		// Handle stripe error
		if stripeErr, ok := err.(*stripe.Error); ok && stripeErr.Code == stripe.ErrorCodeResourceMissing {
			s.Logger.Errorf("Payment method is not found for given id. PaymentMethod = %s", args.Stripe.PaymentMethodId)
			return nil, s.Errors.New("STRP001")
		}
		s.Logger.Error("failed to retrieve Stripe payment method", err)
		return nil, err
	}

	// Find subscription plan type mapping
	priceId, err := s.UserRepository.FindProviderSubscriptionPlanTypeRefId(args.ProviderId, args.SubscriptionPlanTypeId)
	if err != nil {
		if err == sql.ErrNoRows {
			s.Logger.Errorf("unsupported Subscription Plan type. PlanType = %d", args.SubscriptionPlanTypeId)
			return nil, s.Errors.New("USR014")
		}
		s.Logger.Error("failed to retrieve subscription plan type reference id", err)
		return nil, err
	}

	// Create subscription
	sp := stripe.SubscriptionParams{
		Customer: stripe.String(customerId),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan: stripe.String(priceId),
			},
		},
		DefaultPaymentMethod: stripe.String(pm.ID),
	}
	sp.AddExpand("latest_invoice.payment_intent")

	// Set voucher
	if args.VoucherProviderRefId != "" {
		sp.PromotionCode = &args.VoucherProviderRefId
	}

	// Subscribe
	subscription, err := stripeSubscription.New(&sp)
	if err != nil {
		// Handle stripe error
		stripeErr, ok := err.(*stripe.Error)
		if ok {
			// Handle promotion code missing error
			if stripeErr.Code == stripe.ErrorCodeResourceMissing && stripeErr.Param == "promotion_code" {
				s.Logger.Error("Invalid voucher submitted", err)
				return nil, s.Errors.New("USR017")
			}

			// Handle invalid payment method
			if stripeErr.Type == stripe.ErrorTypeInvalidRequest && stripeErr.Param == "payment_method" {
				s.Logger.Error("Invalid payment method", err)
				return nil, s.Errors.New("USR019")
			}

			if stripeErr.Msg == "This promotion code cannot be redeemed because the associated customer has prior transactions." {
				s.Logger.Error("Voucher has been redeemed by another user", err)
				return nil, s.Errors.New("USR018")
			}
		}

		s.Logger.Errorf("failed to create Stripe subscription", err)
		return nil, err
	}

	// Store subscription

	// -- Create subscription snapshot for metadata
	metadata, err := json.Marshal(model.SubscriptionSnapshot{
		Stripe: &model.StripeSubscriptionSnapshot{
			Subscriptions: subscription},
	})
	if err != nil {
		s.Logger.Error("failed to create metadata for subscription", err)
		return nil, err
	}
	s.Logger.Debugf("ProviderMetadata = %s", metadata)

	// -- Init timestamp
	timestamp := time.Now()
	subscriptionModel := model.UserSubscription{
		Id:                      s.IdGen.New(),
		UserId:                  args.UserId,
		PlanTypeId:              args.SubscriptionPlanTypeId,
		ProviderId:              args.ProviderId,
		ProviderSubscriptionRef: subscription.ID,
		ProviderOptions:         json.RawMessage("{}"),
		PeriodStart:             time.Unix(subscription.CurrentPeriodStart, 0),
		PeriodEnd:               time.Unix(subscription.CurrentPeriodEnd, 0),
		StatusId:                s.getSubscriptionStatus(subscription),
		Metadata:                metadata,
		CreatedAt:               timestamp,
		UpdatedAt:               timestamp,
		ModifiedBy: model.ModifierMeta{
			Id:   args.UserId,
			Role: api.ModifierUser,
		},
	}

	// -- Persist subscription
	err = s.UserRepository.InsertSubscription(subscriptionModel)
	if err != nil {
		s.Logger.Error("unable to insert user_subscription", err)
		return nil, err
	}

	// If subscription status is active and plan is Advertiser, send activation email
	if subscriptionModel.PlanTypeId == api.AdvertiserSubscriptionPlanType &&
		subscriptionModel.StatusId == api.SubscriptionActive {
		err := s.TriggerSendAdvertiserActivation(dto.AdvertiserActivationReq{UserId: subscriptionModel.UserId})
		if err != nil {
			s.Logger.Error("failed to trigger send advertiser activation email", err)
		}
	}

	// Compose response
	stripeResp := dto.UserSubscribeStripeResp{}
	if invoice := subscription.LatestInvoice; invoice != nil && invoice.PaymentIntent != nil {
		stripeResp.PaymentIntentStatus = string(invoice.PaymentIntent.Status)
		stripeResp.InvoiceId = invoice.ID
		stripeResp.InvoiceURL = invoice.HostedInvoiceURL
		stripeResp.SubscriptionId = subscription.ID

		// If payment intent is not succeed, then show client secret
		if invoice.PaymentIntent.Status != stripe.PaymentIntentStatusSucceeded {
			stripeResp.PaymentIntentClientSecret = invoice.PaymentIntent.ClientSecret
		}
	}

	respBody := dto.UserSubscribeResp{
		Status:                 subscriptionModel.StatusId,
		SubscriptionPlanTypeId: subscriptionModel.PlanTypeId,
		PeriodStart:            subscriptionModel.PeriodStart.Unix(),
		PeriodEnd:              subscriptionModel.PeriodEnd.Unix(),
		Stripe:                 stripeResp,
	}

	return &respBody, nil
}

func (s *User) GetProfileSnapshot(userId string) (*model.UserSnapshot, error) {
	// Get profile
	profile, err := s.UserRepository.FindProfileById(userId)
	if err != nil {
		s.Logger.Error("unable to retrieve user profile", err)
		return nil, err
	}

	// Convert to snapshot
	snapshot := model.NewUserSnapshot(profile)
	return &snapshot, nil
}

func (s *User) CancelSubscription(args dto.UserSubscriptionReq) error {
	// Get latest subscription
	userSubscription, err := s.UserRepository.FindActiveSubscription(args.UserId, time.Now())
	if err != nil {
		if err == sql.ErrNoRows {
			return s.Errors.New("USR016")
		}
		s.Logger.Error("unable to find latest subscription by user", err)
		return err
	}

	// Check status
	subscription, err := stripeSubscription.Get(userSubscription.ProviderSubscriptionRef, nil)
	if err != nil {
		s.Logger.Error("unable to retrieve subscription from provider Stripe", err)
		return err
	}

	if subscription.Status != stripe.SubscriptionStatusActive {
		s.Logger.Debugf("No active subscription. Stripe.SubscriptionStatus = %s", subscription.Status)
		return s.Errors.New("USR015")
	}

	// Cancel subscription
	result, err := stripeSubscription.Cancel(subscription.ID, nil)
	if err != nil {
		s.Logger.Error("unable to cancel Stripe subscription", err)
		return err
	}

	if result.Status != stripe.SubscriptionStatusCanceled {
		s.Logger.Errorf("unexpected cancel status = %s", result.Status)
		return nhttp.ErrInternalServer
	}

	// Update user subscription
	userSubscription.StatusId = api.SubscriptionCanceled
	userSubscription.UpdatedAt = time.Now()
	userSubscription.ModifiedBy = model.ModifierMeta{
		Id:   args.UserId,
		Role: api.ModifierUser,
	}
	err = s.UserRepository.UpdateSubscriptionStatus(*userSubscription)

	return nil
}

func (s *User) RefreshSession(payload dto.UserRefreshSession) (map[string]string, error) {
	// Get device info
	d := payload.DeviceInfo

	// Validate input
	if d.DevicePlatformId == 0 ||
		d.DeviceId == "" ||
		d.DeviceModel == "" ||
		d.DeviceManufacturer == "" ||
		d.NotificationChannel == 0 ||
		d.NotificationToken == "" {
		return nil, nhttp.ErrBadRequest
	}

	// Get existing session
	session, err := s.UserRepository.FindSessionById(payload.SessionId)
	if err != nil {
		s.Logger.Error("unable to find session by id", err)
		return nil, err
	}

	// Validate device
	if d.DeviceId != session.DeviceId ||
		d.DevicePlatformId != session.DevicePlatformId ||
		d.DeviceManufacturer != session.DeviceManufacturer ||
		d.DeviceModel != session.DeviceModel {
		err := s.Errors.New("USR013")
		return nil, err
	}

	// Delete current session
	err = s.UserRepository.DeleteSessionById(session.Id)
	if err != nil {
		s.Logger.Error("unable to delete session by id", err)
		return nil, err
	}

	// Set auth provider
	d.AuthProviderId = session.AuthProviderId

	// Create session request body
	token, err := s.NewSession(session.UserId, d)
	if err != nil {
		return nil, err
	}

	// Return header
	header := map[string]string{
		api.AccessTokenKey:    token.Token,
		api.AccessTokenExpKey: strconv.FormatInt(token.ExpiredAt, 10),
	}

	return header, err
}

func (s *User) UpdateProfile(req dto.UserUpdateProfileReq) error {
	// Get user
	user, err := s.UserRepository.FindProfileById(req.Data.Id)
	if err != nil {
		return err
	}

	// Copy user
	var newUser model.UserProfile
	err = copier.Copy(&newUser, user)
	if err != nil {
		s.Logger.Error("unable to copy user profile instance", err)
	}

	// Set values
	newUser.FullName = req.Data.FullName
	newUser.GenderId = req.Data.Gender
	newUser.DateOfBirth = pqx.ParseDate(pqx.DateOpt{Input: req.Data.DOB})
	newUser.AvatarFile = nsql.NullString(req.Data.AvatarFile.FileName)
	newUser.UpdatedAt = time.Now()

	// Persist updates
	err = s.UserRepository.UpdateProfile(*user, newUser, req.Changes)
	if err != nil {
		s.Logger.Error("unable to persist user profile update", err)
		return err
	}

	return nil
}

func (s *User) UpdateVerifyEmail(userId string) error {
	// Persist email verify update
	err := s.UserRepository.UpdateVerifyEmail(userId, true, time.Now())
	if err != nil {
		s.Logger.Error("unable to persist status to email verified", err)
		return err
	}

	return err
}

func (s *User) ValidateResetPasswordSignature(session *dto.ResetPasswordSession) (string, error) {
	// Check email
	if session.Email == "" {
		return "", nhttp.ErrBadRequest
	}

	// Check if email exist
	user, err := s.UserRepository.FindAuthByEmail(session.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", s.Errors.New("USR008")
		}
		s.Logger.Error("unable to check if email exists", err)
		return "", err
	}

	// Generate signature
	signature, err := s.AuthService.SignMd5(dto.SignatureReq{
		Format: RawFormatResetPasswordSubjectSignature,
		Args: []interface{}{
			session.RequestId,
			user.Id,
			session.Email,
			user.UpdatedAt.Unix(),
			s.SignatureSaltResetPasswordSubject,
		},
	})

	// Compare signature
	if signature != session.UserSignature {
		return "", s.Errors.New("USR010")
	}

	return user.Id, nil
}

func (s *User) ValidateVerifyEmailSignature(session *dto.VerifyEmailSession) (string, error) {
	// Check email
	if session.Email == "" {
		return "", nhttp.ErrBadRequest
	}

	// Check if email exist
	user, err := s.UserRepository.FindProfileByEmail(session.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", s.Errors.New("USR008")
		}
		s.Logger.Error("unable to check if email exists", err)
		return "", err
	}

	// Generate signature
	signature, err := s.AuthService.SignMd5(dto.SignatureReq{
		Format: RawFormatVerifyEmailUserSignature,
		Args: []interface{}{
			session.RequestId,
			user.Id,
			user.Email,
			user.EmailVerified,
			s.SignatureSaltVerifyEmailSubject,
		},
	})

	// Compare signature
	if signature != session.UserSignature {
		return "", s.Errors.New("USR010")
	}

	return user.Id, nil
}

func (s *User) RequestResetPassword(email string) error {
	// Check email
	if email == "" {
		return nhttp.ErrBadRequest
	}

	// Check if email exist
	user, err := s.UserRepository.FindAuthByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			return s.Errors.New("USR008")
		}
		s.Logger.Error("unable to check if email exists", err)
		return err
	}

	// Validate user status
	if user.StatusId == api.UserSuspended {
		return s.Errors.New("USR002")
	}

	// Create request id as session
	reqId := s.IdGen.New()

	// Sign subject
	signature, err := s.AuthService.SignMd5(dto.SignatureReq{
		Format: RawFormatResetPasswordSubjectSignature,
		Args: []interface{}{
			reqId,
			user.Id,
			email,
			user.UpdatedAt.Unix(),
			s.SignatureSaltResetPasswordSubject,
		},
	})

	// Create purpose token
	token, err := s.AuthService.NewOneTimeToken(dto.JWTOptReq{
		Subject:   email,
		SessionId: reqId,
		Lifetime:  s.ResetPasswordTokenLifetime,
		Purpose:   api.JWTPurposeResetPassword,
		Extras: map[string]string{
			api.UserSignatureKey: signature,
		}})
	if err != nil {
		return err
	}

	// Create reset password url
	resetUrl := fmt.Sprintf("%s/users/reset-password?t=%s", s.BaseUrl, token.Token)

	// Send reset password email
	err = s.Mailer.Send(nmailgun.SendOpt{
		Sender:       s.Mailer.GetDefaultSender(),
		Recipients:   []string{email},
		Subject:      "Running App - Reset Password",
		TemplateFile: "reset_password_request.html",
		TemplateData: struct {
			URL           string
			TokenLifetime string
		}{
			URL:           resetUrl,
			TokenLifetime: "2 hours",
		},
	})
	if err != nil {
		s.Logger.Error("unable to send reset password email", err)
		return err
	}

	return nil
}

func (s *User) validatePassword(input string) bool {
	if input == "" || len(input) < 8 {
		return false
	}
	return true
}

func (s *User) hashPassword(input string) (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(input), BcryptCost)
	return string(password), err
}

func (s *User) ChangePassword(req dto.ChangePasswordReq) error {
	// Get user auth by id
	auth, err := s.UserRepository.FindAuthById(req.UserId)
	if err != nil {
		s.Logger.Error("unable to find user auth by id", err)
		return err
	}

	// If not reset and existing password is set, check existing password
	if !req.Reset && auth.Password != UnsetPassword {
		err = bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(req.ExistingPassword))
		if err != nil {
			return s.Errors.New("USR007")
		}
	}

	// validate password input
	if !s.validatePassword(req.NewPassword) {
		return s.Errors.New("USR005")
	}

	// If old password and new password are the same, return error
	if req.NewPassword == req.ExistingPassword {
		return s.Errors.New("USR006")
	}

	// Hash password
	password, err := s.hashPassword(req.NewPassword)
	if err != nil {
		s.Logger.Error("unable to hash password", err)
		return err
	}

	// Persist new password
	err = s.UserRepository.UpdatePassword(req.UserId, password, time.Now())
	if err != nil {
		s.Logger.Error("unable to persist password update", err)
		return err
	}

	return nil
}

func (s *User) IsPremiumRunner(userId string) (bool, error) {
	// Get user active subscription
	userSubscription, err := s.UserRepository.FindActiveSubscription(userId, time.Now())
	if err != nil {
		// If empty row, then return false
		if err == sql.ErrNoRows {
			return false, err
		}
		s.Logger.Error("unable to find latest subscription by user", err)
		return false, err
	}

	premium := userSubscription.StatusId == api.SubscriptionActive

	return premium, nil
}

func (s *User) GetProfile(userId string) (*dto.UserProfileResp, error) {
	// Get Profile
	profile, err := s.UserRepository.FindProfileById(userId)
	if err != nil {
		s.Logger.Error("unable to retrieve user profile", err)
		return nil, err
	}

	// Get date of birth
	var dob string
	if profile.DateOfBirth.Valid {
		dob = profile.DateOfBirth.Time.Format("2006-01-02")
	}

	isPremium, err := s.IsPremiumRunner(userId)
	if err != nil {
		return nil, err
	}

	// Resolve avatar file url
	avatarUrl := s.AssetService.GetPublicUrl(api.AssetAvatarProfile, profile.AvatarFile.String)

	// Compose response
	respBody := dto.UserProfileResp{
		Id:            profile.Id,
		FullName:      profile.FullName,
		AvatarUrl:     avatarUrl,
		GenderId:      profile.GenderId,
		DateOfBirth:   dob,
		Email:         profile.Email,
		EmailVerified: profile.EmailVerified,
		PremiumRunner: isPremium,
		CreatedAt:     profile.CreatedAt.Unix(),
		UpdatedAt:     profile.UpdatedAt.Unix(),
	}

	return &respBody, nil
}

func (s *User) Logout(userId string) error {
	// Clear user session
	err := s.UserRepository.DeleteAllSession(userId)
	return err
}

func (s *User) ValidateSession(sessionId, userId string) error {
	session, err := s.UserRepository.FindSessionById(sessionId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nhttp.ErrUnauthorized
		}

		s.Logger.Error("unable to find session", err)
		return err
	}

	// Check user
	if userId != session.UserId {
		s.Logger.Debugf("invalid user id. Expected: %s, Actual: %s", session.UserId, userId)
		return nhttp.ErrUnauthorized
	}

	// Check expiry
	if session.ExpiredAt.Unix() < time.Now().Unix() {
		s.Logger.Debugf("token has expired. UserId: %s", userId)
		return s.Errors.New("USR003")
	}

	return nil
}

func (s *User) NewSession(userId string, req dto.UserLoginReq) (*entity.AccessToken, error) {
	// New session id
	sessionId := s.IdGen.New()

	// Create token
	token, err := s.AuthService.NewAccessToken(dto.JWTOptReq{
		Subject:   userId,
		SessionId: sessionId,
		Lifetime:  s.UserAccessLifetime,
	})
	if err != nil {
		return nil, err
	}

	// Init timestamp
	timestamp := time.Now()

	// Generate signature
	signatureRaw := fmt.Sprintf("%d-%s-%s-%s", req.DevicePlatformId, req.DeviceId, req.DeviceModel, req.DeviceManufacturer)
	hasher := md5.New()
	_, err = hasher.Write([]byte(signatureRaw))
	if err != nil {
		return nil, err
	}
	signature := hex.EncodeToString(hasher.Sum(nil))

	// TODO: Create refresh token

	// Create session
	session := model.UserSession{
		Id:                    sessionId,
		UserId:                userId,
		AuthProviderId:        req.AuthProviderId,
		DevicePlatformId:      req.DevicePlatformId,
		DeviceId:              req.DeviceId,
		DeviceManufacturer:    req.DeviceManufacturer,
		DeviceModel:           req.DeviceModel,
		NotificationChannelId: req.NotificationChannel,
		NotificationToken:     req.NotificationToken,
		Signature:             signature,
		ExpiredAt:             time.Unix(token.ExpiredAt, 0),
		CreatedAt:             timestamp,
		UpdatedAt:             timestamp,
	}

	// Persist session
	err = s.UserRepository.InsertSession(session)
	if err != nil {
		s.Logger.Error("unable to persist new session", err)
		return nil, err
	}

	return token, nil
}

func (s *User) Login(req dto.UserLoginReq) (map[string]string, error) {
	// Validate required fields in general
	if req.DevicePlatformId == 0 ||
		req.DeviceId == "" ||
		req.DeviceModel == "" ||
		req.DeviceManufacturer == "" ||
		req.NotificationChannel == 0 ||
		req.NotificationToken == "" {
		return nil, nhttp.ErrBadRequest
	}

	// Switch login
	switch req.AuthProviderId {
	case api.FacebookAuthProvider:
		return s.LoginByFacebook(req)
	default:
		return s.LoginByEmail(req)
	}
}

func (s *User) LoginByEmail(req dto.UserLoginReq) (map[string]string, error) {
	// Validate required fields for login by email
	if req.Email == "" ||
		req.Password == "" {
		return nil, nhttp.ErrBadRequest
	}

	// Get user auth by email
	auth, err := s.UserRepository.FindAuthByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nhttp.ErrUnauthorized
		}

		s.Logger.Error("unable to retrieve UserAuth", err)
		return nil, err
	}

	// Check if password unset
	if auth.Password == UnsetPassword {
		return nil, s.Errors.New("USR011")
	}

	// Validate password
	err = bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(req.Password))
	if err != nil {
		return nil, s.Errors.New("USR007")
	}

	// Validate user status
	if auth.StatusId == api.UserSuspended {
		return nil, s.Errors.New("USR002")
	}

	// Set auth provider
	req.AuthProviderId = api.AppAuthProvider

	// Clear user session
	err = s.UserRepository.DeleteAllSession(auth.Id)
	if err != nil {
		s.Logger.Error("unable to clear user session", err)
		return nil, err
	}

	// Generate Session
	token, err := s.NewSession(auth.Id, req)
	if err != nil {
		return nil, err
	}

	// Return header
	header := map[string]string{
		api.AccessTokenKey:    token.Token,
		api.AccessTokenExpKey: strconv.FormatInt(token.ExpiredAt, 10),
	}

	return header, nil
}

func (s *User) LoginByFacebook(req dto.UserLoginReq) (map[string]string, error) {
	// Validate required fields
	if req.ThirdPartyToken == "" {
		return nil, nhttp.ErrBadRequest
	}

	facebookUser, err := s.Facebook.InspectToken(req.ThirdPartyToken)
	if err != nil {
		s.Logger.Error("unable to inspect facebook user token", err)
		return nil, nhttp.ErrUnauthorized
	}

	// Get user auth by facebook id
	auth, err := s.UserRepository.FindAuthByThirdParty(facebookUser.UserId, api.FacebookAuthProvider)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, s.Errors.New("USR008")
		}

		s.Logger.Error("unable to retrieve UserAuth", err)
		return nil, err
	}

	// Validate user status
	if auth.StatusId == api.UserSuspended {
		return nil, s.Errors.New("USR002")
	}

	// Clear user session
	err = s.UserRepository.DeleteAllSession(auth.Id)
	if err != nil {
		s.Logger.Error("unable to clear user session", err)
		return nil, err
	}

	// Generate Session
	token, err := s.NewSession(auth.Id, req)
	if err != nil {
		return nil, err
	}

	// Return header
	header := map[string]string{
		api.AccessTokenKey:    token.Token,
		api.AccessTokenExpKey: strconv.FormatInt(token.ExpiredAt, 10),
	}

	return header, nil
}

func (s *User) Register(req dto.UserProfileReq) error {
	// Validate required fields
	if req.FullName == "" ||
		req.Email == "" ||
		req.Gender < api.GenderMale ||
		req.Gender > api.GenderFemale {
		return nhttp.ErrBadRequest
	}

	// Validate email is exist
	isExist, err := s.UserRepository.IsExistByEmail(req.Email)
	if err != nil {
		s.Logger.Error("unable to check email is exist", err)
		return err
	}
	if isExist {
		return s.Errors.New("USR001")
	}

	// Switch auth provider
	switch req.AuthProviderId {
	case api.FacebookAuthProvider:
		return s.RegisterByFacebook(req)
	default:
		return s.RegisterByEmail(req)
	}
}

func (s *User) RegisterByEmail(req dto.UserProfileReq) error {
	// Validate password
	if !s.validatePassword(req.Password) {
		return s.Errors.New("USR005")
	}

	// Encrypt password with bcrypt
	password, err := s.hashPassword(req.Password)
	if err != nil {
		return err
	}

	// Init timestamp
	timestamp := time.Now()

	// Create user profile
	userProfile := model.UserProfile{
		Id:          s.IdGen.New(),
		FullName:    req.FullName,
		GenderId:    req.Gender,
		DateOfBirth: pqx.ParseDate(pqx.DateOpt{Input: req.DOB}),
		Email:       req.Email,
		CreatedAt:   timestamp,
		UpdatedAt:   timestamp,
	}

	// Create user auth
	userAuth := model.UserAuth{
		Id:        userProfile.Id,
		Username:  userProfile.Email,
		Password:  password,
		StatusId:  api.UserActive,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	// Persist user
	err = s.UserRepository.Insert(userProfile, userAuth)
	if err != nil {
		s.Logger.Error("unable to persist user", err)
		return err
	}

	// Send email verification
	_ = s.SendEmailVerification(userProfile)

	return nil
}

func (s *User) RegisterByFacebook(req dto.UserProfileReq) error {
	// Validate token input
	if req.ThirdPartyToken == "" {
		return nhttp.ErrBadRequest
	}

	// Validate token to Facebook
	fbUser, err := s.Facebook.InspectToken(req.ThirdPartyToken)
	if err != nil {
		s.Logger.Error("unable to inspect token", err)
		return nhttp.ErrUnauthorized
	}

	// Validate existing facebook id
	isExist, err := s.UserRepository.IsExistBy3rdPartyAcc(api.FacebookAuthProvider, fbUser.UserId)
	if err != nil && err != sql.ErrNoRows {
		s.Logger.Error("unable to retrieve user by third party", err)
		return err
	}

	if isExist {
		return s.Errors.New("USR012")
	}

	// Init timestamp
	timestamp := time.Now()

	// Create user profile
	userProfile := model.UserProfile{
		Id:          s.IdGen.New(),
		FullName:    req.FullName,
		GenderId:    req.Gender,
		DateOfBirth: pqx.ParseDate(pqx.DateOpt{Input: req.DOB}),
		Email:       req.Email,
		CreatedAt:   timestamp,
		UpdatedAt:   timestamp,
	}

	// Create user auth
	userAuth := model.UserAuth{
		Id:        userProfile.Id,
		Username:  userProfile.Email,
		Password:  UnsetPassword,
		StatusId:  api.UserActive,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	userAuthThirdParty := model.UserAuthThirdParty{
		Id:             s.IdGen.New(),
		UserId:         userProfile.Id,
		AuthProviderId: req.AuthProviderId,
		AccessKey:      fbUser.UserId,
		CreatedAt:      timestamp,
		UpdatedAt:      timestamp,
	}

	// Persist user with third party
	err = s.UserRepository.InsertWithThirdParty(userProfile, userAuth, userAuthThirdParty)
	if err != nil {
		s.Logger.Error("unable to persist user by 3rd party", err)
		return err
	}

	return nil
}

func (s *User) SendEmailVerification(profile model.UserProfile) error {
	// Create request id as session
	reqId := s.IdGen.New()

	// Sign user
	signature, err := s.AuthService.SignMd5(dto.SignatureReq{
		Format: RawFormatVerifyEmailUserSignature,
		Args: []interface{}{
			reqId,
			profile.Id,
			profile.Email,
			profile.EmailVerified,
			s.SignatureSaltVerifyEmailSubject,
		},
	})
	if err != nil {
		s.Logger.Error("unable to sign subject for verify email", err)
		return err
	}

	// Create purpose token
	token, err := s.AuthService.NewOneTimeToken(dto.JWTOptReq{
		Subject:   profile.Email,
		SessionId: reqId,
		Lifetime:  s.VerifyEmailTokenLifetime,
		Purpose:   api.JWTPurposeVerifyEmail,
		Extras: map[string]string{
			api.UserSignatureKey: signature,
		}})
	if err != nil {
		return err
	}

	// Create verify email url
	verifyEmailUrl := fmt.Sprintf("%s/users/verify-email?t=%s", s.BaseUrl, token.Token)

	// Send email verification
	err = s.Mailer.Send(nmailgun.SendOpt{
		Sender:       s.Mailer.GetDefaultSender(),
		Recipients:   []string{profile.Email},
		Subject:      "Running App - Email Verification",
		TemplateFile: "verify_email.html",
		TemplateData: struct {
			URL string
		}{
			URL: verifyEmailUrl,
		},
	})
	if err != nil {
		s.Logger.Error("unable to send email verification to user", err)
		return err
	}

	return nil
}

func (s *User) IsEmailExists(email string) (interface{}, error) {
	isExist, err := s.UserRepository.IsExistByEmail(email)
	if err != nil {
		return nil, err
	}

	// Compose response
	resp := map[string]bool{
		"is_exist": isExist,
	}

	return resp, nil
}
