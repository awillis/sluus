package kafka

import (
	"github.com/segmentio/kafka-go"

	"github.com/awillis/sluus/plugin"
)

type Source struct {
	plugin.Base
	reader *kafka.Reader
	opts   *options
}

func (s *Source) Initialize() (err error) {
	return
}

func (s *Source) Execute() (err error) {
	return
}

func (s *Source) Shutdown() (err error) {
	return
}
