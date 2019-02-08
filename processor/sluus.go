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
		batchSize     uint
		inCtr, outCtr uint64
		pollInterval  time.Duration
		wg            *sync.WaitGroup
		queue         *Queue
		ring          map[byte]*ring.RingBuffer
		logger        *zap.SugaredLogger
	}

	SluusOpt func(*Sluus) error
)

func NewSluus(proc Interface) (s *Sluus) {
	return &Sluus{
		wg:    new(sync.WaitGroup),
		queue: NewQueue(),
		ring:  make(map[byte]*ring.RingBuffer),
	}
}

func (s *Sluus) Configure(opts ...SluusOpt) (err error) {
	for _, o := range opts {
		err = o(s)
		if err != nil {
			return
		}
	}
	return
}

func (s *Sluus) Initialize() (err error) {

	if e := s.queue.Initialize(); e != nil {
		return e
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		go s.inputIO()
		go s.outputIO(OUTPUT)
		go s.outputIO(REJECT)
		go s.outputIO(ACCEPT)
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
	return s.ring[INPUT]
}

// Output() is used by the pipeline to connect processors together
func (s *Sluus) Output() *ring.RingBuffer {
	return s.ring[OUTPUT]
}

// Reject() is used by the pipeline to connect processors together
func (s *Sluus) Reject() *ring.RingBuffer {
	return s.ring[REJECT]
}

// Accept() is used by the pipeline to connect processors together
func (s *Sluus) Accept() *ring.RingBuffer {
	return s.ring[ACCEPT]
}

func (s *Sluus) shutdown() {

	for i := range s.ring {
		if s.ring[i] != nil {
			s.ring[i].Dispose()
		}
	}

	s.wg.Wait()

	if err := s.queue.shutdown(); err != nil {
		s.Logger().Error(err)
	}
}

func (s *Sluus) receive() (batch *message.Batch) {
	batch, err := s.queue.Get(INPUT, s.batchSize)
	if err != nil {
		s.Logger().Error(err)
	}
	return
}

func (s *Sluus) send(prefix byte, batch *message.Batch) {
	if err := s.ring[prefix].Put(batch); err != nil {
		s.Logger().Error(err)
	}
}

func (s *Sluus) sendOutput(batch *message.Batch) {
	s.send(OUTPUT, batch)
}

func (s *Sluus) sendReject(batch *message.Batch) {
	s.send(REJECT, batch)
}

func (s *Sluus) sendAccept(batch *message.Batch) {
	s.send(ACCEPT, batch)
}

func (s *Sluus) inputIO() {
	s.wg.Add(1)
	defer s.wg.Done()
	r := s.ring[INPUT]

	for {
		if r.IsDisposed() {
			break
		}

		input, err := r.Poll(s.pollInterval)
		if err != nil && err != ring.ErrTimeout {
			s.logger.Error(err)
			continue
		}

		if batch, ok := input.(*message.Batch); ok {
			if e := s.queue.Put(INPUT, batch); e != nil {
				s.logger.Error(e)
			}
		}
	}
}

func (s *Sluus) outputIO(prefix byte) {
	s.wg.Add(1)
	defer s.wg.Done()
	r := s.ring[prefix]

	for {
		if r.IsDisposed() {
			break
		}

		batch, err := s.queue.Get(prefix, uint(r.Cap()))
		if err != nil {
			s.logger.Error(err)
		}
		if batch.Count() > 0 {
			if e := r.Put(batch); e != nil {
				s.logger.Error(e)
			}
		}
	}

}
