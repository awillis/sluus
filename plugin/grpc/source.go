package grpc

import (
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

func (s *Source) Initialize() (err error) {
	return
}

func (s *Source) Execute(input <-chan message.Batch, accept chan<- message.Batch, reject chan<- message.Batch) (err error) {
	return
}

func (s *Source) Shutdown() (err error) {
	return
}
