package core

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/struct"
)

const (
	URGENT uint32 = iota
	HIGH
	NORMAL
	LOW
)

func NewMessage(priority uint32) *Message {

	m := new(Message)
	m.Priority = priority
	return m
}

func (m *Message) Field(name string) *structpb.Value {
	return m.Body.Fields[name]
}

func (m *Message) SetField(name string, value *structpb.Value) *Message {
	m.Body.Fields[name] = value
	return m
}

func (m *Message) MarkReceived() {
	m.Received = ptypes.TimestampNow()
}
