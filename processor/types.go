package processor

import (
	"container/ring"
	"github.com/golang-collections/go-datastructures/queue"
)

type Event interface {
	Body() string
	Compare(other queue.Item) int
}

type Batch struct {
	ID string
	Ring *ring.Ring
}

type Processor interface {
	Category() string
	Input() chan Event
	Output() chan Event
	Queue() queue.PriorityQueue
}

type Source interface {
	Processor
	Produce() error
}

type Conduit interface {
	Processor
	Process(batch Batch) error
}

type Sink interface {
	Processor
	Consume() error
}

