package cmd

import (
	"github.com/awillis/sluus/pipeline"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the sluus service",
	Long:  "instantiate pipelines and execute them",
	Run: func(cmd *cobra.Command, args []string) {
		pipelineRegistry := pipeline.NewRegistry()
		pipelineRegistry.AddPipeline(pipeline.NewPipeline())
	},
}
