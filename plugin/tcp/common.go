package tcp

import (
	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
	"github.com/google/uuid"
	"go.uber.org/zap"
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
	Port uint64 `toml:"port"`
	// batch size
	BatchSize uint64 `toml:"batch_size"`
	// application buffer size
	BufferSize uint64 `toml:"buffer_size"`
	// OS socket buffer size, a portion of which will be allocated for the app
	SockBufferSize uint64 `toml:"sock_buffer_size"`
	// poll interval, time between checks for new data to process in milliseconds
	PollInterval uint64 `toml:"poll_interval"`
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

func (o *options) logCurrentConfig(logger *zap.SugaredLogger) {
	logger.Infof("port: %d", o.Port)
	logger.Infof("batch size: %d", o.BatchSize)
	logger.Infof("buffer size: %d", o.BufferSize)
	logger.Infof("socket buffer size: %d", o.SockBufferSize)
	logger.Infof("poll interval: %d ms", o.PollInterval)
}

// defaultPort() validates the port value given in the configuration
// file and sets a reasonable default if needed
func (o *options) defaultPort() plugin.Default {
	return func(def plugin.Option) {
		if o.Port < 1 || o.Port > 65535 {
			o.Port = 3030
		}
	}
}

func (o *options) defaultBatchSize() plugin.Default {
	return func(def plugin.Option) {
		if o.BatchSize < 64 {
			o.BatchSize = 64
		}
	}
}

func (o *options) defaultBufferSize() plugin.Default {
	return func(def plugin.Option) {
		if o.BufferSize < 16384 {
			o.BufferSize = 16384
		}
	}
}

func (o *options) defaultSockBufferSize() plugin.Default {
	return func(def plugin.Option) {
		if o.SockBufferSize < 65536 {
			o.SockBufferSize = 65536
		}
	}
}

func (o *options) defaultPollInterval() plugin.Default {
	return func(def plugin.Option) {
		if o.PollInterval == 0 {
			o.PollInterval = 200
		}
	}
}
