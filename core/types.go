package core

import (
	"container/ring"
	"github.com/golang-collections/go-datastructures/queue"
	"github.com/google/uuid"
)

type Message interface {
	PipelineID() uuid.UUID
	Payload() interface{}
	Compare(other queue.Item) int
}

type Batch struct {
	ID   string
	Ring *ring.Ring
}

func (b Batch) AddMessage(message Message) {
	b.Ring.Value = message
	b.Ring.Next()
}

type Processor interface {
	Category() string
	Input() chan Message
	Output() chan Message
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
