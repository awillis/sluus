// +build windows

package plugin

import (
	"sync"
)

var (
	Registry *StaticRegistry
)

type StaticRegistry struct {
	sync.Mutex
	reg map[string]func(Type) (Interface, error)
}

func init() {
	Registry = new(StaticRegistry)
}

func New(name string, pluginType Type) (plug Interface, err error) {
	factory := Registry.Get(name)
	return factory(pluginType)
}

func (r *StaticRegistry) Register(name string, factory func(Type) (Interface, error)) {
	r.reg[name] = factory
}

func (r *StaticRegistry) Get(name string) func(Type) (Interface, error) {
	r.Lock()
	defer r.Unlock()
	return r.reg[name]
}
