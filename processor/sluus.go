package processor

import (
	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/queue"
	"go.uber.org/zap"
)

type (
	Sluus struct {
		logger        *zap.SugaredLogger
		inputCounter  uint64
		outputCounter uint64
		input         chan *message.Batch
		output        chan *message.Batch
		reject        chan *message.Batch
		accept        chan *message.Batch
		databasePath  string
		queue         *queue.Queue
	}

	Option func(*Sluus) error
)

func NewSluus() (sluus *Sluus) {
	return &Sluus{
		queue: queue.New(""),
	}
}

func Initialize() (err error) {
	return
}

func (s *Sluus) Logger() *zap.SugaredLogger {
	return s.logger
}

func (s *Sluus) SetLogger(logger *zap.SugaredLogger) {
	s.logger = logger
}

func (s *Sluus) Input() chan *message.Batch {
	// wire this to queue produce
	// take messages from a batch and write
	return s.input
}

func (s *Sluus) Output() chan *message.Batch {
	// wire this to queue consume
	return s.output
}

func (s *Sluus) Reject() chan *message.Batch {
	return s.reject
}

func (s *Sluus) Accept() chan *message.Batch {
	return s.accept
}

// configuration options

func Configure(sluus *Sluus, opts ...Option) (err error) {
	for _, o := range opts {
		err = o(sluus)
		if err != nil {
			return
		}
	}
	return
}

func Input(input chan *message.Batch) Option {
	return func(s *Sluus) (err error) {
		s.input = input
		return
	}
}

func Output(output chan *message.Batch) Option {
	return func(s *Sluus) (err error) {
		s.output = output
		return
	}
}

func Reject(reject chan *message.Batch) Option {
	return func(s *Sluus) (err error) {
		s.reject = reject
		return
	}
}

func Accept(accept chan *message.Batch) Option {
	return func(s *Sluus) (err error) {
		s.accept = accept
		return
	}
}
