package grpc

import (
	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
)

type Sink struct {
	plugin.Base
	opts options
}

func (s *Sink) Options() interface{} {
	return &s.opts
}

func (s *Sink) Initialize() (err error) {
	return
}

func (s *Sink) Execute(input <-chan message.Batch, accept chan<- message.Batch, reject chan<- message.Batch) (err error) {
	return
}

func (s *Sink) Shutdown() (err error) {
	return
}

// Allows port to be set for sink
func (c options) Port(port int) plugin.Option {
	return func(p plugin.Processor) (err error) {
		if port < 0 || port > 65535 {
			err = plugin.ErrInvalidOption
		}

		s := p.(*Sink)
		s.opts.port = port
		return
	}
}
