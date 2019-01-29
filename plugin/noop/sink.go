package noop

import (
	"context"
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

func (s *Sink) Initialize(ctx context.Context) (err error) {
	return
}

func (s *Sink) Process(message.Batch) (batch message.Batch, err error) {
	return batch
}

func (s *Sink) Shutdown() (err error) {
	return
}
