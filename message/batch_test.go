package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	b  *Batch
	m1 *Message
	m2 *Message
)

func init() {
	m1, m2 = New(), New()
}
func TestNewBatch(t *testing.T) {
	b = NewBatch(1)
	assert.Nil(t, b.Add(m1))
}

func TestNewBatchFull(t *testing.T) {
	b := NewBatch(1)
	assert.Nil(t, b.Add(m1))
	err := b.Add(m2)
	assert.EqualError(t, err, ErrBatchFull.Error())
}

func TestBatch_Iter(t *testing.T) {
	b = NewBatch(2)
	assert.Nil(t, b.Add(m1))
	assert.Nil(t, b.Add(m2))
	for msg := range b.Iter() {
		_ = msg
		break
	}
	b.Cancel()
}
