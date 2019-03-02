package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	b *Batch
)

func TestNewBatch(t *testing.T) {
	msg, err := New(&testStruct)
	assert.NoError(t, err)
	b = NewBatch(1)
	assert.Nil(t, b.Add(msg))
}

func TestNewBatchFull(t *testing.T) {
	m1, err := New(&testStruct)
	assert.NoError(t, err)
	m2, err := NewFromString(testContent)
	assert.NoError(t, err)

	b := NewBatch(1)
	assert.Nil(t, b.Add(m1))
	err = b.Add(m2)
	assert.EqualError(t, err, ErrBatchFull.Error())
}

func TestBatch_Iter(t *testing.T) {
	m1, err := New(&testStruct)
	assert.NoError(t, err)
	m2, err := NewFromString(testContent)
	assert.NoError(t, err)

	b = NewBatch(2)
	assert.Nil(t, b.Add(m1))
	assert.Nil(t, b.Add(m2))
	for msg := range b.Iter() {
		_ = msg
		break
	}
	b.Cancel()
}
