package message

import (
	"os"
	"reflect"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

func NewMessage(priority Message_Priority) *Message {

	m := new(Message)
	m.Priority = priority
	return m
}

func (m *Message) MarkReceived() {
	m.Received = ptypes.TimestampNow()
}

func (m *Message) SetBody(value interface{}) (err error) {

	reflectType := reflect.TypeOf(value)
	reflectValue := reflect.ValueOf(&value)

	body, err := ptypes.MarshalAny(&any.Any{
		TypeUrl: reflectType.PkgPath() + string(os.PathSeparator) + reflectType.String(),
		Value:   reflectValue.Elem().Bytes(),
	})

	m.Body = body
	return err
}
