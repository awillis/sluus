package core

import (
	"github.com/awillis/sluus/plugin"
	"github.com/awillis/sluus/plugin/grpc"
	"github.com/awillis/sluus/plugin/kafka"
	"github.com/awillis/sluus/plugin/noop"
	"github.com/awillis/sluus/plugin/tcp"
)

func init() {
	plugin.Registry.Add("grpc", grpc.New)
	plugin.Registry.Add("kafka", kafka.New)
	plugin.Registry.Add("noop", noop.New)
	plugin.Registry.Add("tcp", tcp.New)
}
