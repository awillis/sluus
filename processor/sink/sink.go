package sink

import (
	"sluus/plugin"
)

type Sink struct {
	plugin.PlugBase
}

func (s *Sink) PluginInit() bool {
	return true
}
