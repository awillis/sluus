package cmd

import (
	"fmt"
	"github.com/awillis/sluus/core"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

func init() {
	signal.Notify(terminate)
}

var (
	rootCmd = &cobra.Command{
		Use:     "sluus",
		Short:   "A data pipeline toolkit.",
		Long:    "A data pipeline toolkit. See http://sluus.io",
		Version: core.VERSION,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Short, "see 'sluus help' for usage")
		},
	}
	terminate = make(chan os.Signal, 1)
)

func Execute() {
	rootCmd.PersistentFlags().
		StringVar(&core.HOMEDIR, "homedir", core.HOMEDIR, "home directory")
	rootCmd.PersistentFlags().
		StringVar(&core.CONFDIR, "confdir", core.CONFDIR, "config directory")
	rootCmd.PersistentFlags().
		StringVar(&core.DATADIR, "datadir", core.DATADIR, "data directory")
	rootCmd.PersistentFlags().
		StringVar(&core.PLUGDIR, "plugdir", core.PLUGDIR, "plugin directory")
	rootCmd.PersistentFlags().
		StringVar(&core.LOGDIR, "logdir", core.LOGDIR, "log directory")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
