package core

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"sort"
)

var errBatchFull = errors.New("batch is at capacity")

type Batch struct {
	sort.Interface
	ID   string
	msgs []Message
}

func NewBatch(size int) Batch {

	batch := Batch{
		ID:   uuid.New().String(),
		msgs: make([]Message, 0, size),
	}

	return batch
}

func (b *Batch) Add(msg Message) error {
	if len(b.msgs) == cap(b.msgs) {
		return errBatchFull
	}
	b.msgs = append(b.msgs, msg)
	return nil
}

func (b Batch) Iter() <-chan Message {
	iter := make(chan Message)
	go func() {
		defer close(iter)
		for i := 0; i < len(b.msgs); i++ {
			iter <- b.msgs[i]
		}
	}()
	return iter
}

// sort.Interface methods

func (b Batch) Len() int {
	return len(b.msgs)
}

func (b Batch) Less(i, j int) bool {
	return b.msgs[i].Meta.Priority < b.msgs[j].Meta.Priority
}

func (b Batch) Swap(i, j int) {
	b.msgs[i], b.msgs[j] = b.msgs[j], b.msgs[i]
}
