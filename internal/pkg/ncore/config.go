package ncore

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

func loadYamlConfig(configPath string, requiredConfig []string) (*viper.Viper, error) {
	// Split file name and remove extension
	splits := strings.Split(filepath.Base(configPath), ".")

	// Read config
	config := viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(splits[0])
	config.AddConfigPath(filepath.Dir(configPath))
	err := config.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("nbs-go/ncore: fails to read yaml configuration in %s", configPath)
	}

	// Assert required config
	var invalidKeys = make([]string, 0)
	for _, v := range requiredConfig {
		// If config key is not configured, push keys
		if !config.IsSet(v) {
			invalidKeys = append(invalidKeys, v)
		}
	}

	// If invalidKeys is exist, return error
	if len(invalidKeys) > 0 {
		errMsg := "nbs-go/ncore: missing configuration:"
		for _, v := range invalidKeys {
			errMsg += "\n  " + v
		}
		return nil, errors.New(errMsg)
	}

	return config, nil
}
