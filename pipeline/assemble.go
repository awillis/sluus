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

		if err = assembleConfig(config); err != nil {
			return err
		}
	}
	return
}

func assembleConfig(config Config) (err error) {

	pipe := New(config.Name)

	source := processor.New(config.Source.Plugin, plugin.SOURCE)
	if err = source.Load(config.Source.Options); err != nil {
		return
	}

	if err = pipe.SetSource(source); err != nil {
		return
	}

	accept := processor.New(config.AcceptSink.Plugin, plugin.SINK)
	if err = accept.Load(config.AcceptSink.Options); err != nil {
		return
	}

	reject := processor.New(config.RejectSink.Plugin, plugin.SINK)
	if err = reject.Load(config.RejectSink.Options); err != nil {
		return
	}

	if err = pipe.SetSinks(accept, reject); err != nil {
		return
	}

	for _, conf := range config.Conduit {
		conduit := processor.New(conf.Plugin, plugin.CONDUIT)
		if err = conduit.Load(conf.Options); err != nil {
			return
		}

		if err = pipe.AddConduit(conduit); err != nil {
			return
		}
	}

	return
}
