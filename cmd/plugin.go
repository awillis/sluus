package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"plugin"
	"reflect"

	"os"
	"path/filepath"

	"github.com/awillis/sluus/core"
	splug "github.com/awillis/sluus/plugin"
)

func init() {
	rootCmd.AddCommand(plugCmd)
}

var plugCmd = &cobra.Command{
	Use:   "plugin",
	Short: "list available plugins",
	Long:  "display information about available plugins",
	Run: func(cmd *cobra.Command, args []string) {

		if err := filepath.Walk(core.PLUGDIR, func(path string, info os.FileInfo, err error) (rerr error) {

			if info.IsDir() {
				return
			}

			symbol, err := splug.LoadByFile(path)

			if err != nil {
				return err
			}

			fmt.Println(path)
			proc, err := pluginLookupBySymbol(symbol)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("%s %s\n", proc.Name(), proc.Version())
			}

			if err != nil {
				rerr = err
			}
			return
		}); err != nil {
			fmt.Println(err)
		}
	},
}

func pluginLookupBySymbol(symbol plugin.Symbol) (plugInt splug.Interface, err error) {

	var callInterfaceNew func(splug.Type) (splug.Interface, error)
	var callProcessorNew func(splug.Type) (splug.Processor, error)

	symType := reflect.TypeOf(symbol)
	callIType := reflect.TypeOf(callInterfaceNew)
	callPType := reflect.TypeOf(callProcessorNew)

	for i := 0; i < 4; i++ {
		typ := splug.Type(i)
		if symType.String() == callIType.String() {
			plugInt, err = symbol.(func(splug.Type) (splug.Interface, error))(typ)
		}

		if symType.String() == callPType.String() {
			plugInt, err = symbol.(func(splug.Type) (splug.Processor, error))(typ)
		}
	}
	return plugInt, err
}
