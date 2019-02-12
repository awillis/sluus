package noop

import (
	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
)

var _ plugin.Interface = new(Source)
var _ plugin.Producer = new(Source)

type Source struct {
	plugin.Base
	opts options
}

func (s *Source) Options() interface{} {
	return &s.opts
}

func (s *Source) Initialize() (err error) {
	return
}

func (s *Source) Produce() (batch *message.Batch, err error) {
	batch = new(message.Batch)
	return
}

func (s *Source) Shutdown() (err error) {
	return
}
