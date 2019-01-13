// +build !windows

package main

import (
	"github.com/awillis/sluus/plugin"
	"github.com/awillis/sluus/processor/tcp"
)

func New(ptype plugin.Type) (plugin.Processor, error) {
	// Plugin builds require exporting the constructor in a separate main package
	return tcp.New(ptype)
}
