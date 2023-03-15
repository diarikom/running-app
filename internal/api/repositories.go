package api

import (
	"github.com/diarikom/running-app/running-app-api/internal/api/model"
	"time"
)

type DiscoverContentRepository interface {
	CountContents() (total int, err error)
	FindContents(limit int8, skip int64) (result []model.DiscoverContent, err error)
}

type RunRepository interface {
	CountRunSessionHistory(userId string) (total int, err error)
	FindRunSessionHistory(userId string, limit int8, skip int64) (result []model.RunSession, err error)
	InsertRunSession(session model.RunSession) error
	UpdateRunSyncStatus(id, userId string, syncStatus int) error
	SumRunSessionDistance(userID string, start time.Time, end time.Time) (result int, err error)
}

type UserRepository interface {
	DeleteAllSession(userId string) error
	DeleteSessionById(id string) error
	FindAuthByEmail(email string) (*model.UserAuth, error)
	FindAuthById(userId string) (*model.UserAuth, error)
	FindAuthByThirdParty(userId string, providerId int) (*model.UserAuth, error)
	FindProfileById(userId string) (*model.UserProfile, error)
	FindProfileByEmail(userId string) (*model.UserProfile, error)
	FindSessionById(sessionId string) (*model.UserSession, error)
	Insert(userProfile model.UserProfile, userAuth model.UserAuth) error
	InsertWithThirdParty(userProfile model.UserProfile, userAuth model.UserAuth, userAuthThirdParty model.UserAuthThirdParty) error
	InsertSession(userSession model.UserSession) error
	IsExistByEmail(email string) (bool, error)
	IsExistBy3rdPartyAcc(authProviderId int, accessKey string) (bool, error)
	UpdateVerifyEmail(userId string, isVerified bool, timestamp time.Time) error
	UpdatePassword(userId, password string, timestamp time.Time) error
	UpdateProfile(user, newUser model.UserProfile, changes []string) error
	FindEmailById(userId string) (string, error)
	FindProviderRefId(providerId int8, userId string) (string, error)
	FindProviderSubscriptionPlanTypeRefId(providerId, subscriptionPlanTypeId int8) (string, error)
	InsertProviderRefId(mapping model.ProviderUserMapping) error
	InsertSubscription(subscription model.UserSubscription) error
	FindLatestSubscriptionByUser(userId string, providerId int8, now time.Time) (*model.UserSubscription, error)
	FindActiveSubscription(userId string, now time.Time) (*model.UserSubscription, error)
	UpdateSubscriptionStatus(subscription model.UserSubscription) error
	InsertAdminInvitation(invitation model.AdminInvitation) error
	IsUserHasSubscribed(userId string, providerId int8) (bool, error)
}

type AdTagRepository interface {
	GetAdTags() (result []model.AdTag, err error)
}

type MilestoneRepository interface {
	FindById(id string) (*model.Milestone, error)
	FindChallengeById(id string) (*model.Challenge, error)
	Current(now time.Time, status int) (result model.Milestone, err error)
	FindUserChallenge(userID string, challengeID string) (*model.UserChallenge, error)
	InsertUserChallenge(uc model.UserChallenge) error
	UpdateUserChallenge(o, n model.UserChallenge, changes []string) error
	FindCurrentChallenges() ([]model.Challenge, error)
	GetChallengesByStatus(userID string, milestoneID string, status int) ([]model.Challenge, error)
	GetChallengesIn(challengesID string) ([]model.Challenge, error)
	GetUserChallengeByMilestoneId(userID string, milestoneID string, status int8) ([]model.UserChallenge, error)
	FindUnaccomplishedChallengeByUser(userID string, milestoneID string) ([]model.Challenge, error)
}

type CreditRepository interface {
	FindTrxById(trxId string) (*model.UserCreditWalletTrx, error)
	FindWalletById(walletId string) (*model.UserCreditWallet, error)
	FindWalletByTrx(trxId string) (*model.UserCreditWallet, error)
	FindWalletByUser(userId string) (*model.UserCreditWallet, error)
	IsExistTrxRef(walletId, trxId string) (bool, error)
	IsExistWalletByUser(userId string) (bool, error)
	InsertWallet(wallet *model.UserCreditWallet, trx *model.UserCreditWalletTrx) error
	InsertTrx(wallet *model.UserCreditWallet, newTrx *model.UserCreditWalletTrx) error
}

type InitiativeRepository interface {
	FindActive(skip int64, limit int8) ([]model.Initiative, error)
	FindById(id string) (*model.Initiative, error)
	FindDonationByUser(userId string, skip int64, limit int8) ([]model.Donation, error)
	Insert(donation model.Donation, donationLog model.DonationLog) error
	UpdateDonation(oldDonation, newDonation model.Donation, changelog []string) error
}

type SubscriptionPlanRepository interface {
	SubscriptionPlans() ([]model.SubscriptionPlan, error)
}

type SiteSettingRepository interface {
	StaticContent(contentType string) (model.SiteSetting, error)
}
