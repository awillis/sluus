package core

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func NewLogger(component string) *zap.SugaredLogger {

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	return logger.Sugar()
}
