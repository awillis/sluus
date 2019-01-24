package pipeline

import (
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

var ErrConfigValue = errors.New("unknown configuration value")
var ErrConfigSection = errors.New("unable to parse config section")

type (
	Config struct {
		Source  ProcessorConfig
		Sink    map[string]ProcessorConfig
		Conduit []ProcessorConfig
	}

	ProcessorConfig struct {
		Plugin  string                 `toml:"plugin"`
		Options map[string]interface{} `toml:"option"`
		Routes  []Route                `toml:",omitempty"`
	}

	Route struct {
		Terminate   bool   `toml:"terminate"`
		Destination string `toml:"destination"`
	}
)

func ReadConfigurationFile(filename string) (config Config, err error) {

	tree, err := toml.LoadFile(filename)

	if err != nil {
		return
	}

	// source
	source := tree.Get("source")
	if config.Source, err = procConfigFromTree(source.(*toml.Tree)); err != nil {
		return config, errors.Wrap(ErrConfigSection, "source")
	}

	// conduit
	conduit := tree.Get("conduit")
	switch conduit.(type) {
	case []*toml.Tree:
		for _, tree := range conduit.([]*toml.Tree) {
			if pc, err := procConfigFromTree(tree); err != nil {
				return config, err
			} else {
				config.Conduit = append(config.Conduit, pc)
			}
		}
	default:
		return config, errors.Wrap(ErrConfigSection, "conduit")
	}

	// sink
	sink := tree.Get("sink")
	switch sink.(type) {
	case *toml.Tree:
		config.Sink = make(map[string]ProcessorConfig)
		sinktree := sink.(*toml.Tree)
		for _, name := range sinktree.Keys() {
			if pc, err := procConfigFromTree(sinktree.Get(name).(*toml.Tree)); err != nil {
				return config, err
			} else {
				config.Sink[name] = pc
			}
		}
	default:
		return config, errors.Wrap(ErrConfigSection, "sink")
	}
	return
}

func procConfigFromTree(tree *toml.Tree) (pc ProcessorConfig, err error) {

	name := tree.Get("plugin")
	switch name.(type) {
	case string:
		pc.Plugin = name.(string)
	default:
		return pc, errors.Wrapf(ErrConfigValue, "found type %T for plugin name", name)
	}

	opt := tree.Get("option")
	switch opt.(type) {
	case *toml.Tree:
		pc.Options = opt.(*toml.Tree).ToMap()
	default:
		return pc, errors.Wrapf(ErrConfigValue, "found type %T for plugin options", opt)
	}

	routes := tree.Get("route")
	switch routes.(type) {
	case []*toml.Tree:
		for _, r := range routes.([]*toml.Tree) {
			var rte Route
			ttext, e := r.ToTomlString()

			if e != nil {
				return pc, e
			}

			if e := toml.Unmarshal([]byte(ttext), &rte); e != nil {
				return pc, errors.Wrapf(ErrConfigValue, "unable to decode: %s, %+v", e, ttext)
			}

			pc.Routes = append(pc.Routes, rte)
		}
	case nil:
		// routes are not always present
		return pc, err
	default:
		return pc, errors.Wrapf(ErrConfigValue, "found type %T for plugin routes", routes)
	}

	return
}
