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
		ringSize            uint64
		inCtr, outCtr       uint64
		pollInterval        time.Duration
		pType               plugin.Type
		queue               *queue
		ring                map[uint64]*ring.RingBuffer
		logger              *zap.SugaredLogger
		poll, input, output *workGroup
	}
	workGroup struct {
		sync.WaitGroup
		cancel context.CancelFunc
	}
)

func newSluus(pType plugin.Type) (sluus *Sluus) {
	return &Sluus{
		pType:  pType,
		queue:  newQueue(pType),
		ring:   make(map[uint64]*ring.RingBuffer),
		poll:   new(workGroup),
		input:  new(workGroup),
		output: new(workGroup),
	}
}

func (s *Sluus) Initialize() (err error) {

	if s.pType == plugin.SINK {
		s.ring[INPUT] = ring.NewRingBuffer(s.ringSize)
	}

	s.ring[OUTPUT] = ring.NewRingBuffer(s.ringSize)
	s.ring[REJECT] = ring.NewRingBuffer(s.ringSize)
	s.ring[ACCEPT] = ring.NewRingBuffer(s.ringSize)
	return s.queue.Initialize()
}

func (s *Sluus) Start() {

	s.queue.Start()

	poll, pCancel := context.WithCancel(context.Background())
	s.poll.cancel = pCancel

	in, iCancel := context.WithCancel(context.Background())
	s.input.cancel = iCancel

	out, oCancel := context.WithCancel(context.Background())
	s.output.cancel = oCancel

	if s.pType != plugin.SOURCE {
		go s.ioInput(in)
		go s.ioPoll(poll, INPUT)
	}

	go s.ioOutput(out)
	go s.ioPoll(poll, OUTPUT)
	go s.ioPoll(poll, REJECT)
	go s.ioPoll(poll, ACCEPT)
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

func (s *Sluus) ioInput(ctx context.Context) {

	s.input.Add(1)
	defer s.input.Done()

	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

loop:
	select {
	case <-ctx.Done():
		break
	case <-ticker.C:
		if s.Input().IsDisposed() {
			break
		}

		b, err := s.Input().Poll(s.pollInterval)

		if err == ring.ErrDisposed {
			break
		}

		if batch, ok := b.(*message.Batch); ok && batch.Count() > 0 {
			s.queue.Put(INPUT, batch)
		}
		goto loop
	}
}

func (s *Sluus) ioOutput(ctx context.Context) {

	s.output.Add(1)
	defer s.output.Done()

	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

loop:
	select {
	case <-ctx.Done():
		break
	case <-ticker.C:
		goto loop
	case batch, ok := <-s.queue.Output():
		if ok {
			if e := s.Output().Put(batch); e != nil {
				s.Logger().Error(e)
			}
		}
		goto loop
	case batch, ok := <-s.queue.Accept():
		if ok {
			if e := s.Accept().Put(batch); e != nil {
				s.Logger().Error(e)
			}
		}
		goto loop
	case batch, ok := <-s.queue.Reject():
		if ok {
			if e := s.Reject().Put(batch); e != nil {
				s.Logger().Error(e)
			}
		}
		goto loop
	}
}

func (s *Sluus) ioPoll(ctx context.Context, prefix uint64) {

	s.poll.Add(1)
	defer s.poll.Done()

	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-ctx.Done():
			close(s.queue.requestChan[prefix])
			break loop
		case <-ticker.C:
			if s.ring[prefix].Len() < s.ring[prefix].Cap() {
				s.queue.requestChan[prefix] <- true
			}

			goto loop
		}
	}
}

func (s *Sluus) shutdown() {

	// shutdown polling threads, then input, then output
	s.poll.Shutdown()
	s.input.Shutdown()
	s.output.Shutdown()

	if err := s.queue.shutdown(); err != nil {
		s.Logger().Error(err)
	}
}

func (t *workGroup) Shutdown() {
	t.cancel()
	t.Wait()
}
