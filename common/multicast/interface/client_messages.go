package multicast

func (m *BroadcastMessage) AsClientRegisterMessage() error {
	if m.Type != RegisterClient {
		return ErrWrongMessageType
	}

	return nil
}

func NewClientRegisterMessage() BroadcastMessage {
	return BroadcastMessage{
		Type:    RegisterClient,
		Content: nil,
	}
}
