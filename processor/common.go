package processor

import (
	"context"
	"github.com/awillis/sluus/message"
	"github.com/ef-ds/deque"
	"github.com/pkg/errors"
	"runtime"
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
	defer g.Unlock()
	if val, ok := g.deque.PopFront(); ok {
		return val.(*message.Batch)
	}
	return
}

func (g *gate) Put(batch *message.Batch) {
	g.Lock()
	defer g.Unlock()
	g.deque.PushBack(batch)
}

func (g *gate) Len() int {
	return g.deque.Len()
}

func (g *gate) Poll(timeout time.Duration) (batch *message.Batch) {

	start := time.Now()

	for {
		batch = g.Get()

		if batch != nil || time.Since(start) >= timeout {
			return
		}

		runtime.Gosched()
	}
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
