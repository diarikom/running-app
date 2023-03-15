package nlogrus

import (
	"fmt"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/sirupsen/logrus"
)

func NewConsoleLogger(opt LoggerOpt) nlog.Logger {
	// Convert level to logrus level
	lv := convertLevel(opt.LogLevel)

	// Initiate writer
	logrus.SetLevel(lv)

	// Initiate logger instance
	logger := ConsoleLogger{
		hostname:  opt.Hostname,
		level:     opt.LogLevel,
		writer:    logrus.StandardLogger(),
		skipTrace: 1,
	}

	return &logger
}

// ConsoleLogger is an implementation of Logger that prints output to console
type ConsoleLogger struct {
	hostname  string
	level     int
	writer    *logrus.Logger
	skipTrace int
}

func (l *ConsoleLogger) addFields() *logrus.Entry {
	return l.writer.
		WithField("hostname", l.hostname)
}

func (l *ConsoleLogger) Debug(msg string) {
	l.addFields().Debug(msg)
}

func (l *ConsoleLogger) Debugf(format string, args ...interface{}) {
	l.addFields().Debugf(format, args...)
}

func (l *ConsoleLogger) Info(msg string) {
	l.addFields().Info(msg)
}

func (l *ConsoleLogger) Infof(format string, args ...interface{}) {
	l.addFields().Infof(format, args...)
}

func (l *ConsoleLogger) Warn(msg string) {
	l.addFields().Warn(msg)
}

func (l *ConsoleLogger) Warnf(format string, args ...interface{}) {
	l.addFields().Warnf(format, args...)
}

func (l *ConsoleLogger) Error(msg string, err error) {
	// Trace error
	file, line := nlog.Trace(l.skipTrace)
	l.addFields().
		WithField("trace", fmt.Sprintf("%s:%d", file, line)).
		WithField("error_msg", err).
		Error(msg)
}

func (l *ConsoleLogger) Errorf(format string, args ...interface{}) {
	// Trace error
	file, line := nlog.Trace(l.skipTrace)
	l.addFields().
		WithField("trace", fmt.Sprintf("%s:%d", file, line)).
		Errorf(format, args...)
}

func (l *ConsoleLogger) Fatal(msg string, err error) {
	// Trace error
	file, line := nlog.Trace(l.skipTrace)
	l.addFields().
		WithField("trace", fmt.Sprintf("%s:%d", file, line)).
		WithField("error_msg", err).
		Fatal(msg)
}

func (l *ConsoleLogger) Fatalf(format string, args ...interface{}) {
	// Trace error
	file, line := nlog.Trace(l.skipTrace)
	l.addFields().
		WithField("trace", fmt.Sprintf("%s:%d", file, line)).
		Fatalf(format, args...)
}
