package processor

import (
	"container/ring"
	uuid2 "github.com/google/uuid"
)

func NewBatch() Batch {

	uuid := uuid2.New()

	batch := Batch{
		uuid.String(),
		ring.New(5),
	}

	return batch
}

func (b Batch) AddEvent(event Event)  {
	b.Ring.Value = event
	b.Ring.Next()
}