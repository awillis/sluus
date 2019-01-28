package grpc

import (
	"context"
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

func (s *Source) Execute() (err error) {
	return
}

func (s *Source) Shutdown() (err error) {
	return
}
