// +build windows

package plugin

import "sync"

var Registry = NewRegistry()

func New(name string, pluginType Type) (plug Interface, err error) {
	factory := Registry.Get(name)
	return factory(pluginType)
}

type registry struct {
	sync.Mutex
	reg map[string]func(Type) (Interface, error)
}

func NewRegistry() *registry {
	return &registry{
		reg: make(map[string]func(Type) (Interface, error)),
	}
}

func (r *registry) Add(name string, factory func(Type) (Interface, error)) {
	r.reg[name] = factory
}

func (r *registry) Get(name string) func(Type) (Interface, error) {
	r.Lock()
	defer r.Unlock()
	return r.reg[name]
}
