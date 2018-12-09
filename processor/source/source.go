package source

import (
	"sluus/plugin"
)

type Source struct {
	plugin.PlugBase
}

type Processor interface {
	Test()
}

func (s *Source) PluginInit() bool {
	return true
}

func (s *Source) Test() {

}
