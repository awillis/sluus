package plugin

import (
	"container/ring"
	uuid2 "github.com/google/uuid"
	"zystus/core"
)

func NewBatch() core.Batch {

	uuid := uuid2.New()

	batch := core.Batch{
		uuid.String(),
		ring.New(5),
	}

	return batch
}

func (b core.Batch) AddEvent(event core.Event) {
	b.Ring.Value = event
	b.Ring.Next()
}
