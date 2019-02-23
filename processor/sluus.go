package processor

import (
	"context"
	"runtime"
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
	go s.ioThread(ctx)
}

func (s *Sluus) Logger() *zap.SugaredLogger {
	return s.logger.With("component", "sluus")
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
	s.Logger().Infof("sending %d output messages", batch.Count())
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

func (s *Sluus) ioThread(ctx context.Context) {
	s.wg.Add(1)
	defer s.wg.Done()

	ticker := time.NewTicker(time.Duration(s.pollInterval) * time.Millisecond)
	defer ticker.Stop()

	input := make(chan *message.Batch)

	go func(s *Sluus, ctx context.Context, input chan *message.Batch) {
		// input ring to input queue
		s.wg.Add(1)
		defer s.wg.Done()
		defer s.Logger().Info("end input ring thread")

		ticker := time.NewTicker(time.Duration(s.pollInterval) * time.Millisecond)
		defer ticker.Stop()

	loop:
		s.Logger().Info("inside input thread loop")
		select {
		case <-ctx.Done():
			break
		case <-ticker.C:
			s.Logger().Info("inside input timer loop")
			b, err := s.Input().Get()

			if err == ring.ErrDisposed {
				break
			}

			if batch, ok := b.(*message.Batch); ok {
				s.Logger().Infof("got batch of size %d from input ring", batch.Count())
				input <- batch
			}
			goto loop
		}

	}(s, ctx, input)

	go func(s *Sluus, ctx context.Context) {
		// output queue to output rings
		s.wg.Add(1)
		defer s.wg.Done()
		ticker := time.NewTicker(time.Duration(s.pollInterval) * time.Millisecond)
		defer ticker.Stop()

	loop:
		select {
		case <-ctx.Done():
			break
		case <-ticker.C:
			s.Logger().Info("output queue iothread")
			for _, typ := range []uint64{OUTPUT, ACCEPT, REJECT} {
				s.Logger().Infof("output ring cap %d len %d", s.ring[typ].Cap(), s.ring[typ].Len())
				size := s.ring[typ].Cap() - s.ring[typ].Len()
				if size > 0 {
					s.Logger().Infof("output queue request size %d", size)
					s.queue.requestChan[typ] <- size
				}
			}
			goto loop
		}
	}(s, ctx)

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
	case <-ticker.C:
		runtime.Gosched()
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
