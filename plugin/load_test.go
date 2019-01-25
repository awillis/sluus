package plugin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProcessor(t *testing.T) {
	p, err := NewProcessor("noop", SINK)
	assert.Nil(t, err)
	assert.NotNil(t, p)
}
