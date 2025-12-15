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

func NewMessageTarget(messageID string, target entities.UserID) MessageTarget {
	return MessageTarget{
		MessageID: messageID,
		TargetID:  target,
	}
}
