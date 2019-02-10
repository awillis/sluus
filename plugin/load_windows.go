// +build windows

package plugin

import "sync"

var (
	WindowsRegistry = NewRegistry()
)

func New(name string, pluginType Type) (plug Interface, err error) {
	factory := WindowsRegistry.Get(name)
	return factory(pluginType)
}

type Registry struct {
	sync.Mutex
	reg map[string]func(Type) (Interface, error)
}

func NewRegistry() *Registry {
	return &Registry{
		reg: make(map[string]func(Type) (Interface, error)),
	}
}

func (r *Registry) Register(name string, factory func(Type) (Interface, error)) {
	r.reg[name] = factory
}

func (r *Registry) Get(name string) func(Type) (Interface, error) {
	r.Lock()
	defer r.Unlock()
	return r.reg[name]
}
