package entities

import "eaglechat/apps/client/internal/domain/entities"

type PendingMessage struct {
	Target  MessageTarget
	Content []byte
}

func NewPendingMessage(target MessageTarget, content []byte) PendingMessage {
	return PendingMessage{
		Target:  target,
		Content: content,
	}
}

type MessageTarget struct {
	MessageID string
	TargetID  entities.UserID
}

func NewMessageTarget(ID string, target entities.UserID) MessageTarget {
	return MessageTarget{
		MessageID: ID,
		TargetID:  target,
	}
}
