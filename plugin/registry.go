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

func (r *Registry) AddPlugin(plugin Plugin) {
	r.Store(plugin.Name(), &plugin)
}