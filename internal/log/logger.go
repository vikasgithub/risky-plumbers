// Package log provides context-aware and structured logging capabilities.
package log

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

type contextKey int

const (
	requestIDKey contextKey = iota
	correlationIDKey
)

// New creates a new logger using the default configuration.
func New() Logger {
	//TODO add logging based on Production using an environment variable
	l, _ := zap.NewDevelopment()
	return Logger{l.Sugar()}
}
