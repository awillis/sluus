package noop

import (
	"context"
	"encoding/json"
	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
	"runtime"
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
		Count     uint64    `json:"count"`
	}
)

func (s *Source) Options() interface{} {
	return s.opts
}

func (s *Source) Initialize() (err error) {
	plugin.Validate(s.opts,
		s.opts.validMessagePerBatch(),
		s.opts.validBatchInterval(),
	)
	return
}

func (s *Source) Start(ctx context.Context) {

	go func(ctx context.Context) {
		s.wg.Add(1)
		defer s.wg.Done()
		interval := time.Duration(s.opts.BatchInterval) * time.Second
		timer := time.NewTimer(interval)

	loop:
		select {
		case <-ctx.Done():
			break
		case <-timer.C:
			batch := message.NewBatch(s.opts.MessagePerBatch)
			for i := 0; uint64(i) < s.opts.MessagePerBatch; i++ {

				content, err := json.Marshal(&noopMsg{
					Timestamp: time.Now(),
					Count:     uint64(i),
				})

				if err != nil {
					s.Logger().Error(err)
				}

				msg, err := message.WithContent(content)

				if err != nil {
					s.Logger().Error(err)
				}

				if err := batch.Add(msg); err != nil {
					s.Logger().Error(err)
				}
			}
			s.output <- batch
			timer.Reset(interval)
			goto loop
		default:
			runtime.Gosched()
		}
	}(ctx)
}

func (s *Source) Produce() <-chan *message.Batch {
	return s.output
}

func (s *Source) Shutdown() (err error) {
	s.wg.Wait()
	return
}
