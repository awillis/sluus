package tcp

import (
	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
	"github.com/google/uuid"
	"net"
	"sync"
)

const (
	NAME  string = "tcp"
	MAJOR uint8  = 0
	MINOR uint8  = 0
	PATCH uint8  = 1
)

type options struct {
	plugin.Option
	// TCP port number to listen on
	port int
	// batch size
	batchSize uint64
	// application buffer size
	bufferSize int
	// OS socket buffer size, a portion of which will be allocated for the app
	sockBufferSize int
}

func New(pluginType plugin.Type) (plug plugin.Interface, err error) {

	switch pluginType {
	case plugin.SOURCE:
		return &Source{
			opts:    new(options),
			wg:      new(sync.WaitGroup),
			batch:   make(chan *message.Batch),
			message: make(chan *message.Message),
			start:   make(chan *net.TCPConn),
			end:     make(chan *net.TCPConn),
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

// defaultPort() validates the port value given in the configuration
// file and sets a reasonable default if needed
func (o *options) defaultPort() plugin.Default {
	return func(def plugin.Option) {
		if o.port < 1 || o.port > 65535 {
			o.port = 3030
		}
	}
}

func (o *options) defaultBatchSize() plugin.Default {
	return func(def plugin.Option) {
		if o.batchSize < 64 {
			o.batchSize = 64
		}
	}
}

func (o *options) defaultBufferSize() plugin.Default {
	return func(def plugin.Option) {
		if o.bufferSize < 16384 {
			o.bufferSize = 16384
		}
	}
}

func (o *options) defaultSockBufferSize() plugin.Default {
	return func(def plugin.Option) {
		if o.sockBufferSize < 65536 {
			o.sockBufferSize = 65536
		}
	}
}
