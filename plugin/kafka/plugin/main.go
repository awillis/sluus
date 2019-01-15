// +build !windows

package main

import (
	"github.com/awillis/sluus/plugin"
	"github.com/awillis/sluus/plugin/kafka"
)

func New(pluginType plugin.Type) (plugin.Processor, error) {
	// Plugin builds require exporting the constructor in a separate main package
	return kafka.New(pluginType)
}
