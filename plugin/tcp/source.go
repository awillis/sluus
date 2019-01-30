package tcp

import (
	"net"
	"sync"

	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
)

var _ plugin.Interface = new(Source)
var _ plugin.Producer = new(Source)

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

func (s *Source) Initialize() (err error) {
	return
}

func (s *Source) Produce() (ch chan message.Batch) {
	return
}

func (s *Source) Shutdown() (err error) {
	return
}
