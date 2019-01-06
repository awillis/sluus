package pipeline

import (
	"github.com/awillis/sluus/processor"
)

type Sluus struct {
	sender   *processor.Processor
	receiver *processor.Processor
	counter  int64
}

func NewSluus(sender, receiver *processor.Processor) *Sluus {
	sluus := new(Sluus)
	sluus.sender = sender
	sluus.receiver = receiver
	return sluus
}

func (s *Sluus) Connect() {
	select {
	case item, ok := <-s.sender.Output():
		if !ok {
			s.sender.Logger.Error("output channel closed")
		}
		s.receiver.Input() <- item
		s.counter++
	}
}
