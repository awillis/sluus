package pipeline

import "sync"

type Registry struct {
	sync.Map
}

func NewRegistry() *Registry {
	return new(Registry)
}

func (r *Registry) AddPipeline(pipe *Pipe) {
	r.Store(pipe.ID(), pipe)
}
