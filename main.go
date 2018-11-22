package main

import (
	"github.com/mitchellh/cli"
	"kapilary/core"
	"os"
)

var version  = "0.0.1"

func main()  {
	cmd := cli.NewCLI("kapilary", version)
	cmd.Args = os.Args[1:]
	cmd.Commands = map[string]cli.CommandFactory{

	}

	status, err := cmd.Run()
	if err != nil {
		core.Logger.Error().Msg(err.Error())
	}

	core.Logger.Info().Msg("initializing")
	os.Exit(status)
}