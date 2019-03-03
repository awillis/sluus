// +build !windows

package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/awillis/sluus/plugin"
)

func TestNewSource(t *testing.T) {
	plug, err := New(plugin.SOURCE)
	assert.Nil(t, err)
	assert.Equal(t, plug.Type(), plugin.SOURCE)
	assert.Implements(t, (*plugin.Producer)(nil), plug)
}

func TestSource(t *testing.T) {
	plug, err := New(plugin.SOURCE)
	assert.Nil(t, err)
	assert.Implements(t, (*plugin.Producer)(nil), plug)
	if source, ok := plug.(plugin.Producer); ok {
		source.Start(context.Background())
		assert.Nil(t, source.Shutdown())
	}
}

func TestNewConduit(t *testing.T) {
	plug, err := New(plugin.CONDUIT)
	assert.Nil(t, err)
	assert.Equal(t, plug.Type(), plugin.CONDUIT)
	assert.Implements(t, (*plugin.Processor)(nil), plug)
}

func TestConduit(t *testing.T) {
	plug, err := New(plugin.CONDUIT)
	assert.Nil(t, err)
	assert.Implements(t, (*plugin.Processor)(nil), plug)
	if conduit, ok := plug.(plugin.Processor); ok {
		conduit.Start(context.Background())
		assert.Nil(t, conduit.Shutdown())
	}
}

func TestNewSink(t *testing.T) {
	plug, err := New(plugin.SINK)
	assert.Nil(t, err)
	assert.Equal(t, plug.Type(), plugin.SINK)
	assert.Implements(t, (*plugin.Consumer)(nil), plug)
}

func TestSink(t *testing.T) {
	plug, err := New(plugin.SINK)
	assert.Nil(t, err)
	assert.Implements(t, (*plugin.Consumer)(nil), plug)
	if sink, ok := plug.(plugin.Consumer); ok {
		sink.Start(context.Background())
		assert.Nil(t, sink.Shutdown())
	}
}

func TestUnimplemented(t *testing.T) {
	_, err := New(42)
	assert.EqualError(t, plugin.ErrUnimplemented, err.Error())
}
