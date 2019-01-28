package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"

	"github.com/awillis/sluus/plugin"
)

type Source struct {
	plugin.Base
	reader *kafka.Reader
	opts   options
}

func (s *Source) Options() interface{} {
	return &s.opts
}

func (s *Source) Initialize(ctx context.Context) (err error) {
	return
}

func (s *Source) Execute() (err error) {
	return
}

func (s *Source) Shutdown() (err error) {
	return
}
