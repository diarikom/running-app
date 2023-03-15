package nlogrus

import (
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

// constants
const DefaultHostname = "unknown"

type LoggerOpt struct {
	Hostname string
	LogLevel int
}

// GetLoggerOptEnv retrieve from env
func GetLoggerOptEnv() LoggerOpt {
	// Retrieve hostname option
	hostname, err := os.Hostname()
	if err != nil {
		hostname = DefaultHostname
	}

	// Retrieve logger option
	var level int
	levelStr, ok := os.LookupEnv(nlog.EnvLogLevel)
	// If env not set, set to default level
	if !ok {
		level = nlog.DefaultLevel
	} else {
		// Parse to int
		tmp, err := strconv.Atoi(levelStr)
		// If level is unable to parse, set to default level
		if err != nil {
			level = nlog.DefaultLevel
		} else {
			// Else, set level
			level = tmp
		}
	}

	// Return option
	return LoggerOpt{
		Hostname: hostname,
		LogLevel: level,
	}
}

func convertLevel(level int) logrus.Level {
	switch level {
	case nlog.LevelPanic:
		return logrus.PanicLevel
	case nlog.LevelFatal:
		return logrus.FatalLevel
	case nlog.LevelError:
		return logrus.ErrorLevel
	case nlog.LevelWarn:
		return logrus.WarnLevel
	case nlog.LevelInfo:
		return logrus.InfoLevel
	case nlog.LevelDebug:
		return logrus.DebugLevel
	default:
		panic("nlogrus: unsupported logger level.")
	}
}
