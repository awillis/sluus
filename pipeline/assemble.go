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

	source := processor.New(config.Source.Plugin, plugin.SOURCE)
	if err := source.Load(); err != nil {
		pipe.Logger().Errorw(err.Error(), "source", source.ID())
		return
	}

	if err := config.Source.Option.SetPluginOptions(source.Plugin().Options()); err != nil {
		pipe.Logger().Errorw(err.Error(), "source", source.ID())
		return
	}

	if err := pipe.AddSource(source); err != nil {
		pipe.Logger().Errorw(err.Error(), "source", source.ID())
		return
	} else {
		pipe.Logger().Infow("set source", "source", source.ID())
	}

	for _, conf := range config.Conduit {
		conduit := processor.New(conf.Plugin, plugin.CONDUIT)
		if err := conduit.Load(); err != nil {
			pipe.Logger().Errorw(err.Error(), "conduit", conduit.ID())

			return
		}

		if err := conf.Option.SetPluginOptions(conduit.Plugin().Options()); err != nil {
			pipe.Logger().Errorw(err.Error(), "conduit", conduit.ID())
			return
		}

		if err := pipe.AddConduit(conduit); err != nil {
			pipe.Logger().Errorw(err.Error(), "conduit", conduit.ID())
			return
		} else {
			pipe.Logger().Infow("added conduit", "id", conduit.ID())
		}
	}

	reject := processor.New(config.RejectSink.Plugin, plugin.SINK)
	if err := reject.Load(); err != nil {
		pipe.Logger().Errorw(err.Error(), "reject", reject.ID())
		return
	}

	if err := config.RejectSink.Option.SetPluginOptions(reject.Plugin().Options()); err != nil {
		pipe.Logger().Errorw(err.Error(), "reject", reject.ID())
		return
	}

	if err := pipe.AddReject(reject); err != nil {
		pipe.Logger().Errorw(err.Error(), "reject", reject.ID())
		return
	} else {
		pipe.Logger().Infow("add reject sink", "reject", reject.ID())
	}

	accept := processor.New(config.AcceptSink.Plugin, plugin.SINK)
	if err := accept.Load(); err != nil {
		pipe.Logger().Errorw(err.Error(), "accept", accept.ID())
		return
	}

	if err := config.AcceptSink.Option.SetPluginOptions(accept.Plugin().Options()); err != nil {
		pipe.Logger().Errorw(err.Error(), "accept", accept.ID())
		return
	}

	if err := pipe.AddAccept(accept); err != nil {
		pipe.Logger().Errorw(err.Error(), "accept", accept.ID())
		return
	} else {
		pipe.Logger().Infow("add accept sink", "accept", accept.ID())
	}

	return
}
