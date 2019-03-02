package kafka

import (
	"context"
	"sync"

	"github.com/segmentio/kafka-go"

	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
)

var _ plugin.Interface = new(Sink)
var _ plugin.Consumer = new(Sink)

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

func (s *Sink) Initialize() (err error) {
	return
}

func (s *Sink) Start(ctx context.Context) {
	return
}

func (s *Sink) Consume(batch *message.Batch) {
	return
}

func (s *Sink) Shutdown() (err error) {
	return
}
