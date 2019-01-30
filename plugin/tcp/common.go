package tcp

import (
	"github.com/awillis/sluus/plugin"
	"github.com/google/uuid"
)

const (
	NAME  string = "tcp"
	MAJOR uint8  = 0
	MINOR uint8  = 0
	PATCH uint8  = 1
)

type options struct {
	// TCP port number to listen on
	port int
	// OS socket buffer size, a portion of which will be allocated for the app
	sockBufferSize int
	// application buffer size
	readBufferSize int
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
	default:
		return plug, plugin.ErrUnimplemented
	}
}
