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

func New(name string, pluginType Type) (plug Interface, err error) {
	if factory, err := LoadByName(name); err == nil {
		if plug, ok := factory.(func(Type) (Interface, error)); ok {
			return plug(pluginType)
		}
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
