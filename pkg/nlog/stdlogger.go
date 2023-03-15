package nlog

import (
	"fmt"
	"io"
	stdLog "log"
	"os"
)

type StdLogger struct {
	level     int
	levelStr  map[int]string
	skipTrace int
	writer    *stdLog.Logger
}

func (l StdLogger) Fatal(msg string, err error) {
	l.printErr(LevelFatal, msg, err)
}

func (l StdLogger) Fatalf(format string, args ...interface{}) {
	l.printf(LevelFatal, format, args...)
}

func (l StdLogger) Error(msg string, err error) {
	l.printErr(LevelError, msg, err)
}

func (l StdLogger) Errorf(format string, args ...interface{}) {
	l.printf(LevelError, format, args...)
}

func (l StdLogger) Warn(msg string) {
	l.print(LevelWarn, msg)
}

func (l StdLogger) Warnf(format string, args ...interface{}) {
	l.printf(LevelWarn, format, args...)
}

func (l StdLogger) Info(msg string) {
	l.print(LevelInfo, msg)
}

func (l StdLogger) Infof(format string, args ...interface{}) {
	l.printf(LevelInfo, format, args...)
}

func (l *StdLogger) Debug(msg string) {
	l.print(LevelDebug, msg)
}

func (l *StdLogger) Debugf(format string, args ...interface{}) {
	l.printf(LevelDebug, format, args...)
}

func NewStdLogger(level int, w io.Writer, prefix string, flags int) Logger {
	// If writer is nil, set default writer to Stdout
	if w == nil {
		w = os.Stdout
	}

	// Init standard logger instance
	l := StdLogger{
		level: level,
		levelStr: map[int]string{
			LevelPanic: "PANIC",
			LevelFatal: "FATAL",
			LevelError: "ERROR",
			LevelWarn:  "WARN",
			LevelInfo:  "INFO",
			LevelDebug: "DEBUG",
		},
		skipTrace: 2,
		writer:    stdLog.New(w, prefix, flags),
	}
	return &l
}

func (l *StdLogger) print(outLevel int, msg string) {
	// if output level is greater than log level, don't print
	if outLevel > l.level {
		return
	}

	// print log
	l.writer.Printf("[%s] %s", l.levelStr[outLevel], msg)
}

func (l *StdLogger) printf(outLevel int, pattern string, args ...interface{}) {
	// if output level is greater than log level, don't print
	if outLevel > l.level {
		return
	}

	// Generate level prefix
	levelStr := fmt.Sprintf("[%s] ", l.levelStr[outLevel])

	// print log pattern
	l.writer.Printf(levelStr+pattern, args...)
}

// printErr trace error and print
func (l *StdLogger) printErr(outLevel int, msg string, err error) {
	// if output level is greater than log level, don't print
	if outLevel > l.level {
		return
	}

	// Trace error
	filePath, line := Trace(l.skipTrace)

	// Get level string
	levelStr := l.levelStr[outLevel]

	l.writer.Printf("[%s] Error: %s", levelStr, err)
	l.writer.Printf("[%s] Description: %s", levelStr, msg)
	l.writer.Printf("[%s] Trace: %s:%d", levelStr, filePath, line)
}
