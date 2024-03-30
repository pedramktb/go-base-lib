package logging

import (
	"github.com/ez-as/go-base-lib/pkg/env"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	if env.IsDebug() {
		logger = zap.Must(zap.NewDevelopment())
	} else {
		logger = zap.Must(zap.NewProduction())
	}
}

func Logger() *zap.Logger {
	return logger
}
