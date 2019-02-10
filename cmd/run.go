package cmd

import (
	"github.com/awillis/sluus/core"
	"github.com/awillis/sluus/pipeline"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var (
	stop = make(chan os.Signal)

	runCmd = &cobra.Command{
		Use:   "run",
		Short: "run the sluus service",
		Long:  "instantiate pipelines and execute them",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			core.Logger = core.SetupLogger(core.LogConfig("core", strconv.Itoa(syscall.Getpid())))
			signal.Notify(stop, syscall.SIGTERM)
			signal.Notify(stop, syscall.SIGINT)
			core.Logger.Infof("sluus version %s", core.VERSION)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := pipeline.Assemble(); err != nil {
				core.Logger.Fatal(err)
			}

			pipeline.Registry.Start()
			core.Logger.Info("sluus started")
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			sig := <-stop
			core.Logger.Infof("received %s", sig.String())
			pipeline.Registry.Stop()
			core.Logger.Info("sluus stopped")
		},
	}
)
