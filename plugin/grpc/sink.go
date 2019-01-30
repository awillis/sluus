package grpc

import (
	"context"
	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
)

var _ plugin.Loader = new(Sink)
var _ plugin.Processor = new(Sink)

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

func (s *Sink) Process(message.Batch) (pass, reject, accept message.Batch, err error) {
	return
}

func (s *Sink) Shutdown() (err error) {
	return
}
