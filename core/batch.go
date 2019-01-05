package core

import (
	"context"
	"errors"
	"sort"

	"github.com/google/uuid"
)

var errBatchFull = errors.New("batch is at capacity")

type Batch struct {
	sort.Interface
	ID         string
	msgs       []Message
	CancelIter context.CancelFunc
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
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	b.CancelIter = cancel

	go func(ctx context.Context) {
		defer close(iter)
		for i := 0; i < len(b.msgs); i++ {
			select {
			case <-ctx.Done():
				break
			case iter <- b.msgs[i]:
				continue
			}
		}
	}(ctx)
	return iter
}

// sort.Interface methods

func (b Batch) Len() int {
	return len(b.msgs)
}

func (b Batch) Less(i, j int) bool {
	return b.msgs[i].Priority < b.msgs[j].Priority || b.msgs[i].ID.Time().Before(b.msgs[j].ID.Time())
}

func (b Batch) Swap(i, j int) {
	b.msgs[i], b.msgs[j] = b.msgs[j], b.msgs[i]
}
