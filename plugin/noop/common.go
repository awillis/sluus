package noop

import (
	"github.com/awillis/sluus/plugin"
	"github.com/google/uuid"
)

const (
	NAME  string = "noop"
	MAJOR uint8  = 0
	MINOR uint8  = 0
	PATCH uint8  = 1
)

type options struct {
}

func New(pluginType plugin.Type) (plug plugin.Interface, err error) {

	switch pluginType {
	case plugin.CONDUIT:
		return &Sink{
			Base: plugin.Base{
				Id:       uuid.New().String(),
				PlugName: NAME,
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
				PlugName: NAME,
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
