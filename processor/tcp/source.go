package tcp

import (
	"github.com/awillis/sluus/plugin"
	"net"
	"sync"
)

type Source struct {
	plugin.Base
	wg        *sync.WaitGroup
	sock      *net.TCPListener
	start     chan *net.TCPConn
	end       chan *net.TCPConn
	conntable *sync.Map
}

type options struct {
	// TCP port number to listen on
	Port int `mapstructure:"port"`
	// OS socket buffer size
	SockBufferSize int `mapstructure:"sock_buffer_size"`
	// application buffer size, used to read from OS
	ReadBufferSize int `mapstructure:"read_buffer_size"`
}

func (s *Source) Initialize() error {
	var err error
	return err
}

func (s *Source) Execute() error {
	var err error
	return err
}

func (s *Source) Shutdown() error {
	var err error
	return err
}
