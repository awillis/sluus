// +build !windows

package main

import (
	"github.com/awillis/sluus/plugin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSink(t *testing.T) {
	plug, err := New(plugin.SINK)
	assert.NotNil(t, plug.ID())
	assert.Equal(t, "grpc", plug.Name())
	assert.NotNil(t, plug.Version())
	assert.Nil(t, err)
	assert.Equal(t, plug.Type(), plugin.SINK)
}

func TestNewSource(t *testing.T) {
	plug, err := New(plugin.SOURCE)
	assert.Nil(t, err)
	assert.Equal(t, plug.Type(), plugin.SOURCE)
}

func TestSinkInitialize(t *testing.T) {
	plug, err := New(plugin.SINK)
	assert.Nil(t, err)
	err = plug.Initialize()
	assert.Nil(t, err)
}

//func TestSinkExecute(t *testing.T) {
//	plug, err := New(plugin.SINK)
//	assert.Nil(t, err)
//	err = plug.Execute()
//	assert.Nil(t, err)
//}

func TestSinkShutdown(t *testing.T) {
	plug, err := New(plugin.SINK)
	assert.Nil(t, err)
	err = plug.Shutdown()
	assert.Nil(t, err)
}

func TestSourceInitialize(t *testing.T) {
	plug, err := New(plugin.SOURCE)
	assert.Nil(t, err)
	err = plug.Initialize()
	assert.Nil(t, err)
}

//func TestSourceExecute(t *testing.T) {
//	plug, err := New(plugin.SOURCE)
//	assert.Nil(t, err)
//	err = plug.Execute()
//	assert.Nil(t, err)
//}

func TestSourceShutdown(t *testing.T) {
	plug, err := New(plugin.SOURCE)
	assert.Nil(t, err)
	err = plug.Shutdown()
	assert.Nil(t, err)
}

func TestUnimplemented(t *testing.T) {
	_, err := New(plugin.MESSAGE)
	assert.EqualError(t, plugin.ErrUnimplemented, err.Error())
}
