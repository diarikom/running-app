package dto

type UserUpdateProfileReq struct {
	Data    UserProfileReq `json:"data"`
	Changes []string       `json:"changes"`
}

type UserProfileReq struct {
	Id              string     `json:"id"`
	Email           string     `json:"email"`
	Password        string     `json:"password"`
	FullName        string     `json:"full_name"`
	DOB             string     `json:"dob"`
	Gender          int        `json:"gender"`
	AvatarFile      UploadResp `json:"avatar_file"`
	AuthProviderId  int        `json:"auth_provider_id"`
	ThirdPartyToken string     `json:"third_party_token"`
}

type UserLoginReq struct {
	Email               string `json:"email"`
	Password            string `json:"password"`
	DevicePlatformId    int    `json:"device_platform_id"`
	DeviceId            string `json:"device_id"`
	DeviceModel         string `json:"device_model"`
	DeviceManufacturer  string `json:"device_manufacturer"`
	NotificationChannel int    `json:"notification_channel"`
	NotificationToken   string `json:"notification_token"`
	AuthProviderId      int    `json:"auth_provider_id"`
	ThirdPartyToken     string `json:"third_party_token"`
}

type UserRefreshSession struct {
	SessionId  string       `json:"session_id"`
	UserId     string       `json:"user_id"`
	DeviceInfo UserLoginReq `json:"device_info"`
}

type JWTOptReq struct {
	Subject   string
	SessionId string
	Lifetime  int
	Purpose   int
	Extras    map[string]string
}

type SignatureReq struct {
	Format string
	Args   []interface{}
}

type ChangePasswordReq struct {
	UserId           string `json:"-"`
	ExistingPassword string `json:"existing_password"`
	NewPassword      string `json:"new_password"`
	Reset            bool   `json:"-"`
}

type ResetPasswordSession struct {
	RequestId     string
	Email         string
	UserSignature string
}

type VerifyEmailSession struct {
	RequestId     string
	Email         string
	UserSignature string
}

type UserSubscriptionReq struct {
	ProviderId             int8                          `json:"provider_id"`
	SubscriptionPlanTypeId int8                          `json:"subscription_plan_type_id"`
	UserId                 string                        `json:"-"`
	Stripe                 *StripeSubscriptionPaymentReq `json:"stripe"`
	VoucherCode            string                        `json:"voucher_code"`
	VoucherProviderRefId   string                        `json:"voucher_provider_ref_id"`
}

type StripeSubscriptionPaymentReq struct {
	PaymentMethodId string `json:"payment_method_id"`
}

type AdvertiserActivationReq struct {
	UserId string `json:"user_id"`
}
