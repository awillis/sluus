package noop

import (
	"context"
	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
	"sync"
	"time"
)

var _ plugin.Interface = new(Source)
var _ plugin.Producer = new(Source)

type (
	Source struct {
		plugin.Base
		opts   *options
		wg     *sync.WaitGroup
		output chan *message.Batch
	}
	noopMsg struct {
		Timestamp time.Time `json:"timestamp"`
		Counter   int       `json:"counter"`
	}
)

func (s *Source) Options() interface{} {
	return s.opts
}

func (s *Source) Initialize() (err error) {
	plugin.Validate(s.opts,
		s.opts.defaultMessagePerBatch(),
		s.opts.defaultBatchInterval(),
	)
	s.opts.logCurrentConfig(s.Logger())
	return
}

func (s *Source) Start(ctx context.Context) {

	go func(ctx context.Context) {
		//s.wg.Add(1)
		//defer s.wg.Done()
		ticker := time.NewTicker(time.Duration(s.opts.BatchInterval) * time.Millisecond)
		counter := 0

	loop:
		select {
		case <-ctx.Done():
			break
		case <-ticker.C:
			batch := message.NewBatch(s.opts.MessagePerBatch)
			for i := 0; uint64(i) < s.opts.MessagePerBatch; i++ {

				msg, err := message.New(&noopMsg{
					Timestamp: time.Now(),
					Counter:   counter,
				})

				if err != nil {
					s.Logger().Error(err)
				}

				if err := batch.Add(msg); err != nil {
					s.Logger().Error(err)
				}

				counter++
			}
			s.output <- batch
			goto loop
		}
	}(ctx)
}

func (s *Source) Produce() <-chan *message.Batch {
	return s.output
}

func (s *Source) Shutdown() (err error) {
	//s.wg.Wait()
	return
}
