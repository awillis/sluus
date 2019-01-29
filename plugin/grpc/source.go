package grpc

import (
	"context"
	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
)

type Source struct {
	plugin.Base
	opts options
}

func (s *Source) Options() interface{} {
	return &s.opts
}

func (s *Source) Initialize(ctx context.Context) (err error) {
	return
}

func (s *Source) Process(message.Batch) (batch message.Batch, err error) {
	return batch
}

func (s *Source) Shutdown() (err error) {
	return
}
