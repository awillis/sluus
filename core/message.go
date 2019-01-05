package core

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	URGENT MessageLevel = iota
	HIGH
	NORMAL
	LOW
)

type MessageLevel uint8

type Message struct {
	ID        bson.ObjectId `bson:"id" json:"id"`
	Timestamp string        `bson:"timestamp" json:"timestamp"`
	Priority  MessageLevel  `bson:"priority" json:"priority"`
	Content   bson.M        `bson:"content" json:"content"`
}

func NewMessage(priority MessageLevel) *Message {

	msg := new(Message)
	msg.ID = bson.NewObjectId()
	msg.Timestamp = msg.ID.Time().Format(time.RFC3339)
	msg.Priority = priority
	return msg
}

func (msg *Message) Field(name string) interface{} {
	return msg.Content[name]
}

func (msg *Message) SetField(name string, value interface{}) *Message {
	msg.Content[name] = value
	return msg
}

func (msg *Message) SetContent(content map[string]interface{}) {
	msg.Content = content
}
