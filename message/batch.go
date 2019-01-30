package message

import (
	"context"
	"errors"
)

var ErrBatchFull = errors.New("batch is at capacity")

type Batch struct {
	msgs       []*Message
	CancelIter context.CancelFunc
}

func NewBatch(size int) *Batch {
	return &Batch{
		msgs: make([]*Message, 0, size),
	}
}

func (b *Batch) Add(m *Message) (err error) {
	if len(b.msgs) == cap(b.msgs) {
		return ErrBatchFull
	}
	b.msgs = append(b.msgs, m)
	return err
}

func (b *Batch) Count() uint64 {
	return uint64(len(b.msgs))
}

func (b *Batch) Iter() <-chan *Message {
	iter := make(chan *Message)
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

func (b *Batch) Cancel() {
	b.CancelIter()
}
