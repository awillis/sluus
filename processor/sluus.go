package processor

import (
	"github.com/awillis/sluus/message"
)

type (
	Sluus struct {
		input  chan<- message.Batch
		output <-chan message.Batch
		reject <-chan message.Batch
		accept <-chan message.Batch
	}

	Option func(*Sluus) error
)

func (s *Sluus) Input() chan<- message.Batch {
	return s.input
}

func (s *Sluus) Output() <-chan message.Batch {
	return s.output
}

func (s *Sluus) Reject() <-chan message.Batch {
	return s.reject
}

func (s *Sluus) Accept() <-chan message.Batch {
	return s.accept
}

func Configure(sluus *Sluus, opts ...Option) (err error) {
	for _, o := range opts {
		err = o(sluus)
		if err != nil {
			return
		}
	}
	return
}

func Input(input chan<- message.Batch) Option {
	return func(s *Sluus) (err error) {
		s.input = input
		return
	}
}

func Output(output <-chan message.Batch) Option {
	return func(s *Sluus) (err error) {
		s.output = output
		return
	}
}

func Reject(reject <-chan message.Batch) Option {
	return func(s *Sluus) (err error) {
		s.reject = reject
		return
	}
}

func Accept(accept <-chan message.Batch) Option {
	return func(s *Sluus) (err error) {
		s.accept = accept
		return
	}
}
