package grpc

import "github.com/awillis/sluus/plugin"

type Sink struct {
	plugin.Base
}

//func (s *Sink) ID() string {
//	return s.Base.ID()
//}
//
//func (s *Sink) Name() string {
//	return s.Base.Name()
//}
//
//func (s *Sink) Type() plugin.Type {
//	return s.Base.Type()
//}
//
//func (s *Sink) Version() string {
//	return s.Base.Version()
//}
func (s *Sink) Version() string {
	return string(s.Major) + "." + string(s.Minor) + "." + string(s.Patch)
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
