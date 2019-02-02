package pipeline

import (
	"github.com/awillis/sluus/plugin"
	"github.com/awillis/sluus/processor"
	"github.com/pkg/errors"
)

var (
	ErrNoConfigFound = errors.New("no configuration files found")
)

func Assemble() (err error) {

	confFiles, err := FindConfigurationFiles()

	if err != nil {
		return
	}

	if len(confFiles) == 0 {
		return ErrNoConfigFound
	}

	for _, file := range confFiles {
		config, err := ReadConfigurationFile(file)

		if err != nil {
			return err
		}

		_ = assembleConfig(config)

		if err != nil {
			return err
		}
	}
	return
}

func assembleConfig(config Config) (pipe *Pipe) {

	pipe = New(config.Name)
	pipe.Logger().Info("assembling pipeline")

	attachProcessorToPipe(pipe, config.Source, plugin.SOURCE)

	for _, conf := range config.Conduit {
		attachProcessorToPipe(pipe, conf, plugin.CONDUIT)
	}

	attachProcessorToPipe(pipe, config.RejectSink, plugin.SINK)
	attachProcessorToPipe(pipe, config.AcceptSink, plugin.SINK)
	pipe.Configure()
	return
}

func attachProcessorToPipe(pipe *Pipe, config ProcessorConfig, pluginType plugin.Type) {

	var proc *processor.Processor

	switch pluginType {
	case plugin.SOURCE:
		proc = processor.New(config.Plugin, plugin.SOURCE)
	case plugin.CONDUIT:
		proc = processor.New(config.Plugin, plugin.CONDUIT)
	case plugin.SINK:
		proc = processor.New(config.Plugin, plugin.SINK)
	}

	if e := proc.Load(); e != nil {
		pipe.Logger().Errorw(e.Error(), "name", proc.Name, "id", proc.ID())
		return
	}

	if e := config.Option.SetPluginOptions(proc.Plugin().Options()); e != nil {
		pipe.Logger().Errorw(e.Error(), "name", proc.Name, "id", proc.ID())
		return
	}

	if e := pipe.Add(proc); e != nil {
		pipe.Logger().Errorw(e.Error(), "name", proc.Name, "id", proc.ID())
	}
}
