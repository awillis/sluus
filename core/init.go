package core

import (
	"go.uber.org/zap"
)

var HOME string
var CONF string
var DATA string

func init() {

	Logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	Logger.Info("initializing logger")

}
