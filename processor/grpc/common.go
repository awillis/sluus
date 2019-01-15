package grpc

import (
	"github.com/awillis/sluus/plugin"
	"github.com/google/uuid"
)

const (
	MAJOR uint8 = 0
	MINOR uint8 = 0
	PATCH uint8 = 1
)

type Config struct {
	port int
}

func New(pluginType plugin.Type) (plug plugin.Processor, err error) {

	switch pluginType {
	case plugin.SOURCE:
		return &Source{
			Base: plugin.Base{
				Id:       uuid.New().String(),
				PlugName: "grpcSource",
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
				PlugName: "grpcSink",
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
