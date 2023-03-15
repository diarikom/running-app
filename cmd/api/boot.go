package main

import (
	"github.com/diarikom/running-app/running-app-api/internal/api/service"
	_ "github.com/lib/pq"

	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/pkg/ncore"
)

func boot(flags AppFlags) api.Api {
	// Load boot option
	bootOpt := loadBootOpt(flags)

	// Boot core
	core := ncore.Boot(bootOpt)

	// Boot datasources
	datasources := initDatasources(core.Config)

	// Boot Components
	components := initComponents(bootOpt, core.Config)

	// Create base url
	baseUrl, port := getBaseUrl(core.Config)

	// Init app
	app := api.Api{
		BaseUrl:     baseUrl,
		Port:        port,
		Core:        core,
		Datasources: &datasources,
		Components:  &components,
		Logger:      log,
		Services: &api.Services{
			Asset:            new(service.Asset),
			Auth:             new(service.Authenticator),
			User:             new(service.User),
			Run:              new(service.Run),
			DiscoverContent:  new(service.DiscoverContent),
			Tag:              new(service.AdTagService),
			MilestoneService: new(service.MilestoneService),
			Credit:           new(service.CreditService),
			Initiative:       new(service.Initiative),
			SubscriptionPlan: new(service.SubscriptionService),
			SiteSetting:      new(service.SiteSettingService),
		},
	}

	// Boot services
	initServices(&app)

	return app
}

// loadBootOpt parse boot options from cli arguments
func loadBootOpt(flags AppFlags) ncore.BootOpt {
	// Init opt
	opt := ncore.BootOpt{
		Environment:    *flags.OptEnvironment,
		Dir:            *flags.OptDir,
		NodeNo:         *flags.OptNodeNo,
		ConfigRequired: api.RequiredConfig,
	}

	// Get Config Path
	if *flags.OptConfig == "" {
		opt.ConfigPath = opt.Dir + "/config.yml"
	} else {
		opt.ConfigPath = *flags.OptConfig
	}

	// Get Error Codes Path
	if *flags.OptErrorCodes == "" {
		opt.ErrorCodesPath = opt.Dir + "/error-codes.yml"
	} else {
		opt.ErrorCodesPath = *flags.OptErrorCodes
	}

	return opt
}
