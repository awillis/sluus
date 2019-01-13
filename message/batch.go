package message

import (
	"context"
	"errors"
	"sort"

	"github.com/google/uuid"
)

var ErrBatchFull = errors.New("batch is at capacity")

type Batch struct {
	sort.Interface
	id         string
	msgs       []Message
	CancelIter context.CancelFunc
}

func NewBatch(size int) Batch {

	batch := Batch{
		id:   uuid.New().String(),
		msgs: make([]Message, 0, size),
	}

	return batch
}

func (b *Batch) Id() string {
	return b.id
}

func (b *Batch) Add(m Message) error {
	if len(b.msgs) == cap(b.msgs) {
		return ErrBatchFull
	}
	b.msgs = append(b.msgs, m)
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

// sort.Processor methods

func (b Batch) Len() int {
	return len(b.msgs)
}

func (b Batch) Less(i, j int) bool {
	return b.msgs[i].Priority < b.msgs[j].Priority ||
		b.msgs[i].GetReceived().GetSeconds() < b.msgs[j].GetReceived().GetSeconds()
}

func (b Batch) Swap(i, j int) {
	b.msgs[i], b.msgs[j] = b.msgs[j], b.msgs[i]
}
