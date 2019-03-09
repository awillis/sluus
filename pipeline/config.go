package pipeline

import (
	"github.com/awillis/sluus/core"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"regexp"
)

var (
	ErrConfigValue   = errors.New("unknown configuration value")
	ErrConfigSection = errors.New("unable to parse config section")
	confPattern      = regexp.MustCompile(".pipe.toml$")
)

type (
	OptionMap map[string]interface{}

	Config struct {
		Name       string
		Pipe       PipeConfig
		Source     ProcessorConfig
		AcceptSink ProcessorConfig
		RejectSink ProcessorConfig
		Conduit    []ProcessorConfig
	}

	PipeConfig struct {
		PollInterval        uint64 `toml:"poll_interval"`
		BatchSize           uint64 `toml:"batch_size"`
		BatchTimeout        uint64 `toml:"batch_timeout"`
		RingSize            uint64 `toml:"ring_size"`
		QueryQueueRequests  uint64 `toml:"qq_requests"`
		TableLoadingMode    string `toml:"table_loading_mode"`
		ValueLogLoadingMode string `toml:"value_log_loading_mode"`
	}

	ProcessorConfig struct {
		Plugin string    `toml:"plugin"`
		Option OptionMap `toml:"option"`
	}
)

func FindConfigurationFiles() (files []string, err error) {

	err = filepath.Walk(core.CONFDIR, func(path string, info os.FileInfo, err error) (rerr error) {
		if info.IsDir() {
			return
		}

		if confPattern.MatchString(path) {
			files = append(files, path)
		}

		return
	})
	return
}

func ReadConfigurationFile(filename string) (config Config, err error) {

	tree, err := toml.LoadFile(filename)

	if err != nil {
		return
	}

	// name
	if name, ok := tree.Get("pipe.name").(string); ok {
		config.Name = name
	} else {
		return config, errors.Wrap(ErrConfigValue, "pipe.name")
	}

	// pipe
	pipe := tree.Get("pipe")
	switch pipe.(type) {
	case *toml.Tree:
		if pipeConf, ok := pipe.(*toml.Tree); ok {
			if e := pipeConf.Unmarshal(&config.Pipe); e != nil {
				return config, e
			}
		}
	default:
		pos := tree.Position()
		return config, errors.Wrapf(ErrConfigSection, "pipe at line %d, column %d", pos.Line, pos.Col)
	}

	// source
	source := tree.Get("source")
	switch source.(type) {
	case *toml.Tree:
		if pc, e := configFromTree(source.(*toml.Tree)); e != nil {
			return config, e
		} else {
			config.Source = pc
		}
	default:
		pos := tree.Position()
		return config, errors.Wrapf(ErrConfigSection, "source at line %d, column %d", pos.Line, pos.Col)
	}

	// conduit
	conduit := tree.Get("conduit")
	switch conduit.(type) {
	case []*toml.Tree:
		for _, conduitTree := range conduit.([]*toml.Tree) {
			if pc, e := configFromTree(conduitTree); e != nil {
				return config, e
			} else {
				config.Conduit = append(config.Conduit, pc)
			}
		}
	default:
		pos := tree.Position()
		return config, errors.Wrapf(ErrConfigSection, "conduit at line %d, column %d", pos.Line, pos.Col)
	}

	// accept sink
	acceptSink := tree.Get("sink.accept")
	switch acceptSink.(type) {
	case *toml.Tree:
		if accept, e := configFromTree(acceptSink.(*toml.Tree)); e != nil {
			return config, e
		} else {
			config.AcceptSink = accept
		}
	default:
		pos := tree.Position()
		return config, errors.Wrapf(ErrConfigSection, "sink.accept at line %d, column %d", pos.Line, pos.Col)
	}

	// reject sink
	rejectSink := tree.Get("sink.reject")
	switch rejectSink.(type) {
	case *toml.Tree:
		if reject, e := configFromTree(rejectSink.(*toml.Tree)); e != nil {
			return config, e
		} else {
			config.RejectSink = reject
		}
	default:
		pos := tree.GetPosition("sink.reject")
		tree.Position()
		return config, errors.Wrapf(ErrConfigSection, "sink.reject at line %d, column %d", pos.Line, pos.Col)
	}

	return
}

func configFromTree(tree *toml.Tree) (pc ProcessorConfig, err error) {

	name := tree.Get("plugin")
	switch name.(type) {
	case string:
		pc.Plugin = name.(string)
	default:
		pos := tree.Position()
		return pc, errors.Wrapf(ErrConfigValue, "found type %T for plugin name at line %d, column %d", name, pos.Line, pos.Col)
	}

	pc.Option = tree.ToMap()
	delete(pc.Option, "plugin")

	return
}
