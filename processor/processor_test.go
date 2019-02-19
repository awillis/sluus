package processor

import (
	"github.com/awillis/sluus/plugin"
	"github.com/stretchr/testify/assert"
	"testing"
)

//func TestNew(t *testing.T) {
//	noop := New("noop", plugin.SINK)
//	assert.Equal(t, "noop", noop.Name )
//	assert.Equal(t, plugin.SINK, noop.Type())
//}

func TestNoPlugin(t *testing.T) {
	none := New("nonexistent", plugin.SINK)
	assert.Implements(t, (*Processor)(nil), none)
}
