package processor

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
)

type (
	Sluus struct {
		inCtr, outCtr             uint64
		queue                     *queue
		gate                      map[uint64]*gate
		pollwg, inputwg, outputwg *workGroup
		cnf                       *config
	}
)

func newSluus(cnf *config) (sluus *Sluus) {

	sluus = &Sluus{
		cnf:      cnf,
		queue:    newQueue(cnf),
		gate:     make(map[uint64]*gate),
		pollwg:   new(workGroup),
		inputwg:  new(workGroup),
		outputwg: new(workGroup),
	}

	if cnf.pluginType == plugin.SINK {
		sluus.gate[INPUT] = newGate()
	}

	sluus.gate[OUTPUT] = newGate()

	return
}

func (s *Sluus) Initialize() (err error) {

	//for _, direction := range compass[s.cnf.pluginType] {
	//	s.gate[direction] = newGate()
	//}

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
func (s *Sluus) Input() *gate {
	return s.gate[INPUT]
}

// Output() is used during pipeline assembly
func (s *Sluus) Output() *gate {
	return s.gate[OUTPUT]
}

// Reject() is used during pipeline assembly
func (s *Sluus) Reject() *gate {
	return s.gate[REJECT]
}

// Accept() is used during pipeline assembly
func (s *Sluus) Accept() *gate {
	return s.gate[ACCEPT]
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

	defer s.Logger().Infof("ioInput exit for %s", plugin.TypeName(s.cnf.pluginType))

loop:
	select {
	case <-ctx.Done():
		break
	case <-ticker.C:

		batch := s.Input().Poll(s.cnf.pollInterval)

		if batch != nil && batch.Count() > 0 {
			s.Logger().Infof("%s input queue len: %d, cap: %d", plugin.TypeName(s.cnf.pluginType))
			s.Logger().Infof("%s input batch: %d", plugin.TypeName(s.cnf.pluginType), batch.Count())
			println(s.Input())
			s.queue.Put(INPUT, batch)
		} else {
			s.Logger().Infof("%s received nil batch", plugin.TypeName(s.cnf.pluginType))
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
		if s.cnf.pluginType == plugin.SOURCE {
			s.Logger().Infof("source batch output: %d", batch.Count())
			s.Logger().Infof("source output ring len: %d", s.Output().Len())
			println(s.Output())
		}

		if ok {
			s.Output().Put(batch)
		}
		goto loop
	case batch, ok := <-s.queue.Accept():
		if ok {
			s.Accept().Put(batch)
		}
		goto loop
	case batch, ok := <-s.queue.Reject():
		if ok {
			s.Reject().Put(batch)
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
			//if s.gate[direction].Len() < 8192 {
			s.queue.requestChan[direction] <- true
			//}

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
