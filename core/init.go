package core

import (
	"github.com/mitchellh/cli"
	"go.uber.org/zap"
)

var HOME string
var CONF string
var DATA string

var Runner *cli.CommandFactory

func init() {

	Logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	Logger.Info("initializing logger")

	Runner = new(cli.CommandFactory)
}
