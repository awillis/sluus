package processor

import (
	"github.com/awillis/sluus/plugin"
)

type Source struct {
	source
	plugin.Plugin
}

type source interface {
	Produce()
}

func (s *Source) Produce() {

}
