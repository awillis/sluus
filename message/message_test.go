package message

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testContent = "{\"foo\":\"bar\",\"a\":1,\"b\":true,\"nest\":{\"x\":\"y\"}}"
)

func TestNew(t *testing.T) {
	assert.IsType(t, &Message{}, New())
}

func TestWithContent(t *testing.T) {
	msg, err := WithContent(testContent)
	assert.NoError(t, err)
	assert.Equal(t, true, msg.FieldValueByName("b").GetBoolValue())
	assert.Equal(t, "y", msg.FieldValueByName("nest").
		GetStructValue().GetFields()["x"].GetStringValue())
	_, err = WithContent(testContent + "}")
	assert.Error(t, err)
}

func TestWithContentByte(t *testing.T) {
	_, err := WithContentByte(json.RawMessage(testContent))
	assert.NoError(t, err)
	_, err = WithContentByte([]byte(":"))
	assert.Error(t, err)
}

func TestMessage_GetContent(t *testing.T) {
	_, err := WithContent(string(testContent))
	assert.NoError(t, err)
}

func TestMessage_Redirect(t *testing.T) {
	msg := New()
	msg.Redirect(Message_REJECT)
	assert.Equal(t, Message_REJECT, msg.Direction)
}

func TestFromString(t *testing.T) {
	msg, err := WithContent(testContent)
	assert.NoError(t, err)
	content, err := msg.ToString()
	assert.NoError(t, err)
	m, err := FromString(content)
	assert.NoError(t, err)
	assert.Equal(t, msg.FieldValueByName("foo"), m.FieldValueByName("foo"))
}

func TestMessage_Fields(t *testing.T) {
	msg, err := WithContent(testContent)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "foo", "nest"}, msg.Fields())
}

func TestMessage_ToString(t *testing.T) {
	msg, err := WithContent(string(testContent))
	assert.NoError(t, err)
	_, err = msg.ToString()
	assert.NoError(t, err)
}
