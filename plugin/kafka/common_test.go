package kafka

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/awillis/sluus/plugin"
)

func TestNewSource(t *testing.T) {
	proc, err := New(plugin.SOURCE)
	assert.NoError(t, err)
	assert.Equal(t, plugin.SOURCE, proc.Type())
	assert.Implements(t, (*plugin.Producer)(nil), proc)
}

//func TestNewConduit(t *testing.T) {
//	proc, err := New(plugin.CONDUIT)
//	assert.NoError(t, err)
//	assert.Equal(t, plugin.CONDUIT, proc.Type())
//	assert.Implements(t, (*plugin.Processor)(nil), proc)
//}

func TestNewSink(t *testing.T) {
	proc, err := New(plugin.SINK)
	assert.NoError(t, err)
	assert.Equal(t, plugin.SINK, proc.Type())
	assert.Implements(t, (*plugin.Consumer)(nil), proc)
}

func TestUnimplemented(t *testing.T) {
	_, err := New(42)
	assert.EqualError(t, plugin.ErrUnimplemented, err.Error())
}
