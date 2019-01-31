package noop

import (
	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
)

var _ plugin.Interface = new(Sink)
var _ plugin.Consumer = new(Sink)

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

func (s *Sink) Consume() (ch chan *message.Batch) {
	return
}

func (s *Sink) Shutdown() (err error) {
	return
}
