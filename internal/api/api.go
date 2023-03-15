package api

import (
	"github.com/diarikom/running-app/running-app-api/internal/pkg/ncore"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
	"net/url"
)

type Api struct {
	BaseUrl     *url.URL
	Port        string
	Datasources *Datasources
	Components  *Components
	Services    *Services
	Logger      nlog.Logger

	*ncore.Core
}

type Datasources struct {
	Db    *nsql.SqlDatabase
	Asset S3Provider
}

type Components struct {
	Errors    *Errors
	Id        *SnowflakeGen
	JWTIssuer JWTIssuerComponent
	Mailer    MailerComponent
	Facebook  FacebookProviderComponent
}

type Services struct {
	Asset            AssetService
	Auth             AuthenticatorService
	User             UserService
	Run              RunService
	DiscoverContent  DiscoverContentService
	Tag              AdTagService
	MilestoneService MilestoneService
	Credit           CreditService
	Initiative       InitiativeService
	SubscriptionPlan SubscriptionPlanService
	SiteSetting      SiteSettingService
}
