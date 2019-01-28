package tcp

import (
	"context"
	"github.com/awillis/sluus/plugin"
	"net"
	"sync"
)

type Source struct {
	plugin.Base
	opts      options
	wg        *sync.WaitGroup
	sock      *net.TCPListener
	start     chan *net.TCPConn
	end       chan *net.TCPConn
	conntable *sync.Map
}

func (s *Source) Options() interface{} {
	return &s.opts
}

func (s *Source) Initialize(ctx context.Context) (err error) {
	return
}

func (s *Source) Execute() (err error) {
	return
}

func (s *Source) Shutdown() (err error) {
	return
}
