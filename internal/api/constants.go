package api

const (
	GenderMale   = 1
	GenderFemale = 2

	UserActive    = 1
	UserSuspended = 2

	RunSummaryStored = 1
	RunDetailsStored = 2

	AccessTokenKey    = "X-Access-Token"
	AccessTokenExpKey = "X-Access-Token-Expiry"

	JWTAudienceUser = "RunningApp.User"
	JWTAudienceApp  = "RunningApp.App"

	UserSignatureKey = "user_signature"
)

const (
	SkipDefault  = 0
	LimitDefault = 10
)

const (
	JWTUser = iota
	JWTApp
	JWTPurposeResetPassword
	JWTPurposeVerifyEmail
)

const (
	AppAuthProvider = 1
	// GoogleAuthProvider
	FacebookAuthProvider = 3
)

const (
	_ = iota
	AssetAvatarProfile
	AssetInitiative
	AssetDiscoverContent
)

var AssetDirs = map[int]string{
	AssetAvatarProfile:   "avatars",
	AssetInitiative:      "initiatives",
	AssetDiscoverContent: "discover-contents",
}

const (
	// MilestoneDraft = iota + 1
	MilestoneStart = 2
	// MilestoneEnd
)

const (
	ChallengeDraft = iota + 1
	ChallengeStart
	ChallengeAchieved
	ChallengeEnd
	ChallengeRewardClaimed
)

const (
	TrxPending = iota + 1
	TrxSuccess
	TrxFailed
	TrxExpired
)

const (
	InitEntryType = iota + 1
	Debit
	Credit
)

const (
	CreditReward = iota + 1
)

const (
	UserAccumulatedRunDistanceFact = "user_acc_run_distance"
)

const (
	_ = iota
	CurrencyDecibel
)

var CurrencyName = map[int8]string{
	CurrencyDecibel: "dB",
}

const (
	ActiveInitiative = 2
)

const (
	PaymentMethodCredit = 1
)

const (
	DonationCreated        = 1
	DonationPaymentPending = 3
	DonationPaymentOK      = 4
)

const (
	ModifierUser = "USER"
)

const (
	ProviderStripe int8 = 1
)

const (
	_ int8 = iota
	PremiumRunnerSubscriptionPlanType
	AdvertiserSubscriptionPlanType
)

const (
	_ int8 = iota
	SubscriptionActive
	SubscriptionPending
	SubscriptionCanceled
	SubscriptionFailed
	SubscriptionInactive
)
