package processor

import (
	ring "github.com/Workiva/go-datastructures/queue"
	"github.com/awillis/sluus/message"
	"go.uber.org/zap"
	"runtime"
	"sync"
	"time"
)

type (
	Sluus struct {
		batchSize                         uint
		inCtr, outCtr                     uint64
		pollInterval                      time.Duration
		wg                                *sync.WaitGroup
		input, output, reject, accept     *ring.RingBuffer
		inputQ, outputQ, rejectQ, acceptQ *Queue
		logger                            *zap.SugaredLogger
	}

	Option func(*Sluus) error
)

func NewSluus(proc Interface) (s *Sluus) {
	return &Sluus{
		wg:      new(sync.WaitGroup),
		inputQ:  NewQueue(proc.ID(), "input"),
		outputQ: NewQueue(proc.ID(), "output"),
		rejectQ: NewQueue(proc.ID(), "reject"),
		acceptQ: NewQueue(proc.ID(), "accept"),
		output:  new(ring.RingBuffer),
	}
}

func (s *Sluus) Initialize() (err error) {
	outIO := s.outputIO(s.outputQ, s.output)
	rejIO := s.outputIO(s.rejectQ, s.reject)
	accIO := s.outputIO(s.acceptQ, s.accept)
	for i := 0; i < runtime.NumCPU(); i++ {
		go s.inputIO()
		go outIO()
		go rejIO()
		go accIO()
	}

	if e := s.inputQ.Initialize(); e != nil {
		return e
	}

	if e := s.outputQ.Initialize(); e != nil {
		return e
	}

	if e := s.rejectQ.Initialize(); e != nil {
		return e
	}

	if e := s.acceptQ.Initialize(); e != nil {
		return e
	}

	return
}

func (s *Sluus) Logger() *zap.SugaredLogger {
	return s.logger
}

func (s *Sluus) SetLogger(logger *zap.SugaredLogger) {
	s.logger = logger
}

// Input() is used by the pipeline to connect processors together
func (s *Sluus) Input() *ring.RingBuffer {
	return s.input
}

// Output() is used by the pipeline to connect processors together
func (s *Sluus) Output() *ring.RingBuffer {
	return s.output
}

// Reject() is used by the pipeline to connect processors together
func (s *Sluus) Reject() *ring.RingBuffer {
	return s.reject
}

// Accept() is used by the pipeline to connect processors together
func (s *Sluus) Accept() *ring.RingBuffer {
	return s.accept
}

func (s *Sluus) shutdown() {
	if s.input != nil {
		s.input.Dispose()
	}

	if s.output != nil {
		s.output.Dispose()
	}

	s.wg.Wait()

	if err := s.inputQ.shutdown(); err != nil {
		s.Logger().Error(err)
	}
}

func (s *Sluus) receive() (batch *message.Batch) {
	batch, err := s.inputQ.Get(s.batchSize)
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

func (s *Sluus) sendOutput(batch *message.Batch) {
	s.send(s.output, batch)
}

func (s *Sluus) sendReject(batch *message.Batch) {
	s.send(s.reject, batch)
}

func (s *Sluus) sendAccept(batch *message.Batch) {
	s.send(s.accept, batch)
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
			s.logger.Error(err)
			continue
		}

		if batch, ok := input.(*message.Batch); ok {
			if e := s.inputQ.Put(batch); e != nil {
				s.logger.Error(e)
			}
		}
	}
}

func (s *Sluus) outputIO(q *Queue, ring *ring.RingBuffer) func() {
	return func() {
		s.wg.Add(1)
		defer s.wg.Done()

		for {
			batch, err := q.Get(uint(ring.Cap()))
			if err != nil {
				s.logger.Error(err)
			}
			if batch.Count() > 0 {
				if e := ring.Put(batch); e != nil {
					s.logger.Error(e)
				}
			}
		}
	}
}
