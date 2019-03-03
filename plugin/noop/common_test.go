package noop

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

func TestNewConduit(t *testing.T) {
	proc, err := New(plugin.CONDUIT)
	assert.NoError(t, err)
	assert.Equal(t, plugin.CONDUIT, proc.Type())
	assert.Implements(t, (*plugin.Processor)(nil), proc)
}

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

func TestMessagePerBatch(t *testing.T) {
	proc, _ := New(plugin.SOURCE)
	if opts, ok := proc.Options().(*options); ok {
		plugin.Validate(opts, opts.defaultMessagePerBatch())
		assert.Equal(t, uint64(5), opts.MessagePerBatch)
	}
}

func TestBatchInterval(t *testing.T) {
	proc, _ := New(plugin.SOURCE)
	if opts, ok := proc.Options().(*options); ok {
		plugin.Validate(opts, opts.defaultBatchInterval())
		assert.Equal(t, uint64(5), opts.BatchInterval)
	}
}

func TestRejectPercentage(t *testing.T) {
	proc, _ := New(plugin.SOURCE)
	if opts, ok := proc.Options().(*options); ok {
		opts.MessagePerBatch = 20
		opts.RejectPercentage = 30
		plugin.Validate(opts, opts.defaultRejectPercentage())
		assert.Equal(t, uint64(20), opts.RejectPercentage)
		opts.RejectPercentage = 10
		assert.Equal(t, uint64(10), opts.RejectPercentage)
	}
}

func TestAcceptPercentage(t *testing.T) {
	proc, _ := New(plugin.SOURCE)
	if opts, ok := proc.Options().(*options); ok {
		opts.MessagePerBatch = 20
		opts.AcceptPercentage = 30
		plugin.Validate(opts, opts.defaultAcceptPercentage())
		assert.Equal(t, uint64(20), opts.AcceptPercentage)
		opts.AcceptPercentage = 10
		assert.Equal(t, uint64(10), opts.AcceptPercentage)
	}
}
