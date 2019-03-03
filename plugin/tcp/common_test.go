package tcp

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

func TestUnimplemented(t *testing.T) {
	_, err := New(42)
	assert.EqualError(t, plugin.ErrUnimplemented, err.Error())
}

func TestPort(t *testing.T) {
	proc, _ := New(plugin.SOURCE)
	if opts, ok := proc.Options().(*options); ok {
		plugin.Validate(opts, opts.defaultPort())
		assert.Equal(t, uint64(3030), opts.Port)
	}
}

func TestBatchSize(t *testing.T) {
	proc, _ := New(plugin.SOURCE)
	if opts, ok := proc.Options().(*options); ok {
		plugin.Validate(opts, opts.defaultBatchSize())
		assert.Equal(t, uint64(64), opts.BatchSize)
	}
}

func TestBufferSize(t *testing.T) {
	proc, _ := New(plugin.SOURCE)
	if opts, ok := proc.Options().(*options); ok {
		plugin.Validate(opts, opts.defaultBufferSize())
		assert.Equal(t, uint64(16384), opts.BufferSize)
	}
}

func TestSockBufferSize(t *testing.T) {
	proc, _ := New(plugin.SOURCE)
	if opts, ok := proc.Options().(*options); ok {
		plugin.Validate(opts, opts.defaultSockBufferSize())
		assert.Equal(t, uint64(65536), opts.SockBufferSize)
	}
}

func TestPollInterval(t *testing.T) {
	proc, _ := New(plugin.SOURCE)
	if opts, ok := proc.Options().(*options); ok {
		plugin.Validate(opts, opts.defaultPollInterval())
		assert.Equal(t, uint64(200), opts.PollInterval)
	}
}
