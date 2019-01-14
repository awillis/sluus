package noop

import (
	"github.com/awillis/sluus/plugin"
	"github.com/google/uuid"
)

const (
	MAJOR uint8 = 0
	MINOR uint8 = 0
	PATCH uint8 = 1
)

func New(pluginType plugin.Type) (plug plugin.Processor, err error) {

	switch pluginType {
	case plugin.CONDUIT:
		return &Sink{
			Base: plugin.Base{
				Id:       uuid.New().String(),
				PlugName: "noopConduit",
				PlugType: pluginType,
				Major:    MAJOR,
				Minor:    MINOR,
				Patch:    PATCH,
			},
		}, err
	case plugin.SINK:
		return &Sink{
			Base: plugin.Base{
				Id:       uuid.New().String(),
				PlugName: "noopSink",
				PlugType: pluginType,
				Major:    MAJOR,
				Minor:    MINOR,
				Patch:    PATCH,
			},
		}, err
	default:
		return plug, plugin.ErrUnimplemented
	}
}

// Config contains common options for all plugin types
type CommonConfig struct {
	Test string
}
