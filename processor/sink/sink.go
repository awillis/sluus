package sink

import (
	"github.com/awillis/sluus/plugin"
)

type Sink struct {
	plugin.PlugBase
}

type Consumer interface {
	Consume()
}

func (s *Sink) Consume() {

}
