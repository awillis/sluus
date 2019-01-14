package grpc

import "github.com/awillis/sluus/plugin"

type Source struct {
	plugin.Base
	Config SourceConfig
}

type SourceConfig struct {
	CommonConfig
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

func (s *SourceConfig) Validate() (err error) {
	return
}

func (s *SourceConfig) Configure() (err error) {
	return
}
