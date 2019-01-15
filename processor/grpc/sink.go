package grpc

import (
	"github.com/awillis/sluus/plugin"
	"reflect"
)

type Sink struct {
	plugin.Base
	conf Config
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

// Allows port to be set for sink
func (c Config) Port(port int) plugin.Option {
	return func(p plugin.Processor) (err error) {
		s := reflect.ValueOf(p).Interface().(*Sink)
		s.conf.port = port
		return
	}
}
