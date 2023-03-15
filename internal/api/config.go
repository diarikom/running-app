package api

const (
	ConfServerBasePath = "server.base_path"
	ConfServerHost     = "server.host"
	ConfServerPort     = "server.port"
	ConfServerScheme   = "server.scheme"

	ConfAppClientSecret                   = "auth.app_client_secret"
	ConfDashboardClientSecret             = "auth.dashboard_client_secret"
	ConfUserAccessLifetime                = "auth.token_lifetime.user_access"
	ConfResetPasswordTokenLifetime        = "auth.token_lifetime.reset_password"
	ConfVerifyEmailTokenLifetime          = "auth.token_lifetime.verify_email"
	ConfSignatureSaltResetPasswordSubject = "auth.signature_salt.reset_password_subject"
	ConfSignatureSaltEmailVerifySubject   = "auth.signature_salt.verify_email_subject"

	ConfJWTAuthKey         = "components.njwt.auth_key"
	ConfJWTDefaultLifetime = "components.njwt.default_lifetime"
	ConfJWTIssuer          = "components.njwt.issuer"

	ConfFacebookAppId     = "components.facebook.app_id"
	ConfFacebookAppSecret = "components.facebook.app_secret"

	ConfAssetEndpoint        = "datasources.asset.endpoint"
	ConfAssetUseSSL          = "datasources.asset.use_ssl"
	ConfAssetAccessKeyId     = "datasources.asset.access_key_id"
	ConfAssetSecretAccessKey = "datasources.asset.secret_access_key"
	ConfAssetBucketName      = "datasources.asset.bucket_name"
	ConfAssetRegion          = "datasources.asset.region"
	ConfAssetBaseUrl         = "asset.base_url"

	ConfStripeSecretKey = "components.stripe.secret_key"

	ConfDashboardUrl                 = "components.dashboard.url"
	ConfAdvertiserActivationLifetime = "components.dashboard.advertiser_activation_lifetime"
)

var RequiredConfig = []string{
	// Authentication
	ConfAppClientSecret,
	ConfDashboardClientSecret,
	ConfUserAccessLifetime,
	ConfResetPasswordTokenLifetime,
	ConfSignatureSaltResetPasswordSubject,
	ConfSignatureSaltEmailVerifySubject,

	// Datasources.Db
	"datasources.db.driver",
	"datasources.db.host",
	"datasources.db.port",
	"datasources.db.username",
	"datasources.db.password",
	"datasources.db.database",

	// Datasources.Asset
	ConfAssetEndpoint,
	ConfAssetAccessKeyId,
	ConfAssetSecretAccessKey,
	ConfAssetBucketName,
	ConfAssetBaseUrl,

	// Components.JWTIssuer
	ConfJWTAuthKey,
	ConfJWTDefaultLifetime,
	ConfJWTIssuer,

	// Components.Facebook
	ConfFacebookAppId,
	ConfFacebookAppSecret,

	// Components.Mailer
	"components.nmailgun.domain",
	"components.nmailgun.private_api_key",
	"components.nmailgun.template_path",

	// Stripe
	ConfStripeSecretKey,

	// Components.Dashboard
	ConfDashboardUrl,
}
