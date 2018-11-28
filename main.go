package main

import (
	"github.com/mitchellh/cli"
	"os"
	"zystus/core"
)

var version = "0.0.1"

func main() {
	cmd := cli.NewCLI("zystus", version)
	cmd.Args = os.Args[1:]
	cmd.Commands = map[string]cli.CommandFactory{}

	status, err := cmd.Run()
	if err != nil {
		core.Logger.Error(err)
	}

	os.Exit(status)
}
