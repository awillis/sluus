package pipeline

import "sync"

var Registry = NewRegistry()

type registry struct {
	sync.Mutex
	reg map[string]*Pipe
}

func NewRegistry() *registry {
	return &registry{
		reg: make(map[string]*Pipe),
	}
}

func (r *registry) Add(pipe *Pipe) {
	r.Lock()
	defer r.Unlock()
	r.reg[pipe.ID()] = pipe
}

func (r *registry) Start() {
	r.Lock()
	defer r.Unlock()
	for id, pipe := range r.reg {
		pipe.Logger().Infow("pipeline start", "id", id)
		pipe.Start()
	}
}

func (r *registry) Stop() {
	r.Lock()
	defer r.Unlock()
	for id, pipe := range r.reg {
		pipe.Logger().Infow("pipeline stop", "id", id)
		pipe.Stop()
	}
}
