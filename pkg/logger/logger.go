// Package logger represents a generic logging interface

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log is a package level variable, every program should access logging function through "Log"
var Log Logger
var atomicLevel zap.AtomicLevel

// Logger represent common interface for logging function
type Logger interface {
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Warnf(format string, args ...interface{})
	Warn(args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
}

// SetLogger is the setter for log variable, it should be the only way to assign value to log
func SetLogger(newLogger Logger) {
	Log = newLogger
}

func SetLogLevel(lvl zapcore.Level) {
	atom.SetLevel(lvl)
}
