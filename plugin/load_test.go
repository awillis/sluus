// +build !windows

package plugin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	p, err := New("noop", SINK)
	assert.Nil(t, err)
	assert.NotNil(t, p)
}
