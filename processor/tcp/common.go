package tcp

import (
	"github.com/awillis/sluus/plugin"
	"github.com/google/uuid"
)

var MAJOR uint8 = 0
var MINOR uint8 = 0
var PATCH uint8 = 1

func New(ptype plugin.Type) (plugin.Processor, error) {

	switch ptype {
	case plugin.SOURCE:
		return &Source{
			Base: plugin.Base{
				Id:       uuid.New().String(),
				PlugName: "tcpSource",
				PlugType: ptype,
				Major:    MAJOR,
				Minor:    MINOR,
				Patch:    PATCH,
			},
		}, nil
	default:
		return nil, plugin.ErrUnimplemented
	}
}
