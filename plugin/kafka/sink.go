package kafka

import (
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

func (s *Sink) Initialize() (err error) {
	return
}

func (s *Sink) Execute() (err error) {
	return
}

func (s *Sink) Shutdown() (err error) {
	return
}
