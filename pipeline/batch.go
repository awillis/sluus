package pipeline

import (
	"container/ring"
	"github.com/awillis/sluus/core"
	uuid2 "github.com/google/uuid"
)

func NewBatch() core.Batch {

	uuid := uuid2.New()

	batch := core.Batch{
		uuid.String(),
		ring.New(5),
	}

	return batch
}
