package noop

import (
	"context"
	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
)

var _ plugin.Interface = new(Sink)
var _ plugin.Consumer = new(Sink)

type Sink struct {
	plugin.Base
	opts *options
}

func (s *Sink) Options() interface{} {
	return s.opts
}

func (s *Sink) Initialize() (err error) {
	plugin.Validate(s.opts,
		s.opts.defaultMessagePerBatch(),
		s.opts.defaultBatchInterval(),
	)
	s.opts.logCurrentConfig(s.Logger())
	return
}

func (s *Sink) Start(ctx context.Context) {
	return
}

func (s *Sink) Consume(batch *message.Batch) {
	for msg := range batch.Iter() {
		s.Logger().Infof("final message: %s", msg.String())
	}
	return
}

func (s *Sink) Shutdown() (err error) {
	return
}
