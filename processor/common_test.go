package processor

import (
	"github.com/awillis/sluus/message"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGate_PutGet(t *testing.T) {
	g := newGate()
	batch := message.NewBatch(3)
	assert.NotNil(t, batch)
	g.Put(batch)
	b := g.Get()
	assert.NotNil(t, b)
}
