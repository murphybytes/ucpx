package common

import (
	"errors"
	"log"
	"os"
)

// Logger wraps Go logger but allows variable verbosity
type Logger interface {
	LogInfo(v ...interface{})
	LogWarn(v ...interface{})
	LogError(v ...interface{})
	LogFatal(v ...interface{})
}

type logger struct {
	verbosity int
	logFunc   func(v ...interface{})
}

// NewLogger creates a logger
func NewLogger(f *Flags) (l Logger, e error) {
	return newLogger(f, log.Print)
}

// NewLogger creates a logger
func newLogger(f *Flags, logFunc func(v ...interface{})) (l Logger, e error) {
	levels := map[string]int{
		logInfo:  3,
		logWarn:  2,
		logError: 1,
	}

	if level, ok := levels[f.LogLevel]; ok {
		l = &logger{
			verbosity: level,
			logFunc:   logFunc,
		}
	} else {
		e = errors.New("Invalid log level")
	}

	return
}

func (l *logger) LogInfo(v ...interface{}) {
	if l.verbosity >= 3 {
		l.logFunc(v...)
	}
}

func (l *logger) LogWarn(v ...interface{}) {
	if l.verbosity >= 2 {
		l.logFunc(v...)
	}

}

func (l *logger) LogError(v ...interface{}) {
	if l.verbosity >= 1 {
		l.logFunc(v...)
	}
}

func (l *logger) LogFatal(v ...interface{}) {
	l.logFunc(v...)
	os.Exit(1)

}
