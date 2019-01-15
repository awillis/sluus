package grpc

import (
	"github.com/awillis/sluus/plugin"
	"reflect"
)

type Sink struct {
	port int
	plugin.Base
	conf SinkConfig
}

type SinkConfig bool
type SinkOption func(*Sink) error

func (s *Sink) Initialize(opts ...plugin.Option) (err error) {
	for _, o := range opts {
		oVal := reflect.ValueOf(o).Interface()
		err = oVal.(SinkOption)(s)
		if err != nil {
			return
		}
	}
	return
}

func (s *Sink) Execute() (err error) {
	return
}

func (s *Sink) Shutdown() (err error) {
	return
}

func (s SinkOption) Error() (err string) {
	return
}

// Allows port to be set for sink
func (sc *SinkConfig) Port(port int) SinkOption {
	return func(s *Sink) (err error) {
		s.port = port
		return
	}
}

func Test(sink *Sink, conf *SinkConfig, opts ...plugin.Option) error {
	return sink.Initialize(opts...)
}
