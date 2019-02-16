package cmd

import (
	"context"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/awillis/sluus/core"
	"github.com/awillis/sluus/pipeline"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "run the sluus service",
		Long:  "instantiate pipelines and execute them",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			core.Logger = core.SetupLogger(core.LogConfig("core", strconv.Itoa(syscall.Getpid())))
			core.Logger.Infof("sluus version %s", core.VERSION)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := pipeline.Assemble(); err != nil {
				core.Logger.Fatal(err)
			}

			pipeline.Registry.Start(rootContext)
			core.Logger.Info("sluus started")
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			select {
			case <-complete:
				pipeline.Registry.Stop(rootCancel)
				core.Logger.Info("sluus stopped")
			}
			return core.Logger.Sync()
		},
	}
	rootContext, rootCancel = context.WithCancel(context.Background())
)
