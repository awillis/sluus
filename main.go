package main

import (
	"github.com/mitchellh/cli"
	"os"
	"sluus/core"
)

var version string

func main() {
	cmd := cli.NewCLI("sluus", version)
	cmd.Args = os.Args[1:]
	cmd.Commands = map[string]cli.CommandFactory{}

	status, err := cmd.Run()
	if err != nil {
		core.Logger.Error(err)
	}

	os.Exit(status)
}
