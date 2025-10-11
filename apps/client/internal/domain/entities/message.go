package entities

import "time"

type Message struct {
	Sender      User
	Content     string
	CreatedTime time.Time
	ChatIndex   uint
}

func NewMessage(sender User, content string, chatIndex uint) Message {
	return Message{
		Sender:      sender,
		Content:     content,
		CreatedTime: time.Now().UTC().Round(0),
		ChatIndex:   chatIndex,
	}
}
