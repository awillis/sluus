package kafka

import (
	"github.com/awillis/sluus/plugin"
	"github.com/google/uuid"
	"net"
	"strings"
)

var MAJOR uint8 = 0
var MINOR uint8 = 0
var PATCH uint8 = 1

func New(pluginType plugin.Type) (plug plugin.Interface, err error) {

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

func bootstrapLookup(endpoint string) ([]string, error) {

	var err error
	var brokers []string

	host, port, err := net.SplitHostPort(endpoint)
	if err != nil {
		return brokers, err
	}

	addrs, err := net.LookupHost(host)

	if err != nil {
		return brokers, err
	}

	for _, ip := range addrs {
		brokers = append(brokers, strings.Join([]string{ip, port}, ":"))
	}

	return brokers, err
}
