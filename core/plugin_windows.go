package core

import (
	"github.com/awillis/sluus/plugin"
	"github.com/awillis/sluus/plugin/grpc"
	"github.com/awillis/sluus/plugin/kafka"
	"github.com/awillis/sluus/plugin/noop"
	"github.com/awillis/sluus/plugin/tcp"
)

func init() {
	plugin.WindowsRegistry.Register("grpc", grpc.New)
	plugin.WindowsRegistry.Register("kafka", kafka.New)
	plugin.WindowsRegistry.Register("noop", noop.New)
	plugin.WindowsRegistry.Register("tcp", tcp.New)
}
