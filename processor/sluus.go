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
		ringSize      uint64
		inCtr, outCtr uint64
		pollInterval  time.Duration
		wg            *sync.WaitGroup
		queue         *Queue
		ring          map[uint64]*ring.RingBuffer
		logger        *zap.SugaredLogger
	}

	SluusOpt func(*Sluus) error
)

func NewSluus() (sluus *Sluus) {
	return &Sluus{
		wg:    new(sync.WaitGroup),
		queue: NewQueue(),
		ring:  make(map[uint64]*ring.RingBuffer),
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

// ring buffers must be initialized early
func (s *Sluus) RingInit() {
	s.ring[INPUT] = ring.NewRingBuffer(s.ringSize)
	s.ring[OUTPUT] = ring.NewRingBuffer(s.ringSize)
	s.ring[REJECT] = ring.NewRingBuffer(s.ringSize)
	s.ring[ACCEPT] = ring.NewRingBuffer(s.ringSize)
}

func (s *Sluus) Initialize() (err error) {
	return s.queue.Initialize()
}

func (s *Sluus) Start() {
	go s.inputIO()
	go s.outputIO(OUTPUT)
	go s.outputIO(REJECT)
	go s.outputIO(ACCEPT)
}

func (s *Sluus) Logger() *zap.SugaredLogger {
	return s.logger.With("sluus")
}

func (s *Sluus) SetLogger(logger *zap.SugaredLogger) {
	s.queue.logger = logger
	s.logger = logger
}

// Input() is used during pipeling assembly
func (s *Sluus) Input() *ring.RingBuffer {
	return s.ring[INPUT]
}

// Output() is used during pipeling assembly
func (s *Sluus) Output() *ring.RingBuffer {
	return s.ring[OUTPUT]
}

// Reject() is used during pipeling assembly
func (s *Sluus) Reject() *ring.RingBuffer {
	return s.ring[REJECT]
}

// Accept() is used during pipeling assembly
func (s *Sluus) Accept() *ring.RingBuffer {
	return s.ring[ACCEPT]
}

// Queue() is used during pipeling assembly
func (s *Sluus) Queue() *Queue {
	return s.queue
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

func (s *Sluus) receive(prefix, size uint64) (batch *message.Batch) {
	timer := time.NewTimer(s.pollInterval)
	c := s.queue.Get(prefix, size)
	batch = message.NewBatch(0)

	select {
	case <-timer.C:
		s.queue.Cancel()
	default:
		batch = <-c
	}
	return
}

func (s *Sluus) receiveInput() (batch *message.Batch) {
	return s.receive(INPUT, 0)
}

func (s *Sluus) send(prefix uint64, batch *message.Batch) {
	s.queue.Put(prefix, batch)
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
	r := s.Input()

	for {
		if r.IsDisposed() {
			break
		}

		input, err := r.Poll(s.pollInterval)
		if err != nil && err != ring.ErrTimeout {
			s.Logger().Error(err)
			continue
		}

		if batch, ok := input.(*message.Batch); ok {
			s.queue.Put(INPUT, batch)
		}
		runtime.Gosched()
	}
}

func (s *Sluus) outputIO(prefix uint64) {
	s.wg.Add(1)
	defer s.wg.Done()
	r := s.ring[prefix]

	for {
		if r.IsDisposed() {
			break
		}

		batch := s.receive(prefix, r.Cap())

		if batch.Count() > 0 {
			if e := r.Put(batch); e != nil {
				s.Logger().Error(e)
			}
		}
		runtime.Gosched()
	}
}
