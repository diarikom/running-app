package main

import (
	"fmt"
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/pkg/ncore"
	"github.com/diarikom/running-app/running-app-api/internal/pkg/nfacebook"
	"github.com/diarikom/running-app/running-app-api/internal/pkg/njwt"
	"github.com/diarikom/running-app/running-app-api/pkg/nmailgun"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"regexp"
	"time"
)

func initComponents(opt ncore.BootOpt, config *viper.Viper) api.Components {
	// Determine debug mode
	var errDebug bool
	if opt.Environment != ncore.ProductionEnvironment {
		errDebug = true
	}

	// Init error util
	errUtil := api.NewErrorUtil(opt.ErrorCodesPath, errDebug)
	log.Debug("Components.Errors initiated")

	// Init snowflake
	idGen := api.NewSnowflakeGenerator(opt.NodeNo)
	log.Debug("Components.Id initiated")

	// Init jwt issuer
	jwtIssuer := initJwtIssuer(config)
	log.Debug("Components.JWTIssuer initiated")

	// Init mailer
	mailer := initMailer(config)
	log.Debug("Components.Mailer initiated")

	// Init facebook
	fb := initFacebook(config)
	log.Debug("Components.Facebook initiated")

	// Create app components
	return api.Components{
		Errors:    errUtil,
		Id:        idGen,
		JWTIssuer: jwtIssuer,
		Mailer:    mailer,
		Facebook:  fb,
	}
}

func initJwtIssuer(config *viper.Viper) *njwt.Issuer {
	// Retrieve config
	authKey := config.GetString(api.ConfJWTAuthKey)
	defaultLifetime := config.GetInt(api.ConfJWTDefaultLifetime)
	issuer := config.GetString(api.ConfJWTIssuer)

	// Init issuer
	return &njwt.Issuer{
		Key:             []byte(authKey),
		DefaultLifetime: time.Duration(defaultLifetime) & time.Minute,
		Issuer:          issuer,
	}
}

func initMailer(config *viper.Viper) *nmailgun.Mailer {
	// Load mailer config
	var mailerConf nmailgun.Config
	raw := config.Get("components.nmailgun")
	err := mapstructure.Decode(raw, &mailerConf)
	if err != nil {
		panic(fmt.Errorf("running-app-api: unable to retrieve config for Components.Mailer (%s)", err))
	}

	// Validate template_path
	match, err := regexp.MatchString("/$", mailerConf.TemplatePath)
	if err != nil {
		panic(fmt.Errorf("running-app-api: unable to parse template_path in Components.Mailer (%s)", err))
	}

	if !match {
		panic(fmt.Errorf("running-app-api: invalid template_path (%s) in Components.Mailer. template_path must ended with /", mailerConf.TemplatePath))
	}

	// Init mailer
	return nmailgun.New(mailerConf)
}

func initFacebook(config *viper.Viper) *nfacebook.Provider {
	// Retrieve config
	appId := config.GetString(api.ConfFacebookAppId)
	appSecret := config.GetString(api.ConfFacebookAppSecret)

	// Init facebook
	provider, err := nfacebook.NewProvider(nfacebook.ProviderOpt{
		AppId:     appId,
		AppSecret: appSecret,
	})
	if err != nil {
		panic(fmt.Errorf("running-app-api: unable to init Components.Facebook (%s)", err))
	}

	return provider
}
