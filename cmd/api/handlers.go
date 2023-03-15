package main

import (
	"fmt"
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/service"
	"github.com/diarikom/running-app/running-app-api/pkg/nhttp"
	"net/http"
)

type Handlers struct {
	ApiStatus        nhttp.HandlerFunc
	Asset            *service.AssetHandler
	User             *service.UserHandler
	Run              *service.RunHandler
	DiscoverContent  *service.DiscoverContentHandler
	AdTag            *service.AdTagHandler
	MilestoneHandler *service.MilestoneHandler
	Initiative       *service.InitiativeHandler
	SubscriptionPlan *service.SubscriptionPlanHandler
	SiteSetting      *service.SiteSettingHandler
}

func initHandlers(app *api.Api) Handlers {
	user := service.NewUserHandler(app)
	asset := service.NewAssetHandler(app)
	run := service.NewRunHandler(app)
	discoverContent := service.NewDiscoverContentHandler(app)
	tag := service.NewTagHandler(app)
	milestoneHandler := service.NewMilestoneHandler(app)
	initiative := service.NewInitiativeHandler(app)
	subscriptionPlan := service.NewSubscriptionPlanHandler(app)
	siteSetting := service.NewSiteSettingHandler(app)

	return Handlers{
		ApiStatus:        newApiStatusHandler(app),
		Asset:            &asset,
		User:             &user,
		Run:              &run,
		DiscoverContent:  &discoverContent,
		AdTag:            &tag,
		MilestoneHandler: &milestoneHandler,
		Initiative:       &initiative,
		SubscriptionPlan: &subscriptionPlan,
		SiteSetting:      &siteSetting,
	}
}

func newApiStatusHandler(app *api.Api) nhttp.HandlerFunc {
	return func(_ *http.Request) (*nhttp.Success, error) {
		res := &nhttp.Success{
			Result: map[string]string{
				"build_version": AppVersion,
				"uptime":        fmt.Sprintf("%s", app.Uptime()),
			},
		}
		return res, nil
	}
}
