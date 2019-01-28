package grpc

import (
	"context"
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

func (s *Sink) Execute() (err error) {
	return
}

func (s *Sink) Shutdown() (err error) {
	return
}
