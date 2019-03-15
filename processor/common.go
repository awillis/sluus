package processor

import (
	"context"
	"github.com/awillis/sluus/message"
	"github.com/ef-ds/deque"
	"github.com/pkg/errors"
	"sync"
	"time"
)

var ErrTimeout = errors.New("operation timed out")

type (
	gate struct {
		sync.Mutex
		deque *deque.Deque
	}

	workGroup struct {
		sync.WaitGroup
		cancel context.CancelFunc
	}

	head struct {
		sync.RWMutex
		m map[uint64][]byte
	}
)

func newGate() *gate {
	return &gate{
		deque: deque.New(),
	}
}

func (g *gate) Get() (batch *message.Batch) {
	g.Lock()
	val, ok := g.deque.PopFront()
	g.Unlock()

	if ok {
		return val.(*message.Batch)
	}
	return
}

func (g *gate) Put(batch *message.Batch) {
	g.Lock()
	g.deque.PushBack(batch)
	g.Unlock()
}

func (g *gate) Len() int {
	return g.deque.Len()
}

func (g *gate) Poll(ctx context.Context, timeout time.Duration) (batch *message.Batch) {

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

loop:
	select {
	case <-ctx.Done():
		break loop
	default:
		batch = g.Get()

		if batch != nil {
			break loop
		}
		goto loop
	}
	return
}

func (w *workGroup) Shutdown() {
	w.cancel()
	w.Wait()
}

func (h *head) Get(direction uint64) []byte {
	h.RLock()
	defer h.RUnlock()
	return h.m[direction]
}

func (h *head) Set(direction uint64, key []byte) {
	h.Lock()
	defer h.Unlock()
	h.m[direction] = key
}

func (h *head) Reset(direction uint64) {
	h.Lock()
	defer h.Unlock()
	h.m[direction] = nil
	h.m[direction] = make([]byte, 0, 32)
}
