package processor

import (
	"github.com/awillis/sluus/plugin"
)

type Source struct {
	plugin.Plugin
}

func (s *Source) Run() {

}

func (s *Source) Execute() error {
	var err error
	return err
}

func (s *Source) Shutdown() error {
	var err error
	return err
}
