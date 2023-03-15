package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/dto"
	"github.com/diarikom/running-app/running-app-api/pkg/nhttp"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/diarikom/running-app/running-app-api/pkg/nstr"
	"github.com/gorilla/mux"
	"net/http"
)

func NewUserHandler(app *api.Api) UserHandler {
	return UserHandler{
		UserService:      app.Services.User,
		MilestoneService: app.Services.MilestoneService,
		CreditService:    app.Services.Credit,
		Logger:           app.Logger,
	}
}

type UserHandler struct {
	UserService      api.UserService
	MilestoneService api.MilestoneService
	CreditService    api.CreditService
	Logger           nlog.Logger
}

func (h *UserHandler) PutUpdateProfile(r *http.Request) (*nhttp.Success, error) {
	// Get profile request
	var reqBody dto.UserUpdateProfileReq
	err := nhttp.ParseJSON(&reqBody, r)
	if err != nil {
		return nil, nhttp.ErrBadRequest
	}

	// Set user id
	reqBody.Data.Id = r.Header.Get(nhttp.KeyUserId)

	// Call service
	err = h.UserService.UpdateProfile(reqBody)
	if err != nil {
		return nil, err
	}

	return nhttp.OK(), nil
}

func (h *UserHandler) PostRefreshToken(r *http.Request) (*nhttp.Success, error) {
	// Parse request body
	var reqBody dto.UserLoginReq
	err := nhttp.ParseJSON(&reqBody, r)
	if err != nil {
		return nil, nhttp.ErrBadRequest
	}

	// Get dto
	payload := dto.UserRefreshSession{
		SessionId:  r.Header.Get(nhttp.KeySessionId),
		DeviceInfo: reqBody,
	}

	// Call service
	header, err := h.UserService.RefreshSession(payload)
	if err != nil {
		return nil, err
	}

	// Compose response
	resp := nhttp.OK()
	resp.Header = header

	return resp, nil
}

func (h *UserHandler) PutVerifyEmail(r *http.Request) (*nhttp.Success, error) {
	// Get user id
	userId := r.Header.Get(nhttp.KeyUserId)

	// Call service
	err := h.UserService.UpdateVerifyEmail(userId)
	if err != nil {
		return nil, err
	}

	return nhttp.OK(), nil
}

func (h *UserHandler) GetVerifyEmail(_ *http.Request) (*nhttp.Success, error) {
	return nhttp.OK(), nil
}

func (h *UserHandler) GetResetPassword(_ *http.Request) (*nhttp.Success, error) {
	return nhttp.OK(), nil
}

func (h *UserHandler) PutResetPassword(r *http.Request) (*nhttp.Success, error) {
	// Get user id
	userId := r.Header.Get(nhttp.KeyUserId)

	// Get request body
	var reqBody struct {
		Password string `json:"password"`
	}
	err := nhttp.ParseJSON(&reqBody, r)
	if err != nil {
		return nil, nhttp.ErrBadRequest
	}

	// Call service
	err = h.UserService.ChangePassword(dto.ChangePasswordReq{
		UserId:      userId,
		NewPassword: reqBody.Password,
		Reset:       true,
	})
	if err != nil {
		return nil, err
	}

	return nhttp.OK(), nil
}

func (h *UserHandler) PostResetPassword(r *http.Request) (*nhttp.Success, error) {
	// Get request body
	var reqBody struct {
		Email string `json:"email"`
	}
	err := nhttp.ParseJSON(&reqBody, r)
	if err != nil {
		return nil, nhttp.ErrBadRequest
	}

	// Call service
	err = h.UserService.RequestResetPassword(reqBody.Email)
	if err != nil {
		return nil, err
	}

	return nhttp.OK(), nil
}

func (h *UserHandler) PutChangePassword(r *http.Request) (*nhttp.Success, error) {
	// Get user id
	userId := r.Header.Get(nhttp.KeyUserId)

	// Get password
	var reqBody dto.ChangePasswordReq
	err := nhttp.ParseJSON(&reqBody, r)
	if err != nil {
		return nil, nhttp.ErrBadRequest
	}
	// Set user id
	reqBody.UserId = userId

	// Call service
	err = h.UserService.ChangePassword(reqBody)
	if err != nil {
		return nil, err
	}

	// Return response
	return nhttp.OK(), nil
}

func (h *UserHandler) GetProfile(r *http.Request) (*nhttp.Success, error) {
	// Get user id
	userId := r.Header.Get(nhttp.KeyUserId)

	// Call service
	respBody, err := h.UserService.GetProfile(userId)
	if err != nil {
		return nil, err
	}

	// Return response
	return &nhttp.Success{
		Result: respBody,
	}, nil
}

func (h *UserHandler) DeleteLogout(r *http.Request) (*nhttp.Success, error) {
	// Get user id
	userId := r.Header.Get(nhttp.KeyUserId)

	// Call service
	err := h.UserService.Logout(userId)
	if err != nil {
		return nil, err
	}

	// Return response
	return nhttp.OK(), nil
}

func (h *UserHandler) GetCheckEmail(r *http.Request) (*nhttp.Success, error) {
	// Get email
	email := r.URL.Query().Get("email")
	if email == "" {
		return nil, nhttp.ErrBadRequest
	}

	// Call service
	respBody, err := h.UserService.IsEmailExists(email)
	if err != nil {
		return nil, err
	}

	// Return response
	return &nhttp.Success{
		Result: respBody,
	}, nil
}

func (h *UserHandler) PostRegister(r *http.Request) (*nhttp.Success, error) {
	// Get Request
	var reqBody dto.UserProfileReq
	err := nhttp.ParseJSON(&reqBody, r)
	if err != nil {
		return nil, nhttp.ErrBadRequest
	}

	// Call Service
	err = h.UserService.Register(reqBody)
	if err != nil {
		return nil, err
	}

	// Return response
	return nhttp.OK(), nil
}

func (h *UserHandler) PostLogin(r *http.Request) (*nhttp.Success, error) {
	// Get Request
	var reqBody dto.UserLoginReq
	err := nhttp.ParseJSON(&reqBody, r)
	if err != nil {
		return nil, nhttp.ErrBadRequest
	}

	// Call service
	header, err := h.UserService.Login(reqBody)
	if err != nil {
		return nil, err
	}

	// Compose response
	resp := nhttp.OK()
	resp.Header = header

	// Return response
	return resp, nil
}

func (h *UserHandler) GetClaimCredit(r *http.Request) (*nhttp.Success, error) {
	challengeId := mux.Vars(r)["id"]
	if challengeId == "" {
		return nil, nhttp.ErrBadRequest
	}

	userId := r.Header.Get(nhttp.KeyUserId)
	claimed, err := h.MilestoneService.ClaimCredit(dto.UserChallengeReq{
		UserId:      userId,
		ChallengeId: challengeId,
	})
	if err != nil {
		return nil, err
	}

	resp := &nhttp.Success{
		Result: claimed,
	}
	return resp, nil
}

func (h *UserHandler) GetCreditBalance(r *http.Request) (*nhttp.Success, error) {
	// Get user id from session header
	userId := r.Header.Get(nhttp.KeyUserId)

	// Call service
	resp, err := h.CreditService.GetUserBalance(userId)
	if err != nil {
		return nil, err
	}

	return &nhttp.Success{Result: resp}, nil
}

func (h *UserHandler) GetUserProviderRefId(r *http.Request) (*nhttp.Success, error) {
	// Get user id
	reqBody := dto.UserSubscriptionReq{
		ProviderId: nstr.ParseInt8(mux.Vars(r)["providerId"], api.ProviderStripe),
		UserId:     r.Header.Get(nhttp.KeyUserId),
	}

	// Call service
	resp, err := h.UserService.GetUserProviderRefId(reqBody)
	if err != nil {
		return nil, err
	}

	return &nhttp.Success{Result: resp}, nil
}

func (h *UserHandler) PostUserSubscribe(r *http.Request) (*nhttp.Success, error) {
	var reqBody dto.UserSubscriptionReq
	err := nhttp.ParseJSON(&reqBody, r)
	if err != nil {
		return nil, nhttp.ErrBadRequest
	}

	// Get user id
	reqBody.UserId = r.Header.Get(nhttp.KeyUserId)

	// Call service
	respBody, err := h.UserService.Subscribe(reqBody)
	if err != nil {
		return nil, err
	}

	return &nhttp.Success{Result: respBody}, nil
}

func (h *UserHandler) DeleteUserCancelSubscription(r *http.Request) (*nhttp.Success, error) {
	var reqBody dto.UserSubscriptionReq
	err := nhttp.ParseJSON(&reqBody, r)
	if err != nil {
		return nil, nhttp.ErrBadRequest
	}

	// Get user id
	reqBody.UserId = r.Header.Get(nhttp.KeyUserId)

	// Call service
	err = h.UserService.CancelSubscription(reqBody)
	if err != nil {
		return nil, err
	}

	return nhttp.OK(), nil
}

func (h *UserHandler) GetUserSubscription(r *http.Request) (*nhttp.Success, error) {
	// Get user id
	userId := r.Header.Get(nhttp.KeyUserId)

	// Call service
	respData, err := h.UserService.GetSubscriptionDetail(dto.UserSubscriptionReq{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	return &nhttp.Success{Result: respData}, nil
}

func (h *UserHandler) PostSendAdvertiserActivation(r *http.Request) (success *nhttp.Success, err error) {
	// Get user id from body
	var reqBody dto.AdvertiserActivationReq
	err = nhttp.ParseJSON(&reqBody, r)
	if err != nil {
		return nil, nhttp.ErrBadRequest
	}

	// Validate body
	if reqBody.UserId == "" {
		return nil, nhttp.ErrBadRequest
	}

	// Call service
	err = h.UserService.TriggerSendAdvertiserActivation(reqBody)

	if err != nil {
		return
	}

	return nhttp.OK(), nil
}

func (h *UserHandler) GetValidateVoucher(r *http.Request) (*nhttp.Success, error) {
	// Get query
	query := r.URL.Query()

	voucherCode := query.Get("voucher_code")
	if voucherCode == "" {
		return nil, nhttp.ErrBadRequest
	}

	subscriptionPlanType := nstr.ParseInt8(query.Get("subscription_plan_type_id"),
		api.PremiumRunnerSubscriptionPlanType)

	// Get user id from body
	reqBody := dto.UserSubscriptionReq{
		UserId:                 r.Header.Get(nhttp.KeyUserId),
		ProviderId:             api.ProviderStripe,
		SubscriptionPlanTypeId: subscriptionPlanType,
		VoucherCode:            voucherCode,
	}

	// Call service
	resp, err := h.UserService.ValidateVoucher(reqBody)
	if err != nil {
		return nil, err
	}

	return &nhttp.Success{
		Result: resp,
	}, nil
}
