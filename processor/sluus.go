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
		inCtr, outCtr     uint64
		queue             *queue
		gate              map[uint64]*gate
		inputwg, outputwg *workGroup
		cnf               *config
	}
)

func newSluus(cnf *config) (sluus *Sluus) {

	sluus = &Sluus{
		cnf:      cnf,
		queue:    newQueue(cnf),
		gate:     make(map[uint64]*gate),
		inputwg:  new(workGroup),
		outputwg: new(workGroup),
	}

	if cnf.pluginType == plugin.SINK {
		sluus.gate[INPUT] = newGate()
	} else {
		sluus.gate[OUTPUT] = newGate()
	}
	return
}

func (s *Sluus) Initialize() (err error) {
	return s.queue.Initialize()
}

func (s *Sluus) Start() {

	s.queue.Start()

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
		s.Logger().Infof("%s input shutdown signalled", plugin.TypeName(s.cnf.pluginType))
		for s.Input().Len() > 0 {
			batch := s.Input().Get()
			if batch != nil {
				s.queue.Put(INPUT, batch)
			}
		}
		break loop
	default:
		batch := s.Input().Poll(ctx, s.cnf.pollInterval)

		// if there's no data received
		if batch != nil && batch.Count() > 0 {
			s.Logger().Infof("%s input queue len: %d", plugin.TypeName(s.cnf.pluginType), s.Input().Len())
			s.Logger().Infof("%s input batch: %d", plugin.TypeName(s.cnf.pluginType), batch.Count())
			s.queue.Put(INPUT, batch)
		} else {
			s.Logger().Infof("%s received nil batch", plugin.TypeName(s.cnf.pluginType))
			time.Sleep(s.cnf.pollInterval)
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
		s.Logger().Infof("%s output shutdown signalled", plugin.TypeName(s.cnf.pluginType))
		break
	case <-ticker.C:
		goto loop
	case batch, ok := <-s.queue.Output():
		if s.cnf.pluginType == plugin.SOURCE {
			s.Logger().Infof("source batch output: %d", batch.Count())
			s.Logger().Infof("source output ring len: %d", s.Output().Len())
		}

		if ok {
		retryOutput:
			if s.Output().Len() < 128 {
				s.Output().Put(batch)
			} else {
				time.Sleep(time.Second)
				goto retryOutput
			}
		}
		goto loop
	case batch, ok := <-s.queue.Accept():
		if ok {
		retryAccept:
			if s.Accept().Len() < 128 {
				s.Accept().Put(batch)
			} else {
				time.Sleep(time.Second)
				goto retryAccept
			}

		}
		goto loop
	case batch, ok := <-s.queue.Reject():
		if ok {
		retryReject:
			if s.Reject().Len() < 128 {
				s.Reject().Put(batch)
			} else {
				time.Sleep(time.Second)
				goto retryReject
			}

		}
		goto loop
	}
}

func (s *Sluus) shutdown() {

	// shutdown input then output
	if s.cnf.pluginType != plugin.SOURCE {
		s.Logger().Info("sluus input shutdown")
		s.inputwg.Shutdown()
	}

	if s.cnf.pluginType != plugin.SINK {
		s.Logger().Info("sluus output shutdown")
		s.outputwg.Shutdown()
	}

	if err := s.queue.shutdown(); err != nil {
		s.Logger().Error(err)
	}
}
