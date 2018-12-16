package pipeline

import "sync"

type Registry struct {
	sync.Map
}

func NewRegistry() *Registry {
	return new(Registry)
}

func (r *Registry) AddPipeline(pipeline Pipeline) {
	r.Store(pipeline.ID(), &pipeline)
}
