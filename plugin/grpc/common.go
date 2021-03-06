package grpc

import (
	"github.com/awillis/sluus/plugin"
	"github.com/google/uuid"
)

const (
	NAME  string = "grpc"
	MAJOR uint8  = 0
	MINOR uint8  = 0
	PATCH uint8  = 1
)

type options struct {
	Port uint64 `toml:"port"`
}

func New(pluginType plugin.Type) (plug plugin.Interface, err error) {

	switch pluginType {
	case plugin.SOURCE:
		return &Source{
			Base: plugin.Base{
				Id:       uuid.New().String(),
				PlugName: NAME,
				PlugType: pluginType,
				Major:    MAJOR,
				Minor:    MINOR,
				Patch:    PATCH,
			},
		}, err
	case plugin.CONDUIT:
		return &Conduit{
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
