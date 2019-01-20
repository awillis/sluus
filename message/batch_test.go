package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBatch(t *testing.T) {
	b := NewBatch(1)
	m := New()
	err := b.Add(m)
	assert.Nil(t, err)
}

func TestNewBatchFull(t *testing.T) {
	b := NewBatch(1)
	m := New()
	err := b.Add(m)
	assert.Nil(t, err)
	n := New()
	err = b.Add(n)
	assert.EqualError(t, err, ErrBatchFull.Error())
}
