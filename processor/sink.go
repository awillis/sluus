package processor

import (
	"github.com/awillis/sluus/plugin"
)

type Sink struct {
	sink
	plugin.Plugin
}

type sink interface {
	Consume()
}

func (s *Sink) Consume() {

}
