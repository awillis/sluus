// +build !windows

package main

import (
	"github.com/awillis/sluus/plugin"
	"github.com/awillis/sluus/plugin/noop"
)

func New(pluginType plugin.Type) (plugin.Interface, error) {
	// Plugin builds require exporting the constructor in a separate main package
	return noop.New(pluginType)
}
