package processor

import (
	"context"
	"sync"
	"time"

	ring "github.com/Workiva/go-datastructures/queue"
	"go.uber.org/zap"

	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
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

func (s *Sluus) Initialize() (err error) {
	s.ring[INPUT] = ring.NewRingBuffer(s.ringSize)
	s.ring[OUTPUT] = ring.NewRingBuffer(s.ringSize)
	s.ring[REJECT] = ring.NewRingBuffer(s.ringSize)
	s.ring[ACCEPT] = ring.NewRingBuffer(s.ringSize)
	return s.queue.Initialize()
}

func (s *Sluus) Start(ctx context.Context) {
	s.queue.Start(ctx)
	go s.ioThread(ctx, OUTPUT)
	go s.ioThread(ctx, REJECT)
	go s.ioThread(ctx, ACCEPT)
}

func (s *Sluus) Logger() *zap.SugaredLogger {
	return s.logger.With("sluus")
}

func (s *Sluus) SetLogger(logger *zap.SugaredLogger) {
	s.queue.logger = logger
	s.logger = logger
}

// Input() is used during pipeline assembly
func (s *Sluus) Input() *ring.RingBuffer {
	return s.ring[INPUT]
}

// Output() is used during pipeline assembly
func (s *Sluus) Output() *ring.RingBuffer {
	return s.ring[OUTPUT]
}

// Reject() is used during pipeline assembly
func (s *Sluus) Reject() *ring.RingBuffer {
	return s.ring[REJECT]
}

// Accept() is used during pipeline assembly
func (s *Sluus) Accept() *ring.RingBuffer {
	return s.ring[ACCEPT]
}

// receiveInput() is used by the processor runner
func (s *Sluus) receiveInput() <-chan *message.Batch {
	return s.queue.Input()
}

// sendOutput() is used by the processor runner
func (s *Sluus) sendOutput(batch *message.Batch) {
	s.send(OUTPUT, batch)
}

// sendReject() is used by the processor runner
func (s *Sluus) sendReject(batch *message.Batch) {
	s.send(REJECT, batch)
}

// sendAccept() is used by the processor runner
func (s *Sluus) sendAccept(batch *message.Batch) {
	s.send(ACCEPT, batch)
}

// send() is used by the processor runner
func (s *Sluus) send(prefix uint64, batch *message.Batch) {
	s.queue.Put(prefix, batch)
}

func (s *Sluus) ioThread(ctx context.Context, prefix uint64) {
	s.wg.Add(1)
	defer s.wg.Done()
	input := make(chan *message.Batch)

	go func(ctx context.Context, input chan *message.Batch) {
		// input ring to input queue
	loop:
		select {
		case <-ctx.Done():
			break
		default:
			b, err := s.Input().Poll(s.pollInterval)

			if err == ring.ErrDisposed {
				break
			}

			if batch, ok := b.(*message.Batch); ok {
				input <- batch
			}
			goto loop
		}

	}(ctx, input)

	go func(ctx context.Context, s *Sluus) {
		// output queue to output rings
		timer := time.NewTimer(s.pollInterval)

		select {
		case <-ctx.Done():
			break
		case <-timer.C:
			for _, typ := range []uint64{OUTPUT, ACCEPT, REJECT} {
				size := s.ring[typ].Len()
				if size > 0 {
					s.queue.requestChan[typ] <- size
				}
			}
			timer.Reset(s.pollInterval)
		}
	}(ctx, s)

loop:
	select {
	case <-ctx.Done():
		break
	case batch := <-input:
		s.queue.Put(INPUT, batch)
		goto loop
	case batch, ok := <-s.queue.Output():
		if ok {
			if e := s.ring[OUTPUT].Put(batch); e != nil {
				s.Logger().Error(e)
			}
		}
		goto loop
	case batch, ok := <-s.queue.Accept():
		if ok {
			if e := s.ring[ACCEPT].Put(batch); e != nil {
				s.Logger().Error(e)
			}
		}
		goto loop
	case batch, ok := <-s.queue.Reject():
		if ok {
			if e := s.ring[REJECT].Put(batch); e != nil {
				s.Logger().Error(e)
			}
		}
		goto loop
	}
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
