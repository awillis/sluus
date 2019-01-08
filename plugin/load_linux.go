package plugin

import (
	"fmt"
	"os"
	"plugin"
	"strings"

	"github.com/awillis/sluus/core"
)

func Load(name string, ptype core.PluginType) (Processor, error) {

	var err error

	plugfile := strings.Join([]string{core.PLUGDIR, name + ".so"}, string(os.PathSeparator))
	plug, err := plugin.Open(plugfile)

	if err != nil {
		return nil, fmt.Errorf("error loading plugin: %s", err)
	}

	factory, err := plug.Lookup("New")

	if err != nil {
		return nil, err
	}

	proc := factory.(func(core.PluginType) Processor)(ptype)
	return proc, err
}
