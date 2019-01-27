package grpc

import (
	"github.com/awillis/sluus/plugin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateNewSink(t *testing.T) {
	_, err := New(plugin.SINK)
	assert.Nil(t, err, "no errors")
}

func TestConfigureSinkPort(t *testing.T) {
	sink := new(Sink)

	err := plugin.Configure(sink, sink.opts.Port(42))
	assert.Equal(t, 42, sink.opts.port, "port is correctly set")
	assert.Nil(t, err, "no errors")

	err = plugin.Configure(sink, sink.opts.Port(-1))
	assert.EqualError(t, plugin.ErrInvalidOption, err.Error())
}

func TestSink_Initialize(t *testing.T) {
	sink, err := New(plugin.SINK)
	assert.Nil(t, err)
	err = sink.Initialize(nil, nil, nil)
	assert.Nil(t, err)
}

func TestSink_Execute(t *testing.T) {
	sink, err := New(plugin.SINK)
	assert.Nil(t, err)
	err = sink.Execute()
	assert.Nil(t, err)
}

func TestSink_Shutdown(t *testing.T) {
	sink, err := New(plugin.SINK)
	assert.Nil(t, err)
	err = sink.Shutdown()
	assert.Nil(t, err)
}
