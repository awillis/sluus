package noop

import "github.com/awillis/sluus/plugin"

type Sink struct {
	plugin.Base
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
