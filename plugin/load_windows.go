// +build windows

package plugin

import (
	"sync"
)

var (
	WinPlugReg *sync.Map
)

func init() {
	WinPlugReg = new(sync.Map)
}

type (
	WindowsPluginRegistry sync.Map
	pConstructor          func(Type) (Interface, error)
	iConstructor          func(Type) (Interface, error)
)

func NewProcessor(name string, pluginType Type) (procInt Interface, err error) {
	var factory pConstructor
	if f, ok := WinPlugReg.Load(name); ok {
		factory = f.(pConstructor)
	}
	return factory(pluginType)
}

func NewMessage(name string) (plugInt Interface, err error) {
	var factory iConstructor
	if f, ok := WinPlugReg.Load(name); ok {
		factory = f.(iConstructor)
	}
	return factory(MESSAGE)
}
