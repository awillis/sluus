package message

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	msg := New()
	assert.IsType(t, &Message{}, msg)
}

func TestMessage_GetContent(t *testing.T) {
	msg := New()
	assert.Nil(t, msg.GetContent())
	content := json.RawMessage("{\"content\":\"testytester\"}")
	msg, err := WithContent(content)
	assert.Nil(t, err)
	assert.IsType(t, &Message{}, msg)
	assert.Equal(t, "testytester", msg.GetContent().GetStringValue())
}

func TestMessage_GetReceived(t *testing.T) {
	msg := New()
	assert.Nil(t, msg.GetReceived())
	content := json.RawMessage("{\"received\":\"2019-01-01T01:01:01Z\"}")
	msg, err := WithContent(content)
	assert.Nil(t, err)
	assert.IsType(t, &Message{}, msg)
	assert.Equal(t, int64(1546304461), msg.GetReceived().GetSeconds())
}

func TestMessage_GetDirection(t *testing.T) {
	msg := New()
	assert.Equal(t, Message_IN, msg.GetDirection())
	content := json.RawMessage("{\"direction\": 1}")
	msg, err := WithContent(content)
	assert.Nil(t, err)
	assert.IsType(t, &Message{}, msg)
	assert.Equal(t, Message_OUT, msg.GetDirection())
}
