// +build !windows

package plugin

import (
	"github.com/pkg/errors"
	"os"
	"plugin"

	"github.com/awillis/sluus/core"
)

var (
	ErrFileNotFound = errors.New("file not found")
	ErrFileIsDir    = errors.New("not a plugin, directory found")
	ErrPluginLoad   = errors.New("error loading plugin")
)

/// NewProcessor loads plugins that implement processor types (e.g. source, sink and conduit).
// It takes the name and type of the processor plugin and invokes its factory constructor.
func NewProcessor(name string, pluginType Type) (procInt Processor, err error) {

	factory, err := LoadByName(name)
	if err != nil {
		return
	}

	return factory.(func(Type) (Processor, error))(pluginType)
}

/// NewMessage loads plugins that implement message types.
// It takes the name of the plugin and invokes its factory constructor.
func NewMessage(name string) (plugInt Interface, err error) {

	if factory, err := LoadByName(name); err == nil {
		plugInt, err = factory.(func(Type) (Interface, error))(MESSAGE)
	}
	return
}

/// LoadByName takes a plugin name and returns the plugin.Symbol for its New factory
func LoadByName(name string) (factory plugin.Symbol, err error) {
	plugFile := core.PLUGDIR + string(os.PathSeparator) + name + ".so"
	return LoadByFile(plugFile)
}

func LoadByFile(plugFile string) (factory plugin.Symbol, err error) {

	if info, err := os.Stat(plugFile); err != nil {
		if os.IsNotExist(err) {
			return factory, errors.Wrap(ErrFileNotFound, plugFile)
		}
		if info.IsDir() {
			return factory, errors.Wrap(ErrFileIsDir, plugFile)
		}
	}

	plug, err := plugin.Open(plugFile)

	if err != nil {
		return factory, errors.Wrap(ErrPluginLoad, err.Error())
	}

	return plug.Lookup("New")
}
