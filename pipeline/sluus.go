package pipeline

import (
	"runtime"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/awillis/sluus/plugin"
	"github.com/awillis/sluus/processor"
)

type Sluus struct {
	Id      string
	wg      *sync.WaitGroup
	logger  *zap.SugaredLogger
	flume   *processor.Flume
	counter int64
}

func NewSluus() (sluus *Sluus) {
	sluus = new(Sluus)
	sluus.Id = uuid.New().String()
	sluus.wg = new(sync.WaitGroup)
	sluus.flume = new(processor.Flume)
	return
}

func (s *Sluus) Run() {

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			s.wg.Add(1)
		shutdown:
			for {
				select {
				case item, ok := <-s.Flume().Output():
					if !ok {
						s.Logger().Error("output channel closed")
						break shutdown
					} else {
						s.Flume().Input() <- item
						s.counter++
					}
				}
			}
			s.wg.Done()
		}()
	}
}

func (s *Sluus) ID() string {
	return s.Id
}
func (s *Sluus) Type() plugin.Type {
	return plugin.CONDUIT
}

func (s *Sluus) Options() interface{} {
	return nil
}

func (s *Sluus) Flume() *processor.Flume {
	return s.flume
}

func (s *Sluus) Logger() *zap.SugaredLogger {
	return s.logger.With("sluus", s.ID())
}

func (s *Sluus) SetLogger(logger *zap.SugaredLogger) {
	s.logger = logger
}
