package tcp

import (
	"github.com/awillis/sluus/plugin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPort(t *testing.T) {
	proc, _ := New(plugin.SOURCE)
	if opts, ok := proc.Options().(*options); ok {
		opts.port = 5
		plugin.Validate(opts, opts.defaultPort())
		assert.Equal(t, 5, opts.port)
	}
}
