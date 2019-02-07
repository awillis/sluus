package processor

import (
	ring "github.com/Workiva/go-datastructures/queue"
	"github.com/awillis/sluus/message"
	"go.uber.org/zap"
	"sync"
	"time"
)

type (
	Sluus struct {
		logger                        *zap.SugaredLogger
		inCtr, outCtr                 uint64
		wg                            *sync.WaitGroup
		queue                         *Queue
		input, output, reject, accept *ring.RingBuffer
		pollInterval                  time.Duration
		batchSize                     uint
	}

	Option func(*Sluus) error
)

func NewSluus(proc Interface) (s *Sluus) {
	return &Sluus{
		wg:     new(sync.WaitGroup),
		queue:  NewQueue(proc.ID()),
		output: new(ring.RingBuffer),
	}
}

func (s *Sluus) Initialize() (err error) {
	go s.inputIO()
	return s.queue.Initialize()
}

func (s *Sluus) Logger() *zap.SugaredLogger {
	return s.logger
}

func (s *Sluus) SetLogger(logger *zap.SugaredLogger) {
	s.logger = logger
}

func (s *Sluus) Receive() (batch *message.Batch) {
	batch, err := s.queue.Get(s.batchSize)
	if err != nil {
		s.Logger().Error(err)
	}
	return
}

func (s *Sluus) send(ring *ring.RingBuffer, batch *message.Batch) {
	if err := ring.Put(batch); err != nil {
		s.Logger().Error(err)
	}
}

func (s *Sluus) Pass(batch *message.Batch) {
	s.send(s.output, batch)
}

func (s *Sluus) Reject(batch *message.Batch) {
	s.send(s.reject, batch)
}

func (s *Sluus) Accept(batch *message.Batch) {
	s.send(s.accept, batch)
}

func (s *Sluus) Shutdown() {
	if s.input != nil {
		s.input.Dispose()
	}

	if s.output != nil {
		s.output.Dispose()
	}

	s.wg.Wait()

	if err := s.queue.Shutdown(); err != nil {
		s.Logger().Error(err)
	}
}

func (s *Sluus) inputIO() {
	s.wg.Add(1)
	defer s.wg.Done()

	for {
		if s.input.IsDisposed() {
			break
		}

		input, err := s.input.Poll(s.pollInterval)
		if err != nil && err != ring.ErrTimeout {
			s.Logger().Error(err)
			continue
		}

		if batch, ok := input.(*message.Batch); ok {
			if e := s.queue.Put(batch); e != nil {
				s.Logger().Error(e)
			}
		}
	}
}

func (s *Sluus) Configure(opts ...Option) (err error) {
	for _, o := range opts {
		err = o(s)
		if err != nil {
			return
		}
	}
	return
}

func Input(input *ring.RingBuffer) Option {
	return func(s *Sluus) (err error) {
		s.input = input
		return
	}
}

func Output(output *ring.RingBuffer) Option {
	return func(s *Sluus) (err error) {
		s.output = output
		return
	}
}

func Reject(reject *ring.RingBuffer) Option {
	return func(s *Sluus) (err error) {
		s.reject = reject
		return
	}
}

func Accept(accept *ring.RingBuffer) Option {
	return func(s *Sluus) (err error) {
		s.accept = accept
		return
	}
}

func PollInterval(duration time.Duration) Option {
	return func(s *Sluus) (err error) {
		if duration < time.Second {
			duration = time.Second
		}
		s.pollInterval = duration
		return
	}
}

func BatchSize(size uint) Option {
	return func(s *Sluus) (err error) {
		s.batchSize = size
		return
	}
}
