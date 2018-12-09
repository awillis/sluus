package conduit

import (
	"sluus/plugin"
)

type Conduit struct {
	plugin.PlugBase
}

func (c *Conduit) PluginInit() bool {
	return true
}
