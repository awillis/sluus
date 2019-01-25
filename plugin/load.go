// +build !windows

package plugin

import (
	"fmt"
	"github.com/awillis/sluus/core"
	"os"
	"plugin"
)

type (
	pConstructor func(Type) (Processor, error)
	iConstructor func(Type) (Interface, error)
)

/// NewProcessor loads plugins that implement processor types (e.g. source, sink and conduit).
// It takes the name and type of the processor plugin and invokes its factory constructor.
func NewProcessor(name string, pluginType Type) (procInt Processor, err error) {

	if factory, err := LoadByName(name); err != nil {
		procInt, err = factory.(pConstructor)(pluginType)
	}
	return
}

/// NewMessage loads plugins that implement message types.
// It takes the name of the plugin and invokes its factory constructor.
func NewMessage(name string) (plugInt Interface, err error) {

	if factory, err := LoadByName(name); err == nil {
		plugInt, err = factory.(iConstructor)(MESSAGE)
	}
	return
}

/// LoadByName takes a plugin name and returns the plugin.Symbol for its New factory
func LoadByName(name string) (factory plugin.Symbol, err error) {
	plugFile := core.PLUGDIR + string(os.PathSeparator) + name + ".so"
	return LoadByFile(plugFile)
}

func LoadByFile(plugFile string) (factory plugin.Symbol, err error) {
	plug, err := plugin.Open(plugFile)

	if err != nil {
		return factory, fmt.Errorf("error loading plugin: %s", err)
	}

	return plug.Lookup("New")
}
