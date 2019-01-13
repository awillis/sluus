package grpc

import "github.com/awillis/sluus/plugin"

type Source struct {
	plugin.Base
}

func (s *Source) Version() string {
	return string(s.Major) + "." + string(s.Minor) + "." + string(s.Patch)
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
