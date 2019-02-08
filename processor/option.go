package processor

import (
	"github.com/dgraph-io/badger/options"
	"time"

	ring "github.com/Workiva/go-datastructures/queue"
)

func Input(input *ring.RingBuffer) SluusOpt {
	return func(s *Sluus) (err error) {
		s.ring[INPUT] = input
		return
	}
}

func Output(output *ring.RingBuffer) SluusOpt {
	return func(s *Sluus) (err error) {
		s.ring[OUTPUT] = output
		return
	}
}

func Reject(reject *ring.RingBuffer) SluusOpt {
	return func(s *Sluus) (err error) {
		s.ring[REJECT] = reject
		return
	}
}

func Accept(accept *ring.RingBuffer) SluusOpt {
	return func(s *Sluus) (err error) {
		s.ring[ACCEPT] = accept
		return
	}
}

func PollInterval(duration time.Duration) SluusOpt {
	return func(s *Sluus) (err error) {
		if duration < time.Second {
			duration = time.Second
		}
		s.pollInterval = duration
		return
	}
}

func BatchSize(size uint) SluusOpt {
	return func(s *Sluus) (err error) {
		s.batchSize = size
		return
	}
}

func DataDir(path string) QueueOpt {
	return func(q *Queue) (err error) {
		q.opts.Dir = path
		q.opts.ValueDir = path
		return
	}
}

func TableLoadingMode(mode string) QueueOpt {
	return func(q *Queue) (err error) {
		q.opts.TableLoadingMode = options.FileIO
		return
	}
}

func ValueLogLoadingMode(mode string) QueueOpt {
	return func(q *Queue) (err error) {
		q.opts.ValueLogLoadingMode = options.FileIO
		return
	}
}
