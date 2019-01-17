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

func TestWithContent(t *testing.T) {
	content := json.RawMessage("{\"content\":\"testytester\"}")
	msg, err := WithContent(content)
	assert.Nil(t, err)
	assert.IsType(t, &Message{}, msg)
	assert.Equal(t, "testytester", msg.GetContent().GetStringValue())
}
