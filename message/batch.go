package message

import (
	"context"
	"errors"
)

var (
	routes       = []Route{Route_PASS, Route_ACCEPT, Route_REJECT}
	ErrBatchFull = errors.New("batch is at capacity")
)

type Batch struct {
	msgs   []*Message
	index  map[Route][]int
	cancel context.CancelFunc
}

func NewBatch(size uint64) (batch *Batch) {
	batch = new(Batch)
	batch.msgs = make([]*Message, 0, size)
	batch.index = make(map[Route][]int)

	for _, rte := range routes {
		batch.index[rte] = make([]int, 0, size)
	}

	return
}

func (b *Batch) add(m *Message) {
	b.index[m.Direction] = append(b.index[m.Direction], len(b.msgs))
	b.msgs = append(b.msgs, m)
}

func (b *Batch) Add(m *Message) {
	if len(b.msgs) == cap(b.msgs) {
		panic(ErrBatchFull)
	}
	b.add(m)
}

func (b *Batch) AddE(m *Message) (err error) {
	if len(b.msgs) == cap(b.msgs) {
		return ErrBatchFull
	}
	b.add(m)
	return
}

func (b *Batch) Clear() {
	b.msgs = b.msgs[:0]
	b.clearIndex()
}

func (b *Batch) clearIndex() {
	for _, rte := range routes {
		b.index[rte] = b.index[rte][:0]
	}
}

func (b *Batch) reIndex() {
	for i, m := range b.msgs {
		b.index[m.Direction] = append(b.index[m.Direction], i)
	}
}

func (b *Batch) route(route Route) (batch *Batch) {
	batch = NewBatch(uint64(len(b.index[route])))
	for _, i := range b.index[route] {
		batch.Add(b.msgs[i])
	}
	return
}

func (b *Batch) Pass() *Batch {
	return b.route(Route_PASS)
}

func (b *Batch) Accept() *Batch {
	return b.route(Route_ACCEPT)
}

func (b *Batch) Reject() *Batch {
	return b.route(Route_REJECT)
}

func (b *Batch) Count() uint64 {
	return uint64(len(b.msgs))
}

func (b *Batch) PassCount() uint64 {
	return uint64(len(b.index[Route_PASS]))
}

func (b *Batch) AcceptCount() uint64 {
	return uint64(len(b.index[Route_ACCEPT]))
}

func (b *Batch) RejectCount() uint64 {
	return uint64(len(b.index[Route_REJECT]))
}

func (b *Batch) Iter() <-chan *Message {
	iter := make(chan *Message)
	ctx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel

	go func(ctx context.Context) {
		defer close(iter)
	cancel:
		for i := 0; i < len(b.msgs); i++ {
			select {
			case <-ctx.Done():
				b.msgs = b.msgs[i:]
				b.clearIndex()
				b.reIndex()
				break cancel
			case iter <- b.msgs[i]:
				continue
			}
		}
	}(ctx)
	return iter
}

func (b *Batch) Cancel() {
	b.cancel()
}
