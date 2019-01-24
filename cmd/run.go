package cmd

import (
	"github.com/awillis/sluus/core"
	"github.com/awillis/sluus/pipeline"
	"github.com/spf13/cobra"
	"strconv"
	"syscall"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the sluus service",
	Long:  "instantiate pipelines and execute them",
	Run: func(cmd *cobra.Command, args []string) {
		core.Logger = core.SetupLogger(core.LogConfig("core", strconv.Itoa(syscall.Getpid())))
		if err := pipeline.Assemble(); err != nil {
			core.Logger.Fatal(err)
		}
	},
}
