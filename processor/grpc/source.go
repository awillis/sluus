package grpc

import (
	"github.com/awillis/sluus/plugin"
	"reflect"
)

type Source struct {
	port int
	plugin.Base
}

type SourceConfig bool
type SourceOption func(*Source) error

func (s *Source) Initialize(opts ...plugin.Option) (err error) {
	for _, o := range opts {
		oVal := reflect.ValueOf(o).Interface()
		err = oVal.(func(*Source) error)(s)
		if err != nil {
			return
		}
	}
	return
}

func (s *Source) Execute() (err error) {
	return
}

func (s *Source) Shutdown() (err error) {
	return
}

func (sc *SourceConfig) Port(port int) (opt plugin.Option) {
	var o SourceOption
	return o
}

func SourceTest(source *Source, conf *SourceConfig, opts ...plugin.Option) error {
	return source.Initialize(conf.Port(5))
}
