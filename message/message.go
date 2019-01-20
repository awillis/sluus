package message

import (
	"bytes"
	"encoding/json"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
)

var (
	marshaler   *jsonpb.Marshaler
	unmarshaler *jsonpb.Unmarshaler
)

func init() {
	marshaler = new(jsonpb.Marshaler)
	marshaler.EmitDefaults = true
	unmarshaler = new(jsonpb.Unmarshaler)
	unmarshaler.AllowUnknownFields = true
}

func New() (msg *Message) {
	return new(Message)
}

func WithContent(content json.RawMessage) (msg *Message, err error) {
	msg = new(Message)
	err = unmarshaler.Unmarshal(bytes.NewBuffer(content), msg)
	return
}

func (m *Message) MarkReceived() {
	m.Received = ptypes.TimestampNow()
}

func (m *Message) Redirect(direction Message_Direction) {
	m.Direction = direction
}
