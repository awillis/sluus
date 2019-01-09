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

func New(ptype plugin.Type) (plugin.Interface, error) {

	var plug plugin.Interface

	switch ptype {
	case plugin.SINK:
		plug = &Sink{
			Base: plugin.Base{
				Id:       uuid.New().String(),
				PlugName: "kafkaSink",
				Major:    MAJOR,
				Minor:    MINOR,
				Patch:    PATCH,
			},
		}
	case plugin.SOURCE:
		plug = &Source{
			Base: plugin.Base{
				Id:       uuid.New().String(),
				PlugName: "kafkaSource",
				Major:    MAJOR,
				Minor:    MINOR,
				Patch:    PATCH,
			},
		}
	default:
		return nil, plugin.ErrUnimplemented
	}

	return plug, nil
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
