package pipeline

import (
	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
	"github.com/awillis/sluus/processor"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var ErrNoReceiver = errors.New("no receiver")

type Sluus struct {
	Id          string
	logger      *zap.SugaredLogger
	sender      processor.Interface
	receiver    processor.Interface
	hasReceiver bool
	counter     int64
}

func NewSluus(sender processor.Interface) (sluus *Sluus) {
	sluus = new(Sluus)
	sluus.Id = uuid.New().String()
	sluus.sender = sender
	return
}

func (s *Sluus) Connect() (err error) {

	if !s.hasReceiver {
		return ErrNoReceiver
	}

	select {
	case item, ok := <-s.Output():
		if !ok {
			s.Logger().Error("output channel closed")
		}
		s.Input() <- item
		s.counter++
	}
	return
}

func (s *Sluus) ID() string {
	return s.Id
}
func (s *Sluus) Type() plugin.Type {
	return plugin.CONDUIT
}

func (s *Sluus) Options() interface{} {
	return nil
}

func (s *Sluus) Input() chan<- message.Batch {
	return s.receiver.Input()
}

func (s *Sluus) Output() <-chan message.Batch {
	return s.sender.Output()
}

func (s *Sluus) Logger() *zap.SugaredLogger {
	return s.logger.With("sluus", s.ID())
}

func (s *Sluus) SetLogger(logger *zap.SugaredLogger) {
	s.logger = logger
}

func (s *Sluus) SetReceiver(receiver processor.Interface) {
	if s.receiver == nil {
		s.receiver = receiver
		s.hasReceiver = true
	}
}
