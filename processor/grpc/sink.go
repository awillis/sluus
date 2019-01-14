package grpc

import (
	"github.com/awillis/sluus/plugin"
	"net"
)

type Sink struct {
	plugin.Base
	Config SinkConfig
}

type SinkConfig struct {
	ListenAddr net.Addr
	CommonConfig
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

func (s *SinkConfig) Validate() (err error) {
	return
}

func (s *SinkConfig) Configure() (err error) {
	return
}
