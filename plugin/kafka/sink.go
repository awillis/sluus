package kafka

import (
	"context"
	"github.com/awillis/sluus/message"
	"sync"

	"github.com/segmentio/kafka-go"

	"github.com/awillis/sluus/plugin"
)

type Sink struct {
	plugin.Base
	msgs   chan []byte
	wg     *sync.WaitGroup
	writer *kafka.Writer
	opts   options
}

func (s *Sink) Options() interface{} {
	return &s.opts
}

func (s *Sink) Initialize(ctx context.Context) (err error) {
	return
}

func (s *Sink) Process(message.Batch) (batch message.Batch, err error) {
	return batch
}

func (s *Sink) Shutdown() (err error) {
	return
}
