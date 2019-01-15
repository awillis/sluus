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
	Port int `mapstructure:"port"`
	// OS socket buffer size
	SockBufferSize int `mapstructure:"sock_buffer_size"`
	// application buffer size, used to read from OS
	ReadBufferSize int `mapstructure:"read_buffer_size"`
}

func New(pluginType plugin.Type) (plug plugin.Processor, err error) {

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
