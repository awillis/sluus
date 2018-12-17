package message

import (
	"container/ring"
	"github.com/google/uuid"
)

type Message interface {
	PipelineID() uuid.UUID
	Payload() interface{}
}

type Batch struct {
	ID   string
	Ring *ring.Ring
}

func (b Batch) AddMessage(message Message) {
	b.Ring.Value = message
	b.Ring.Next()
}
