package main

import (
	"github.com/awillis/sluus/plugin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSink(t *testing.T) {
	plug, err := New(plugin.SINK)
	assert.Nil(t, err)
	assert.Equal(t, plug.Type(), plugin.SINK)
}

func TestNewConduit(t *testing.T) {
	plug, err := New(plugin.CONDUIT)
	assert.Nil(t, err)
	assert.Equal(t, plug.Type(), plugin.CONDUIT)
}

func TestSinkInitialize(t *testing.T) {
	plug, err := New(plugin.SINK)
	assert.Nil(t, err)
	err = plug.Initialize()
	assert.Nil(t, err)
}

func TestSinkExecute(t *testing.T) {
	plug, err := New(plugin.SINK)
	assert.Nil(t, err)
	err = plug.Execute()
	assert.Nil(t, err)
}

func TestSinkShutdown(t *testing.T) {
	plug, err := New(plugin.SINK)
	assert.Nil(t, err)
	err = plug.Shutdown()
	assert.Nil(t, err)
}

func TestConduitInitialize(t *testing.T) {
	plug, err := New(plugin.CONDUIT)
	assert.Nil(t, err)
	err = plug.Initialize()
	assert.Nil(t, err)
}

func TestConduitExecute(t *testing.T) {
	plug, err := New(plugin.CONDUIT)
	assert.Nil(t, err)
	err = plug.Execute()
	assert.Nil(t, err)
}

func TestConduitShutdown(t *testing.T) {
	plug, err := New(plugin.CONDUIT)
	assert.Nil(t, err)
	err = plug.Shutdown()
	assert.Nil(t, err)
}

func TestUnimplemented(t *testing.T) {
	_, err := New(plugin.MESSAGE)
	assert.EqualError(t, plugin.ErrUnimplemented, err.Error())
}
