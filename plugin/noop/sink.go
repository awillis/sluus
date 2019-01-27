package noop

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
