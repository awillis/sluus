package grpc

import (
	"github.com/awillis/sluus/plugin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigurePort(t *testing.T) {
	sink := new(Sink)

	err := plugin.Configure(sink, sink.opts.Port(42))
	assert.Equal(t, 42, sink.opts.port, "port is correctly set")
	assert.Nil(t, err, "no errors")

	err = plugin.Configure(sink, sink.opts.Port(-1))
	assert.Errorf(t, err, plugin.ErrInvalidOption.Error())
}
