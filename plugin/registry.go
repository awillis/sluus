package plugin

import (
	"sync"
)

type Registry struct {
	sync.Map
}

func NewRegistry() *Registry {
	return new(Registry)
}

func (r *Registry) AddPlugin(proc Processor) {
	r.Store(proc.Name(), &proc)
}
