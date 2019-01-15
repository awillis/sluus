package grpc

import (
	"github.com/awillis/sluus/plugin"
)

type Source struct {
	plugin.Base
	conf Config
}

func (s *Source) Initialize() (err error) {
	return
}

func (s *Source) Execute() (err error) {
	return
}

func (s *Source) Shutdown() (err error) {
	return
}
