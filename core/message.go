package core

import (
	"bytes"
	"encoding/json"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/struct"
)

func NewMessage(priority Message_Priority) *Message {

	m := new(Message)
	m.Priority = priority
	return m
}

func (m *Message) MarkReceived() {
	m.Received = ptypes.TimestampNow()
}

func (m *Message) SetBody(value interface{}) error {
	// Marshal value to json to ensure validity before setting it as protobuf value
	var err error
	jsval, err := json.Marshal(value)
	if err != nil {
		return err
	}

	jsbuf := bytes.NewReader(jsval)
	m.Body = &structpb.Value{}
	jsm := new(jsonpb.Unmarshaler)
	err = jsm.Unmarshal(jsbuf, m.Body)
	return err
}
