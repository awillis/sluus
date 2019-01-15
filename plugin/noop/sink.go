package noop

import "github.com/awillis/sluus/plugin"

type Sink struct {
	plugin.Base
	opts *options
}

func (s *Sink) Initialize() (err error) {
	return
}

func (s *Sink) Execute() (err error) {
	return
}

func (s *Sink) Shutdown() (err error) {
	return
}
