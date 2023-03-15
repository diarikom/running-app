package main

import (
	"flag"
	"github.com/diarikom/running-app/running-app-api/internal/pkg/ncore"
)

// AppFlags represents list of available commands and options
type AppFlags struct {
	CmdShowHelp    *bool
	CmdShowVersion *bool
	OptEnvironment *int
	OptDir         *string
	OptConfig      *string
	OptErrorCodes  *string
	OptNodeNo      *int64
}

// initAppFlags loads commands name, default value and description
func initAppFlags() AppFlags {
	return AppFlags{
		CmdShowHelp:    flag.Bool("help", false, "Command: Show available commands and options"),
		CmdShowVersion: flag.Bool("version", false, "Command: Show version"),
		OptEnvironment: flag.Int("env", ncore.DevelopmentEnvironment, "Option: Set app environment"),
		OptConfig:      flag.String("config", "", "Option: Set config file"),
		OptErrorCodes:  flag.String("error_codes", "", "Option: Set error codes file"),
		OptDir:         flag.String("dir", ".", "Option: Set working directory"),
		OptNodeNo:      flag.Int64("node_no", 1, "Option: App instance number"),
	}
}
