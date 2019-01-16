package message

import (
	"bytes"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
)

var unmarshaler *jsonpb.Unmarshaler

func init() {
	unmarshaler = new(jsonpb.Unmarshaler)
	unmarshaler.AllowUnknownFields = true
}

func New(priority Message_Priority, content []byte) (msg *Message, err error) {
	err = unmarshaler.Unmarshal(bytes.NewBuffer(content), msg)
	msg.Priority = priority
	return
}

func (m *Message) MarkReceived() {
	m.Received = ptypes.TimestampNow()
}
