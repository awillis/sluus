package cmd

import (
	"github.com/awillis/sluus/pipeline"
)

func Run() {
	pipelineRegistry := pipeline.NewRegistry()
	_ = pipelineRegistry
}
