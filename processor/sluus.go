package processor

import (
	"context"
	ring "github.com/Workiva/go-datastructures/queue"
	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
	"go.uber.org/zap"
	"sync"
	"time"
)

type (
	Sluus struct {
		ringSize      uint64
		inCtr, outCtr uint64
		pollInterval  time.Duration
		pType         plugin.Type
		wg            *sync.WaitGroup
		queue         *queue
		ring          map[uint64]*ring.RingBuffer
		logger        *zap.SugaredLogger
	}
)

func newSluus(pType plugin.Type) (sluus *Sluus) {
	return &Sluus{
		pType: pType,
		wg:    new(sync.WaitGroup),
		queue: newQueue(pType),
		ring:  make(map[uint64]*ring.RingBuffer),
	}
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
	go s.ioThread(OUTPUT)
	go s.ioThread(REJECT)
	go s.ioThread(ACCEPT)
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

// queue() is used during pipeling assembly
func (s *Sluus) Queue() *queue {
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
	return s.queue.Get(prefix, size)
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

func (s *Sluus) ioThread(ctx context.Context, prefix uint64) {
	s.wg.Add(1)
	defer s.wg.Done()
	shutdown := make(chan bool)
	input := make(chan *message.Batch)

	go func(chan *message.Batch) {
		// input ring to input queue
		for {
			b, err := s.Input().Poll(s.pollInterval)

			if err == ring.ErrDisposed {
				break
			}

			if batch, ok := b.(*message.Batch); ok {
				input <- batch
			}
		}
	}(input)

	go func(s *Sluus, ctx context.Context) {
		// output queue to output rings
		timer := time.NewTimer(s.pollInterval)

		select {
		case <-ctx.Done():

		case <-timer.C:
			for _, typ := range []uint64{INPUT, OUTPUT, ACCEPT, REJECT} {
				size := s.ring[typ].Len()
				if size > 0 {
					s.queue.requestChan[typ] <- size
				}
			}
		}
	}(s, ctx)

	select {
	case <-shutdown:
		break
	case batch := <-input:
		s.queue.Put(INPUT, batch)
	case batch, ok := <-s.queue.Output():
		if ok {
			if e := s.ring[OUTPUT].Put(batch); e != nil {
				s.Logger().Error(e)
			}
		} else {
			shutdown <- true
		}
	case batch, _ := <-s.queue.Accept():
		if e := s.ring[ACCEPT].Put(batch); e != nil {
			s.Logger().Error(e)
		}
	case batch, _ := <-s.queue.Reject():
		if e := s.ring[REJECT].Put(batch); e != nil {
			s.Logger().Error(e)
		}
	}
}
