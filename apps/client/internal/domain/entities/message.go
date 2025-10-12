package entities

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID          string
	Sender      User
	Target      User
	Content     string
	CreatedTime time.Time
}

func NewMessage(sender User, target User, content string) Message {
	return Message{
		ID:          uuid.New().String(),
		Sender:      sender,
		Target:      target,
		Content:     content,
		CreatedTime: time.Now().UTC().Round(0),
	}
}
