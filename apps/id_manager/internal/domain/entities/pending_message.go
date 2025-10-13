package entities

// Represents a the info of a pending message in EagleChat system
type PendingMessage struct {
	MessageId string   `json:"message_id"` // the message unique id
	TargetId  string   `json:"target_id"`  // the recipient unique id
	CachersId []string `json:"cachers_id"` // the ids of users caching the message
}

// NewPendingMessage creates a new PendingMessage object.
func NewPendingMessage(message_id string, target_id string, cachers_id []string) *PendingMessage {
	return &PendingMessage{
		MessageId: message_id,
		TargetId:  target_id,
		CachersId: cachers_id,
	}
}
