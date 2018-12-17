package pipeline

import (
	"container/ring"
	"github.com/awillis/sluus/message"
	uuid2 "github.com/google/uuid"
)

func NewBatch() message.Batch {

	uuid := uuid2.New()

	batch := message.Batch{
		ID:   uuid.String(),
		Ring: ring.New(5),
	}

	return batch
}
