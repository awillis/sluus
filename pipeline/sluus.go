package pipeline

import (
	"github.com/awillis/sluus/core"
	"github.com/awillis/sluus/plugin"
	"github.com/awillis/sluus/processor"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Sluus struct {
	Id       string
	logger   *zap.SugaredLogger
	sender   *processor.Processor
	receiver *processor.Processor
	counter  int64
}

func NewSluus(sender, receiver *processor.Processor) *Sluus {
	sluus := new(Sluus)
	sluus.Id = uuid.New().String()
	sluus.sender = sender
	sluus.receiver = receiver
	return sluus
}

func (s *Sluus) Connect() {
	select {
	case item, ok := <-s.Output():
		if !ok {
			s.Logger().Error("output channel closed")
		}
		s.Input() <- item
		s.counter++
	}
}

func (s Sluus) ID() string {
	return s.Id
}
func (s Sluus) Type() plugin.Type {
	return plugin.CONDUIT
}

func (s Sluus) Input() chan<- core.Batch {
	return s.receiver.Input()
}

func (s Sluus) Output() <-chan core.Batch {
	return s.sender.Output()
}

func (s Sluus) SetLogger(logger *zap.SugaredLogger) {
	s.logger = logger
}

func (s Sluus) Logger() *zap.SugaredLogger {
	return s.logger.With("sluus", s.ID())
}
