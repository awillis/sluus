package processor

import "sluus/plugin"

type Processor interface {
	plugin.Plugin
	Process()
}
