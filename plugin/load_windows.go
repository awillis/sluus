// +build windows

package plugin

import (
	"sync"
)

var (
	winPlugReg *sync.Map
)

func init() {
	winPlugReg = new(sync.Map)
}

type WindowsPluginRegistry sync.Map
type pConstructor func(Type) (Processor, error)
type iConstructor func(Type) (Interface, error)

func NewProcessor(name string, pluginType Type) (procInt Processor, err error) {
	var factory pConstructor
	if f, ok := winPlugReg.Load(name); ok {
		factory = f.(pConstructor)
	}
	return factory(pluginType)
}

func NewMessage(name string) (plugInt Interface, err error) {
	var factory iConstructor
	if f, ok := winPlugReg.Load(name); ok {
		factory = f.(iConstructor)
	}
	return factory(MESSAGE)
}
