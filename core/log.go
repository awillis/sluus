package core

import (
	"github.com/hashicorp/go-hclog"
	"github.com/rs/zerolog"
	"os"
	"time"
)

var Logger zerolog.Logger

func init() {

	logfh, err := os.OpenFile("kapilary.log", os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		panic(err)
	}

	hclog.DefaultOutput = logfh
	hclog.DefaultLevel = hclog.Info

	Logger := hclog.New(&hclog.LoggerOptions{
		Name: "core",
		Level: hclog.Debug,
		TimeFormat: time.RFC3339,
	})

	Logger.Info("initializing logger")
}

func NewLogger(component string) hclog.Logger {

	logger := hclog.New(&hclog.LoggerOptions{
		Name: component,
		TimeFormat: time.RFC3339,
	})
	return logger
}