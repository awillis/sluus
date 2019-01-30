// +build !windows

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/spf13/cobra"

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

			return displayPlugin(path)
		}); err != nil {
			fmt.Println(err)
		}
	},
}

func displayPlugin(path string) (err error) {

	var callInterfaceNew func(splug.Type) (splug.Interface, error)
	var callProcessorNew func(splug.Type) (splug.Interface, error)

	symbol, err := splug.LoadByFile(path)

	if err != nil {
		return
	}

	symType := reflect.TypeOf(symbol)
	callIType := reflect.TypeOf(callInterfaceNew)
	callPType := reflect.TypeOf(callProcessorNew)

	for i := 0; i < 4; i++ {
		typ := splug.Type(i)

		if symType.String() == callIType.String() {
			if plugInt, perror := symbol.(func(splug.Type) (splug.Interface, error))(typ); perror == nil {
				fmt.Printf("name: %s, version: %s, type: %d\n", plugInt.Name(), plugInt.Version(), plugInt.Type())
			} else if perror.Error() == "unimplemented plugin" {
				continue
			} else {
				err = perror
			}
		}

		if symType.String() == callPType.String() {
			if plugInt, perror := symbol.(func(splug.Type) (splug.Interface, error))(typ); perror == nil {
				fmt.Printf("name: %s, version: %s, type: %d\n", plugInt.Name(), plugInt.Version(), plugInt.Type())
			} else if perror.Error() == "unimplemented plugin" {
				continue
			} else {
				err = perror
			}
		}
	}
	return err
}
