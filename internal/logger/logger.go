package logger

import (
	"go.uber.org/zap"
)

// global var for logging
var logger *zap.Logger

// logger ctor
func Logger(){
    logger, _ = zap.NewProduction()
}

// return global logger
func L() *zap.Logger {
	return logger
}
