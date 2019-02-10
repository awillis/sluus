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
		if duration < 3*time.Second {
			duration = 3 * time.Second
		}
		s.pollInterval = duration
		return
	}
}

func BatchSize(size uint64) SluusOpt {
	return func(s *Sluus) (err error) {
		if size == 0 {
			size = 64
		}
		s.batchSize = size
		return
	}
}

func RingSize(size uint64) SluusOpt {
	return func(s *Sluus) (err error) {
		if size == 0 {
			size = 128
		}
		s.ringSize = size
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
		switch mode {
		case "file":
			q.opts.TableLoadingMode = options.FileIO
		case "memory":
			q.opts.TableLoadingMode = options.LoadToRAM
		case "mmap":
			q.opts.TableLoadingMode = options.MemoryMap
		default:
			q.opts.TableLoadingMode = options.LoadToRAM
		}
		return
	}
}

func ValueLogLoadingMode(mode string) QueueOpt {
	return func(q *Queue) (err error) {
		switch mode {
		case "file":
			q.opts.ValueLogLoadingMode = options.FileIO
		case "memory":
			q.opts.ValueLogLoadingMode = options.LoadToRAM
		case "mmap":
			q.opts.ValueLogLoadingMode = options.MemoryMap
		default:
			q.opts.ValueLogLoadingMode = options.FileIO
		}
		return
	}
}
