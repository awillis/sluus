package cmd

import (
	"github.com/awillis/sluus/pipeline"
	"github.com/awillis/sluus/plugin"
)

func Run() {
	pipelineRegistry := pipeline.NewRegistry()
	_ = pipelineRegistry

	pluginRegistry := plugin.NewRegistry()
	_ = pluginRegistry
}
