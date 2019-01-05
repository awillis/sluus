package cmd

import (
	"fmt"

	"github.com/awillis/sluus/core"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(envCmd)
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "show environment",
	Long:  "show global environment variables",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("HOMEDIR: %s\n", core.HOMEDIR)
		fmt.Printf("CONFDIR: %s\n", core.CONFDIR)
		fmt.Printf("DATADIR: %s\n", core.DATADIR)
		fmt.Printf("PLUGDIR: %s\n", core.PLUGDIR)
		fmt.Printf("LOGDIR: %s\n", core.LOGDIR)
	},
}
