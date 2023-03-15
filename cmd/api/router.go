package main

import (
	"fmt"
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/pkg/nhttp"
	"github.com/spf13/viper"
	"net/url"
)

func getBaseUrl(config *viper.Viper) (*url.URL, string) {
	// Retrieve port
	port := config.GetString(api.ConfServerPort)
	basePath := config.GetString(api.ConfServerBasePath)

	// If host is set, parse host
	if config.IsSet(api.ConfServerHost) {
		host := config.GetString(api.ConfServerHost)
		rawUrl := host + basePath

		u, err := url.Parse(rawUrl)
		if err != nil {
			panic(fmt.Errorf("running-app-api: unable to parse server.host value (%s)", err))
		}

		return u, port
	}

	// Else, build url
	urlConf := nhttp.UrlConfig{
		Scheme:   config.GetString(api.ConfServerScheme),
		Port:     config.GetString(api.ConfServerPort),
		BasePath: config.GetString(api.ConfServerBasePath),
	}

	return nhttp.BuildUrl(urlConf), port
}
