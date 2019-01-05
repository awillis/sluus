package processor

import (
	"github.com/awillis/sluus/plugin"
)

type Sink struct {
	plugin.Plugin
}

func (s *Sink) Run() {

}

func (s *Sink) Execute() error {
	var err error
	return err
}

func (s *Sink) Shutdown() error {
	var err error
	return err
}
