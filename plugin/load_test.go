// +build !windows

package plugin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	_, err := New("noop", SINK)
	assert.NoError(t, err)
	// assert.NotNil(t, p)
}
