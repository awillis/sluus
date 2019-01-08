package plugin

import (
	"runtime"
)

type Sink struct {
	Plugin
}

func (s *Sink) Run() {

}

func (s *Sink) Execute() error {
	runtime.LockOSThread()
	var err error
	return err
}

func (s *Sink) Shutdown() error {
	var err error
	return err
}
