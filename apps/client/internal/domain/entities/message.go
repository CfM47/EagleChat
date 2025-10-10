package entities

import "time"

type Message struct {
	Sender      User
	Content     string
	CreatedTime time.Time
	ChatIndex   uint
}
