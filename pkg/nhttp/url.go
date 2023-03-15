package nhttp

import "net/url"

type UrlConfig struct {
	Scheme   string `yaml:"scheme"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	BasePath string `yaml:"api_path"`
}

func BuildUrl(config UrlConfig) *url.URL {
	// Determine base path
	if config.BasePath == "" {
		config.BasePath = "/"
	}

	// Determine scheme
	switch config.Scheme {
	case "https":
		config.Port = ""
	case "":
		config.Scheme = "http"
	}

	// Determine host
	if config.Host == "" {
		config.Host = "localhost"
	}

	// Determine port and host
	if config.Port != "" {
		config.Host += ":" + config.Port
	}

	// Build url
	baseUrl := url.URL{
		Scheme: config.Scheme,
		Host:   config.Host,
		Path:   config.BasePath,
	}

	return &baseUrl
}
