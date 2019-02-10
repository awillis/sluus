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

func TestSink_Initialize(t *testing.T) {
	sink, err := New(plugin.SINK)
	assert.Nil(t, err)
	err = sink.Initialize()
	assert.Nil(t, err)
}

//func TestSink_Execute(t *testing.T) {
//	sink, err := New(plugin.SINK)
//	assert.Nil(t, err)
//	err = sink.Execute()
//	assert.Nil(t, err)
//}

func TestSink_Shutdown(t *testing.T) {
	proc, err := New(plugin.SINK)
	assert.Nil(t, err)
	assert.Implements(t, (*plugin.Consumer)(nil), proc)
	sink := proc.(plugin.Consumer)
	err = sink.Shutdown()
	assert.Nil(t, err)
}
