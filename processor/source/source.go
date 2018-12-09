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

func (s *Source) Test() {

}
