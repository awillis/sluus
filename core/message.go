package core

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Message struct {
	ID      bson.ObjectId `bson:"id" json:"id"`
	Meta    MessageMeta   `bson:"meta,omitempty" json:"meta,omitempty"`
	Content bson.M        `bson:"content" json:"content"`
}

type MessageMeta struct {
	Timestamp string `bson:"timestamp" json:"timestamp"`
	Priority  int    `bson:"priority" json:"priority"`
}

func NewMessage(priority int) *Message {
	msg := new(Message)
	msg.ID = bson.NewObjectId()
	msg.ID.Machine()
	msg.Meta = MessageMeta{
		msg.ID.Time().Format(time.RFC3339),
		priority,
	}
	return msg
}

func (msg *Message) Field(name string) interface{} {
	return msg.Content[name]
}

func (msg *Message) SetField(name string, value interface{}) *Message {
	msg.Content[name] = value
	return msg
}

func (msg *Message) Contents() bson.M {
	return msg.Content
}

func (msg *Message) SetContent(content map[string]interface{}) {
	msg.Content = content
}
