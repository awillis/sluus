package message

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testContent = "{\"foo\":\"bar\",\"a\":1,\"b\":true,\"nest\":{\"x\":\"y\"},\"list\":[false, 5, \"baz\"]}"
)

func TestNew(t *testing.T) {
	msg, err := New(testContent)
	assert.NoError(t, err)
	assert.Equal(t, true, BoolValue(FieldValue(msg.Body(), "b")))
	StringValue(FieldValue(msg.Body(), "b"))
	assert.Equal(t, "y", StringValue(FieldValue(StructValue(msg.FieldValue("nest")), "x")))
	_, err = New(testContent + "}")
	assert.Error(t, err)
}

func TestWithContentByte(t *testing.T) {
	_, err := NewFromBytes(json.RawMessage(testContent))
	assert.NoError(t, err)
	_, err = NewFromBytes([]byte(":"))
	assert.Error(t, err)
}

func TestMessage_GetContent(t *testing.T) {
	_, err := New(string(testContent))
	assert.NoError(t, err)
}

func TestMessage_Redirect(t *testing.T) {
	msg, err := New(testContent)
	assert.NoError(t, err)
	msg.Redirect(Message_REJECT)
	assert.Equal(t, Message_REJECT, msg.Direction)
}

func TestFromString(t *testing.T) {
	msg, err := New(testContent)
	assert.NoError(t, err)
	content, err := msg.ToString()
	assert.NoError(t, err)
	m, err := FromString(content)
	assert.NoError(t, err)
	assert.Equal(t, msg.FieldValue("foo"), m.FieldValue("foo"))
}

func TestMessage_Fields(t *testing.T) {
	msg, err := New(testContent)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "foo", "list", "nest"}, msg.Fields())
}

func TestMessage_ToString(t *testing.T) {
	msg, err := New(string(testContent))
	assert.NoError(t, err)
	_, err = msg.ToString()
	assert.NoError(t, err)
}

func TestValues(t *testing.T) {
	msg, err := New(string(testContent))
	assert.NoError(t, err)
	list := msg.FieldValue("list")
	for _, val := range Values(list) {
		_ = val
	}
}
