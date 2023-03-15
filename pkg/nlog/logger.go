package nlog

import (
	stdLog "log"
	"os"
	"runtime"
	"sync"
)

// Log Levels as defined in RFC5424.
const (
	LevelPanic = iota
	LevelFatal
	LevelCritical
	LevelError
	LevelWarn
	LevelNotice
	LevelInfo
	LevelDebug
	LevelTrace
)

// Configuration constants.
const (
	// Env Config Keys
	EnvLogLevel = "NLOG_LEVEL"
	// Default Values
	DefaultLevel = LevelError
)

/// Logger contract defines methods that must be available for a Logger.
///
/// Fatal must write an error, message that explaining the error and where its occurred in FATAL level.
/// To trace message, use Trace function and skip 1.
///
/// Fatalf must write a formatted message and where its occurred in FATAL level.
/// To trace message, use Trace function and skip 1.
///
/// Error must write an error, message that explaining the error and where its occurred in ERROR level.
/// To trace message, use Trace function and skip 1.
///
/// Errorf must write a formatted message and where its occurred in ERROR level.
/// To trace message, use Trace function and skip 1.
///
/// Warn must write a message in WARN level.
///
/// Warnf must write a formatted message in WARN level.
///
/// Info must write a message in INFO level.
///
/// Infof must write a formatted message in INFO level.
///
/// Debug must write a message in DEBUG level.
///
/// Debugf must write a formatted message in DEBUG level.
type Logger interface {
	Fatal(msg string, err error)
	Fatalf(format string, args ...interface{})
	Error(msg string, err error)
	Errorf(format string, args ...interface{})
	Warn(msg string)
	Warnf(format string, args ...interface{})
	Info(msg string)
	Infof(format string, args ...interface{})
	Debug(msg string)
	Debugf(format string, args ...interface{})
}

/// log is a singleton logger instance
var log Logger
var logMutex sync.RWMutex

/// Get retrieve singleton logger instance
func Get() Logger {
	// If log is nil, initiate standard logger
	if log == nil {
		// Init standard logger
		l := NewStdLogger(DefaultLevel, os.Stderr, "", stdLog.LstdFlags)

		// Register logger
		Register(l)
		log.Debug("No logger found. StdLogger initiated")
	}
	return log
}

/// Register logger instance
func Register(l Logger) {
	// If logger is nil, return error
	if l == nil {
		panic("nbs-go/nlog: logger to be registered is nil")
	}

	// Set logger
	logMutex.Lock()
	defer logMutex.Unlock()
	log = l
}

/// Trace retrieve where the code is being called and returns full path of file where the error occurred
func Trace(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		file = "<???>"
		line = 1
	}
	return file, line
}
