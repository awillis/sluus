package message

import (
	"bytes"
	"encoding/json"
	"sort"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/pkg/errors"
)

var (
	marshaler      *jsonpb.Marshaler
	unmarshaler    *jsonpb.Unmarshaler
	ErrInvalidJson = errors.New("invalid json")
)

func init() {
	marshaler = new(jsonpb.Marshaler)
	marshaler.EmitDefaults = true
	unmarshaler = new(jsonpb.Unmarshaler)
	unmarshaler.AllowUnknownFields = true
}

func New(content interface{}) (msg *Message, err error) {

	msi := make(map[string]interface{})
	msi["content"] = content

	js, err := json.Marshal(msi)

	if err != nil {
		return
	}

	// protobuf unmarshal
	msg = new(Message)
	err = unmarshaler.Unmarshal(bytes.NewReader(js), msg)
	msg.Received = ptypes.TimestampNow()
	return msg, err
}

func NewFromString(content string) (msg *Message, err error) {

	// the runtime uses a 32 byte buffer for string concatenation
	// a string builder should result in a single allocation
	var sb strings.Builder
	head, foot := `{"content":`, `}`
	sb.Grow(len(head) + len(content) + len(foot) + 4)

	sb.WriteString(head)
	if !json.Valid([]byte(content)) {
		// quote content if not valid json
		sb.WriteString(`"`)
		sb.WriteString(content)
		sb.WriteString(`"`)
	} else {

		sb.WriteString(content)
	}
	sb.WriteString(foot)

	// json unmarshal / marshal to wrap content as json
	msi := make(map[string]interface{})
	if err = json.Unmarshal(json.RawMessage(sb.String()), &msi); err != nil {
		return
	}

	js, err := json.Marshal(msi)

	if err != nil {
		return
	}

	// protobuf unmarshal
	msg = new(Message)
	err = unmarshaler.Unmarshal(bytes.NewReader(js), msg)
	msg.Received = ptypes.TimestampNow()
	return msg, err
}

func NewFromBytes(content []byte) (msg *Message, err error) {
	if !json.Valid(content) {
		return msg, ErrInvalidJson
	}
	return NewFromString(string(content))
}

func FromString(payload string) (msg *Message, err error) {
	msg = new(Message)
	err = unmarshaler.Unmarshal(strings.NewReader(payload), msg)
	return
}

func FromBytes(payload []byte) (msg *Message, err error) {
	msg = new(Message)
	err = unmarshaler.Unmarshal(bytes.NewReader(payload), msg)
	return
}

func (m *Message) ToString() (content string, err error) {
	return marshaler.MarshalToString(m)
}

func (m *Message) ToBytes() (content []byte, err error) {
	buf := bytes.NewBuffer(content)
	err = marshaler.Marshal(buf, m)
	return buf.Bytes(), err
}

func (m *Message) Redirect(direction Route) {
	m.Direction = direction
}

func (m *Message) Body() *structpb.Struct {
	return m.GetContent().GetStructValue()
}

func (m *Message) FieldValue(name string) *structpb.Value {
	return m.Body().Fields[name]
}

func (m *Message) Fields() (fields []string) {
	for key := range m.Body().Fields {
		fields = append(fields, key)
	}
	sort.Strings(fields)
	return
}

func StructValue(value *structpb.Value) *structpb.Struct {
	return value.GetStructValue()
}

func FieldValue(value *structpb.Struct, name string) *structpb.Value {
	return value.Fields[name]
}

func Values(value *structpb.Value) []*structpb.Value {
	return value.GetListValue().Values
}

func StringValue(value *structpb.Value) string {
	return value.GetStringValue()
}

func BoolValue(value *structpb.Value) bool {
	return value.GetBoolValue()
}

func NumberValue(value *structpb.Value) float64 {
	return value.GetNumberValue()
}
