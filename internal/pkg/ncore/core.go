package ncore

import (
	"encoding/json"
	"github.com/spf13/viper"
	"time"
)

// App Environments
const (
	DevelopmentEnvironment = iota
	ProductionEnvironment
	TestingEnvironment
)

type Core struct {
	Environment int
	Dir         string
	Config      *viper.Viper
	// Metadata
	bootStart time.Time
	bootOpt   BootOpt
}

type BootOpt struct {
	Environment    int
	Dir            string
	ConfigPath     string
	ConfigRequired []string `json:"-"`
	ErrorCodesPath string
	NodeNo         int64
}

func Boot(bootOpt BootOpt) *Core {
	// Init Core
	core := Core{
		Environment: bootOpt.Environment,
		Dir:         bootOpt.Dir,
		bootStart:   time.Now(),
		bootOpt:     bootOpt,
	}

	// Load config
	config, err := loadYamlConfig(bootOpt.ConfigPath, bootOpt.ConfigRequired)
	if err != nil {
		panic(err)
	}
	core.Config = config

	return &core
}

func (c *Core) Uptime() time.Duration {
	return time.Since(c.bootStart)
}

func (c *Core) GetBootOpt() string {
	bootOpt, _ := json.MarshalIndent(c.bootOpt, "", "  ")
	return string(bootOpt)
}
