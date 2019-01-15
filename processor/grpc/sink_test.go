package grpc

import (
	"github.com/awillis/sluus/plugin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigurePort(t *testing.T) {
	sink := new(Sink)
	_ = plugin.Configure(sink, sink.conf.Port(42))
	assert.Equal(t, 42, sink.conf.port, "port is correctly set")
}
