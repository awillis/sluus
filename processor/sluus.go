package processor

import (
	ring "github.com/Workiva/go-datastructures/queue"
	"github.com/awillis/sluus/message"
	"go.uber.org/zap"
	"time"
)

type (
	Sluus struct {
		logger             *zap.SugaredLogger
		inCtr, outCtr      uint64
		queue              *Queue
		io, reject, accept *ring.RingBuffer
	}

	Option func(*Sluus) error
)

// io threads are built on a ring buffer
// thread looks at self and pushes to a queue
// work is pulled from a queue and pushes onto next ring buffer
// work must only proceed if there are slots available on output ring buffer
// output, accept and reject are all the same
// shutdown stops all work
// flush ensures all ring buffers have written out to next queue
// size of accept and reject ring buffers is
// S = size of output ring buffer
// N = number of Conduits between source and reject sink
// (N^2) + (S/2)
// (X^2) + (Y/2) = 80
// (y/2) = 80 - X^2
// Y = 40 - (X^2)/2
// Y = -(X^2)/2 + 40
// Given 8 conduits of size 32 = 80
// Given 10 conduits of size 64 = 132
// Given 20 conduits of size 128 = 464
// Given 10 conduits of size 12 = 106
// Given 3 conduits of size 96 = 57
// Given 2 conduits of size 8 = 8
// Size = Y, Number of conduits = X
// (Y^2) + (X/2)

func NewSluus() (sluus *Sluus) {
	return &Sluus{
		queue: NewQueue(""),
	}
}

func Initialize() (err error) {
	// setup queue and ring buffers
	return
}

func (s *Sluus) Logger() *zap.SugaredLogger {
	return s.logger
}

func (s *Sluus) SetLogger(logger *zap.SugaredLogger) {
	s.logger = logger
}

func (s *Sluus) Input() (batch *message.Batch) {
	if msg, err := s.input.Get(); err != nil {
		s.Logger().Error(err)
	} else {
		if m, ok := msg.(*message.Message); ok {
			if e := s.queue.Put(m); e != nil {
				s.Logger().Error(e)
			}
		}
	}
	return
}

func (s *Sluus) Output(batch *message.Batch) {
	// wire this to queue consume
	return s.output
}

func (s *Sluus) Reject() chan *message.Batch {
	return s.reject
}

func (s *Sluus) Accept() chan *message.Batch {
	return s.accept
}

func (s *Sluus) Flush() {

}

func RingIO(
	ring *ring.RingBuffer,
) func(s *Sluus) error {
	return func(s *Sluus) error {
		for {
			batch, e := ring.Poll(time.Millisecond)
			if e != nil {
				s.Logger().Error(e)
			}
			for m := range batch.(*message.Batch).Iter() {
				if e = s.queue.Put(m); e != nil {
					s.Logger().Error(e)
				}
			}
		}
	}
}
