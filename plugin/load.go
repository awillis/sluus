// +build !windows

package plugin

import (
	"fmt"
	"os"
	"plugin"

	"github.com/awillis/sluus/core"
)

/// NewProcessor loads plugins that implement processor types (e.g. source, sink and conduit).
// It takes the name and type of the processor plugin and invokes its constructor.
func NewProcessor(name string, pluginType Type) (procInt Processor, err error) {

	if constructor, err := LoadByName(name); err != nil {
		procInt, err = constructor.(func(Type) (Processor, error))(pluginType)
	}
	return
}

/// NewMessage loads plugins that implement message types.
// It takes the name of the plugin and invokes its constructor.
func NewMessage(name string) (plugInt Interface, err error) {

	if constructor, err := LoadByName(name); err == nil {
		plugInt, err = constructor.(func(Type) (Interface, error))(MESSAGE)
	}
	return
}

/// LoadByName takes a plugin name and returns the plugin.Symbol for its New constructor
func LoadByName(name string) (constructor plugin.Symbol, err error) {
	plugFile := core.PLUGDIR + string(os.PathSeparator) + name + ".so"
	return LoadByFile(plugFile)
}

func LoadByFile(plugFile string) (constructor plugin.Symbol, err error) {
	plug, err := plugin.Open(plugFile)

	if err != nil {
		return constructor, fmt.Errorf("error loading plugin: %s", err)
	}

	return plug.Lookup("New")
}
