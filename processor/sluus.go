package processor

import (
	"context"
	"time"

	ring "github.com/Workiva/go-datastructures/queue"
	"go.uber.org/zap"

	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
)

type (
	Sluus struct {
		inCtr, outCtr             uint64
		queue                     *queue
		ring                      map[uint64]*ring.RingBuffer
		pollwg, inputwg, outputwg *workGroup
		cnf                       *config
	}
)

func newSluus(cnf *config) (sluus *Sluus) {
	return &Sluus{
		cnf:      cnf,
		queue:    newQueue(cnf),
		ring:     make(map[uint64]*ring.RingBuffer),
		pollwg:   new(workGroup),
		inputwg:  new(workGroup),
		outputwg: new(workGroup),
	}
}

func (s *Sluus) Initialize() (err error) {

	for _, direction := range compass[s.cnf.pluginType] {
		s.ring[direction] = ring.NewRingBuffer(s.cnf.ringSize)
	}

	return s.queue.Initialize()
}

func (s *Sluus) Start() {

	s.queue.Start()

	poll, pCancel := context.WithCancel(context.Background())
	s.pollwg.cancel = pCancel

	in, iCancel := context.WithCancel(context.Background())
	s.inputwg.cancel = iCancel

	out, oCancel := context.WithCancel(context.Background())
	s.outputwg.cancel = oCancel

	if s.cnf.pluginType != plugin.SOURCE {
		s.Logger().Info("start ioInput")
		go s.ioInput(in)
	}

	if s.cnf.pluginType != plugin.SINK {
		s.Logger().Info("start ioOutput")
		go s.ioOutput(out)
	}

	for _, direction := range compass[s.cnf.pluginType] {
		s.Logger().Infof("start ioPoll for %+v", direction)
		go s.ioPoll(poll, direction)
	}
}

func (s *Sluus) Logger() *zap.SugaredLogger {
	return s.cnf.logger.With("component", "sluus")
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
func (s *Sluus) send(direction uint64, batch *message.Batch) {
	s.queue.Put(direction, batch)
}

func (s *Sluus) ioInput(ctx context.Context) {

	s.inputwg.Add(1)
	defer s.inputwg.Done()

	ticker := time.NewTicker(s.cnf.pollInterval)
	defer ticker.Stop()

loop:
	select {
	case <-ctx.Done():
		break
	case <-ticker.C:
		if s.Input().IsDisposed() {
			break
		}

		b, err := s.Input().Poll(s.cnf.pollInterval)

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

	s.outputwg.Add(1)
	defer s.outputwg.Done()

	ticker := time.NewTicker(s.cnf.pollInterval)
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

func (s *Sluus) ioPoll(ctx context.Context, direction uint64) {

	s.pollwg.Add(1)
	defer s.pollwg.Done()

	ticker := time.NewTicker(s.cnf.pollInterval)
	defer ticker.Stop()

	defer close(s.queue.requestChan[direction])

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-ticker.C:
			if s.ring[direction].Len() < s.ring[direction].Cap() {
				s.queue.requestChan[direction] <- true
			}

			goto loop
		}
	}
}

func (s *Sluus) shutdown() {

	// shutdown polling threads, then inputwg, then outputwg
	s.Logger().Info("sluus poll shutdown")
	s.pollwg.Shutdown()
	s.Logger().Info("sluus input shutdown")
	s.inputwg.Shutdown()
	s.Logger().Info("sluus output shutdown")
	s.outputwg.Shutdown()

	if err := s.queue.shutdown(); err != nil {
		s.Logger().Error(err)
	}
}
