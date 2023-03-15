package api

import (
	"github.com/diarikom/running-app/running-app-api/internal/api/dto"
	"github.com/diarikom/running-app/running-app-api/internal/api/entity"
	"github.com/diarikom/running-app/running-app-api/internal/api/model"
	"github.com/diarikom/running-app/running-app-api/pkg/nhttp"
)

type ServiceInitiator interface {
	Init(app *Api) error
}

type AssetService interface {
	GetPublicUrl(assetType int, fileName string) string
	GetUploadRule(assetType int) (*nhttp.UploadRule, error)
	UploadFile(req dto.UploadReq) (*dto.UploadResp, error)
}

type AuthenticatorService interface {
	NewAccessToken(req dto.JWTOptReq) (*entity.AccessToken, error)
	NewOneTimeToken(req dto.JWTOptReq) (*entity.AccessToken, error)
	SignMd5(req dto.SignatureReq) (string, error)
	ValidateUserAccess(bearer string) (sessionId, userId string, err error)
	ValidateResetPasswordToken(token string) (*dto.ResetPasswordSession, error)
	ValidateVerifyEmailToken(token string) (*dto.VerifyEmailSession, error)
	ValidateClient(secret string) (err error)
	ValidateClientDashboard(secret string) (err error)
}

type DiscoverContentService interface {
	GetContents(skip int64, limit int8) (resp *dto.DiscoverContentResp, err error)
}

type RunService interface {
	GetRunSessionHistory(userId string, skip int64, limit int8) (resp *dto.RunSessionHistoryResp, err error)
	NewRunSession(userId string, req *dto.RunSessionReq) error
	UpdateRunSyncStatus(id, userId string, status int) error
}

type UserService interface {
	ChangePassword(req dto.ChangePasswordReq) error
	GetProfile(userId string) (*dto.UserProfileResp, error)
	GetProfileSnapshot(userId string) (*model.UserSnapshot, error)
	IsEmailExists(email string) (interface{}, error)
	Login(req dto.UserLoginReq) (map[string]string, error)
	LoginByFacebook(req dto.UserLoginReq) (map[string]string, error)
	Logout(userId string) error
	Register(req dto.UserProfileReq) error
	RefreshSession(req dto.UserRefreshSession) (map[string]string, error)
	RequestResetPassword(email string) error
	UpdateProfile(req dto.UserUpdateProfileReq) error
	UpdateVerifyEmail(userId string) error
	ValidateSession(sessionId, userId string) error
	ValidateResetPasswordSignature(session *dto.ResetPasswordSession) (string, error)
	ValidateVerifyEmailSignature(session *dto.VerifyEmailSession) (string, error)
	GetUserProviderRefId(req dto.UserSubscriptionReq) (*dto.UserSubscriptionRequestResp, error)
	Subscribe(req dto.UserSubscriptionReq) (*dto.UserSubscribeResp, error)
	CancelSubscription(req dto.UserSubscriptionReq) error
	GetSubscriptionDetail(args dto.UserSubscriptionReq) (*dto.UserSubscribeResp, error)
	IsPremiumRunner(userId string) (bool, error)
	TriggerSendAdvertiserActivation(req dto.AdvertiserActivationReq) error
	ValidateVoucher(args dto.UserSubscriptionReq) (*dto.SubscriptionVoucherResp, error)
}

type AdTagService interface {
	GetAdTags() (resp *dto.AdTagResp, err error)
}

type MilestoneService interface {
	ClaimCredit(opt dto.UserChallengeReq) (*dto.ChallengeRewardClaimResp, error)
	Current(userID string) (resp *dto.MilestoneResp, err error)
	CheckChallengeAchieve(req dto.UserChallengeReq) (*dto.MilestoneAchievementResp, error)
	LoadMilestone()
	TriggerCheckChallengeAchieved(req dto.UserChallengeReq) error
	UpdateUserChallengeAchievement(opt dto.UserChallengeReq) error
}

type CreditService interface {
	Charge(opt dto.CreditChargeOpt) (string, error)
	CheckChargeAmount(opt dto.CreditChargeOpt) (int, error)
	GetUserWallet(userId string) (*model.UserCreditWallet, error)
	GetUserBalance(userId string) (*dto.UserCreditBalanceResp, error)
	InsertPendingTrx(opt dto.CreditTrxOpt) (string, error)
	SettlePendingTrx(opt dto.CreditSettleOpt) error
}

type InitiativeService interface {
	Donate(opt dto.DonateReq) (*dto.DonateResp, error)
	List(opt dto.PageReq) ([]dto.InitiativeResp, error)
	ListUserDonation(opt dto.UserResourcesReq) ([]dto.DonationHistoryResp, error)
}

type SubscriptionPlanService interface {
	List() ([]dto.SubscriptionPlanResp, error)
}

type SiteSettingService interface {
	StaticContent(typeContent string) (string, error)
}
