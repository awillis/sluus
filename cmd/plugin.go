package cmd

import (
	"fmt"
	"github.com/awillis/sluus/plugin"
	"reflect"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(plugCmd)
}

var plugCmd = &cobra.Command{
	Use:   "plugins",
	Short: "list available plugins",
	Long:  "display information about available plugins",
	Run: func(cmd *cobra.Command, args []string) {
		proc, err := plugin.NewProcessor("kafka", plugin.SINK)
		if err != nil {
			panic(err)
		}

		typ := reflect.TypeOf(proc)
		fmt.Printf("name: %s, version: %s, ptype: %+v, type: %s\n", proc.Name(), proc.Version(), proc.Type(), typ.String())

		proc, err = plugin.NewProcessor("kafka", plugin.SOURCE)
		if err != nil {
			panic(err)
		}

		typ = reflect.TypeOf(proc)
		fmt.Printf("name: %s, version: %s, ptype: %+v, type: %s\n", proc.Name(), proc.Version(), proc.Type(), typ.String())
	},
}
