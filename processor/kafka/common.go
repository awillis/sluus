package kafka

import (
	"github.com/awillis/sluus/plugin"
	"github.com/google/uuid"
	"net"
)

var MAJOR uint8 = 0
var MINOR uint8 = 0
var PATCH uint8 = 1

func New(pluginType plugin.Type) (plug plugin.Processor, err error) {

	switch pluginType {
	case plugin.SINK:
		return &Sink{
			Base: plugin.Base{
				Id:       uuid.New().String(),
				PlugName: "kafkaSink",
				PlugType: pluginType,
				Major:    MAJOR,
				Minor:    MINOR,
				Patch:    PATCH,
			},
		}, err
	case plugin.SOURCE:
		return &Source{
			Base: plugin.Base{
				Id:       uuid.New().String(),
				PlugName: "kafkaSource",
				PlugType: pluginType,
				Major:    MAJOR,
				Minor:    MINOR,
				Patch:    PATCH,
			},
		}, err
	default:
		return plug, plugin.ErrUnimplemented
	}
}

func bootstrapLookup(endpoint string) (brokers []string, err error) {

	host, port, err := net.SplitHostPort(endpoint)
	if err != nil {
		return brokers, err
	}

	addrs, err := net.LookupHost(host)

	if err != nil {
		return brokers, err
	}

	for _, ip := range addrs {
		brokers = append(brokers, ip+":"+port)
	}

	return brokers, err
}
