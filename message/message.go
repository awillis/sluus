package message

import (
	"container/ring"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type Message interface {
	PipelineID() uuid.UUID
	Content() string
}

type JSONMessage struct {
	json.RawMessage
}

func (jm *JSONMessage) PipelineID() uuid.UUID {
	return uuid.New()
}

func (jm *JSONMessage) Content() string {
	if err := jm.UnmarshalJSON(jm.RawMessage); err != nil {
		fmt.Printf("unable to unmarshall json: %s", err.Error())
	}
	return string(jm.RawMessage)
}

type Batch struct {
	ID   string
	Ring *ring.Ring
}

func (b Batch) AddMessage(message Message) {
	b.Ring.Value = message
	b.Ring.Next()
}
