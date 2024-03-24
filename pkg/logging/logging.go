package logging

import (
	"github.com/ez-as/ironlink-base-lib/env"
	"go.uber.org/zap"
)

func init() {
	var logger *zap.Logger
	if env.IsProd() {
		logger = zap.Must(zap.NewProduction())
	} else {
		logger = zap.Must(zap.NewDevelopment())
	}
	zap.ReplaceGlobals(logger)
}

func Logger() *zap.Logger {
	return zap.L()
}
