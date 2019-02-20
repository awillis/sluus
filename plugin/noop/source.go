package noop

import (
	"context"
	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
)

var _ plugin.Interface = new(Source)
var _ plugin.Producer = new(Source)

type Source struct {
	plugin.Base
	opts *options
}

func (s *Source) Options() interface{} {
	return s.opts
}

func (s *Source) Initialize() (err error) {
	plugin.Validate(s.opts,
		s.opts.validMessagePerBatch(),
	)
	s.Logger().Infof("message per batch: %d", s.opts.MessagePerBatch)
	return
}

func (s *Source) Start(ctx context.Context) {

}

func (s *Source) Produce() <-chan *message.Batch {
	return make(chan *message.Batch)
}

func (s *Source) Shutdown() (err error) {
	return
}
