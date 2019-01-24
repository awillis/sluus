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
	Config struct {
		Name       string
		Source     ProcessorConfig
		AcceptSink ProcessorConfig
		RejectSink ProcessorConfig
		Conduit    []ProcessorConfig
	}

	ProcessorConfig struct {
		Plugin  string                 `toml:"plugin"`
		Options map[string]interface{} `toml:"option"`
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

	yggtree, err := toml.LoadFile(filename)

	if err != nil {
		return
	}

	// source
	source := yggtree.Get("source")
	switch source.(type) {
	case *toml.Tree:
		if pc, e := configFromTree(source.(*toml.Tree)); e != nil {
			return config, e
		} else {
			config.Source = pc
		}
	default:
		pos := yggtree.Position()
		return config, errors.Wrapf(ErrConfigSection, "source at line %d, column %d", pos.Line, pos.Col)
	}

	// conduit
	conduit := yggtree.Get("conduit")
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
		pos := yggtree.Position()
		return config, errors.Wrapf(ErrConfigSection, "conduit at line %d, column %d", pos.Line, pos.Col)
	}

	// accept sink
	acceptSink := yggtree.Get("sink.accept")
	switch acceptSink.(type) {
	case *toml.Tree:
		if accept, e := configFromTree(acceptSink.(*toml.Tree)); e != nil {
			return config, e
		} else {
			config.AcceptSink = accept
		}
	default:
		pos := yggtree.Position()
		return config, errors.Wrapf(ErrConfigSection, "sink.accept at line %d, column %d", pos.Line, pos.Col)
	}

	// reject sink
	rejectSink := yggtree.Get("sink.reject")
	switch rejectSink.(type) {
	case *toml.Tree:
		if reject, e := configFromTree(rejectSink.(*toml.Tree)); e != nil {
			return config, e
		} else {
			config.RejectSink = reject
		}
	default:
		pos := yggtree.GetPosition("sink.reject")
		yggtree.Position()
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

	pc.Options = tree.ToMap()
	delete(pc.Options, "plugin")

	return
}
