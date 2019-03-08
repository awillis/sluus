package processor

import (
	"github.com/dgraph-io/badger/options"
	"path/filepath"
	"time"

	ring "github.com/Workiva/go-datastructures/queue"
)

type Option func(*Processor) error

func (p *Processor) Configure(opts ...Option) (err error) {
	for _, o := range opts {
		err = o(p)
		if err != nil {
			return
		}
	}
	return
}

func Input(input *ring.RingBuffer) Option {
	return func(p *Processor) (err error) {
		p.sluus.ring[INPUT] = input
		return
	}
}

func Output(output *ring.RingBuffer) Option {
	return func(p *Processor) (err error) {
		p.sluus.ring[OUTPUT] = output
		return
	}
}

func Reject(reject *ring.RingBuffer) Option {
	return func(p *Processor) (err error) {
		p.sluus.ring[REJECT] = reject
		return
	}
}

func Accept(accept *ring.RingBuffer) Option {
	return func(p *Processor) (err error) {
		p.sluus.ring[ACCEPT] = accept
		return
	}
}

func PollInterval(duration time.Duration) Option {
	return func(p *Processor) (err error) {
		if duration < 100*time.Millisecond {
			duration = 100 * time.Millisecond
		}
		p.pollInterval = duration
		p.sluus.pollInterval = duration
		p.sluus.queue.pollInterval = duration
		return
	}
}

func RingSize(size uint64) Option {
	return func(p *Processor) (err error) {
		if size == 0 {
			size = 128
		}
		p.sluus.ringSize = size
		return
	}
}

func BatchSize(size uint64) Option {
	return func(p *Processor) (err error) {
		if size == 0 {
			size = 64
		}
		p.sluus.queue.batchSize = size
		return
	}
}

func BatchTimeout(duration time.Duration) Option {
	return func(p *Processor) (err error) {
		if duration < time.Second {
			duration = time.Second
		}
		p.sluus.queue.batchTimeout = duration
		return
	}
}

func QueryInFlight(depth uint64) Option {
	return func(p *Processor) (err error) {
		if depth < 10 {
			depth = 10
		}
		p.sluus.queue.numInFlight = depth
		return
	}
}

func DataDir(path string) Option {
	return func(p *Processor) (err error) {
		path = filepath.Clean(path)
		p.sluus.queue.opts.Dir = path
		p.sluus.queue.opts.ValueDir = path
		return
	}
}

func TableLoadingMode(mode string) Option {
	return func(p *Processor) (err error) {
		switch mode {
		case "file":
			p.sluus.queue.opts.TableLoadingMode = options.FileIO
		case "memory":
			p.sluus.queue.opts.TableLoadingMode = options.LoadToRAM
		case "mmap":
			p.sluus.queue.opts.TableLoadingMode = options.MemoryMap
		default:
			p.sluus.queue.opts.TableLoadingMode = options.LoadToRAM
		}
		return
	}
}

func ValueLogLoadingMode(mode string) Option {
	return func(p *Processor) (err error) {
		switch mode {
		case "file":
			p.sluus.queue.opts.ValueLogLoadingMode = options.FileIO
		case "memory":
			p.sluus.queue.opts.ValueLogLoadingMode = options.LoadToRAM
		case "mmap":
			p.sluus.queue.opts.ValueLogLoadingMode = options.MemoryMap
		default:
			p.sluus.queue.opts.ValueLogLoadingMode = options.FileIO
		}
		return
	}
}
