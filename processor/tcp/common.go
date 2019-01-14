package tcp

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
	case plugin.SOURCE:
		return &Source{
			Base: plugin.Base{
				Id:       uuid.New().String(),
				PlugName: "tcpSource",
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
