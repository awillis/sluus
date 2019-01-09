package main

import (
	"github.com/awillis/sluus/plugin"
	"github.com/awillis/sluus/processor/kafka"
)

func New(ptype plugin.Type) (plugin.Interface, error) {
	// Plugin builds require exporting the constructor in a separate main package
	return kafka.New(ptype)
}
