package sink

import (
	"sluus/plugin"
)

type Sink struct {
	plugin.PlugBase
}

type Consumer interface {
	Consume()
}

func (s *Sink) Consume() {

}
