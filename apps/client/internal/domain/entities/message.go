package entities

import "time"

type Message struct {
	Sender      User
	Content     string
	CreatedTime time.Time
}

func NewMessage(sender User, content string) Message {
	return Message{
		Sender:      sender,
		Content:     content,
		CreatedTime: time.Now().UTC().Round(0),
	}
}
