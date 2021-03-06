package processor

import (
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger/options"
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

func Input(input *gate) Option {
	return func(p *Processor) (err error) {
		p.sluus.gate[INPUT] = input
		return
	}
}

func Output(output *gate) Option {
	return func(p *Processor) (err error) {
		p.sluus.gate[OUTPUT] = output
		return
	}
}

func Reject(reject *gate) Option {
	return func(p *Processor) (err error) {
		p.sluus.gate[REJECT] = reject
		return
	}
}

func Accept(accept *gate) Option {
	return func(p *Processor) (err error) {
		p.sluus.gate[ACCEPT] = accept
		return
	}
}

func PollInterval(duration time.Duration) Option {
	return func(p *Processor) (err error) {
		if duration < 100*time.Millisecond {
			duration = 100 * time.Millisecond
		}
		p.cnf.pollInterval = duration
		return
	}
}

func RingSize(size uint64) Option {
	return func(p *Processor) (err error) {
		if size == 0 {
			size = 128
		}
		p.cnf.ringSize = size
		return
	}
}

func BatchSize(size uint64) Option {
	return func(p *Processor) (err error) {
		if size == 0 {
			size = 64
		}
		p.cnf.batchSize = size
		return
	}
}

func BatchTimeout(duration time.Duration) Option {
	return func(p *Processor) (err error) {
		if duration < time.Second {
			duration = time.Second
		}
		p.cnf.batchTimeout = duration
		return
	}
}

func QueryQueueRequests(requests uint64) Option {
	return func(p *Processor) (err error) {
		if requests < 256 {
			requests = 256
		}
		p.cnf.qqRequests = requests
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
