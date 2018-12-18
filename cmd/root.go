package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/awillis/sluus/core"
)

var rootCmd = &cobra.Command{
	Use:     "sluus",
	Short:   "A data pipeline toolkit.",
	Long:    "A data pipeline toolkit. See http://sluus.io",
	Version: core.VERSION,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Short, "see 'sluus help' for usage")
	},
}

func init() {
	rootCmd.PersistentFlags().
		StringVar(&core.HOMEDIR, "homedir", core.DEFAULT_HOME, "home directory")
	rootCmd.PersistentFlags().
		StringVar(&core.CONFDIR, "confdir", core.DEFAULT_CONF, "config directory")
	rootCmd.PersistentFlags().
		StringVar(&core.DATADIR, "datadir", core.DEFAULT_DATA, "data directory")
	rootCmd.PersistentFlags().
		StringVar(&core.PLUGDIR, "plugdir",
			strings.Join([]string{core.DEFAULT_HOME, "plugin"}, string(os.PathSeparator)), "plugin directory")
	rootCmd.PersistentFlags().
		StringVar(&core.LOGDIR, "logdir",
			strings.Join([]string{core.DEFAULT_HOME, "log"}, string(os.PathSeparator)), "log directory")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
