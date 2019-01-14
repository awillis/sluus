// +build windows

package plugin

import (
	"sync"

	"github.com/awillis/sluus/processor/grpc"
	"github.com/awillis/sluus/processor/kafka"
	"github.com/awillis/sluus/processor/noop"
	"github.com/awillis/sluus/processor/tcp"
)

var (
	winPlugReg WindowsPluginRegistry
)

func init() {
	winPlugReg = new(WindowsPluginRegistry)
	winPlugReg["grpc"] = grpc.New()
	winPlugReg["kafka"] = kafka.New()
	winPlugReg["noop"] = noop.New()
	winPlugReg["tcp"] = tcp.New()
}

type WindowsPluginRegistry sync.Map
type pConstructor func(Type) (Processor, error)
type iConstructor func(Type) (Interface, error)

func NewProcessor(name string, pluginType Type) (procInt Processor, err error) {
	var factory pConstructor = winPlugReg[name]
	return factory(name)(pluginType)
}

func NewMessage(name string) (plugInt Interface, err error) {
	var factory iConstructor = winPlugReg[name]
	return factory(name)(MESSAGE)
}
