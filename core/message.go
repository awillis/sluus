package core

import (
	"github.com/golang/protobuf/ptypes"
)

func NewMessage(priority Message_Priority) *Message {

	m := new(Message)
	m.Priority = priority
	return m
}

func (m *Message) MarkReceived() {
	m.Received = ptypes.TimestampNow()
}
