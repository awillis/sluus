// +build !windows

package main

import (
	"github.com/awillis/sluus/plugin"
	"github.com/awillis/sluus/processor/noop"
)

func New(pluginType plugin.Type) (plugin.Processor, error) {
	// Plugin builds require exporting the constructor in a separate main package
	return noop.New(pluginType)
}
