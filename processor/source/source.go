package source

import (
	"github.com/awillis/sluus/plugin"
)

type Source struct {
	plugin.PlugBase
}

type Producer interface {
	Produce()
}

func (s *Source) Produce() {

}
