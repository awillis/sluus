package processor

import (
	"time"

	ring "github.com/Workiva/go-datastructures/queue"
)

func (s *Sluus) Configure(opts ...Option) (err error) {
	for _, o := range opts {
		err = o(s)
		if err != nil {
			return
		}
	}
	return
}

func Input(input *ring.RingBuffer) Option {
	return func(s *Sluus) (err error) {
		s.input = input
		return
	}
}

func Output(output *ring.RingBuffer) Option {
	return func(s *Sluus) (err error) {
		s.output = output
		return
	}
}

func Reject(reject *ring.RingBuffer) Option {
	return func(s *Sluus) (err error) {
		s.reject = reject
		return
	}
}

func Accept(accept *ring.RingBuffer) Option {
	return func(s *Sluus) (err error) {
		s.accept = accept
		return
	}
}

func PollInterval(duration time.Duration) Option {
	return func(s *Sluus) (err error) {
		if duration < time.Second {
			duration = time.Second
		}
		s.pollInterval = duration
		return
	}
}

func BatchSize(size uint) Option {
	return func(s *Sluus) (err error) {
		s.batchSize = size
		return
	}
}
