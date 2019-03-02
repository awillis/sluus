package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testStruct  = msgTest{Foo: "abc", Bar: true, Baz: 42}
	testContent = "{\"foo\":\"bar\",\"a\":1,\"b\":true,\"nest\":{\"x\":\"y\"},\"list\":[false, 5, \"baz\"]}"
)

type msgTest struct {
	Foo string `json:"foo"`
	Bar bool   `json:"bar"`
	Baz uint64 `json:"baz"`
}

func TestNew(t *testing.T) {
	msg, err := New(&testStruct)
	assert.NoError(t, err)
	assert.Equal(t, float64(42), NumberValue(msg.FieldValue("baz")))
}

func TestNewFromString(t *testing.T) {
	msg, err := NewFromString(testContent)
	assert.NoError(t, err)
	assert.Equal(t, true, BoolValue(FieldValue(msg.Body(), "b")))
	StringValue(FieldValue(msg.Body(), "b"))
	assert.Equal(t, "y", StringValue(FieldValue(StructValue(msg.FieldValue("nest")), "x")))
	_, err = NewFromString(testContent + "}")
	assert.Error(t, err)
}

func TestNewFromBytes(t *testing.T) {
	_, err := NewFromBytes([]byte(testContent))
	assert.NoError(t, err)
	_, err = NewFromBytes([]byte(":"))
	assert.Error(t, err)
}

func TestMessage_GetContent(t *testing.T) {
	_, err := NewFromString(testContent)
	assert.NoError(t, err)
}

func TestMessage_Redirect(t *testing.T) {
	msg, err := New(&testStruct)
	assert.NoError(t, err)
	msg.Redirect(Message_REJECT)
	assert.Equal(t, Message_REJECT, msg.Direction)
}

func TestFromString(t *testing.T) {
	msg, err := NewFromString(testContent)
	assert.NoError(t, err)
	content, err := msg.ToString()
	assert.NoError(t, err)
	m, err := FromString(content)
	assert.NoError(t, err)
	assert.Equal(t, StringValue(msg.FieldValue("foo")), StringValue(m.FieldValue("foo")))
}

func TestToFromBytes(t *testing.T) {
	msg, err := New(&testStruct)
	assert.NoError(t, err)
	content, err := msg.ToBytes()
	assert.NoError(t, err)
	m, err := FromBytes(content)
	assert.NoError(t, err)
	assert.Equal(t, StringValue(msg.FieldValue("foo")), StringValue(m.FieldValue("foo")))
}

func TestMessage_Fields(t *testing.T) {
	msg, err := NewFromString(testContent)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "foo", "list", "nest"}, msg.Fields())
}

func TestMessage_ToString(t *testing.T) {
	msg, err := New(&testStruct)
	assert.NoError(t, err)
	_, err = msg.ToString()
	assert.NoError(t, err)
}

func TestValues(t *testing.T) {
	msg, err := NewFromString(testContent)
	assert.NoError(t, err)
	list := msg.FieldValue("list")
	for _, val := range Values(list) {
		_ = val
	}
}
