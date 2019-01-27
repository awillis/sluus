package kafka

import (
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

func (s *Sink) Initialize() (err error) {
	return
}

func (s *Sink) Execute(input <-chan message.Batch, accept chan<- message.Batch, reject chan<- message.Batch) (err error) {
	return
}

func (s *Sink) Shutdown() (err error) {
	return
}
