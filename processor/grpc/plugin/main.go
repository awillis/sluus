// +build !windows

package main

import (
	"github.com/awillis/sluus/plugin"
	"github.com/awillis/sluus/processor/grpc"
)

func New(pluginType plugin.Type) (plugin.Processor, error) {
	// Plugin builds require exporting the constructor in a separate main package
	return grpc.New(pluginType)
}
