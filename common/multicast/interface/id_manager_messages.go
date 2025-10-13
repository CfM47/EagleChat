package multicast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

type IDManagerMessage struct {
	Type      BroadcastType `json:"type"`
	ID        string        `json:"id"`
	IP        string        `json:"ip"`
	PublicKey []byte        `json:"public_key"`
	Port      string        `json:"port"`
}

func (m *BroadcastMessage) AsIDManagerMessage() (IDManagerMessage, error) {
	if m.Type != AnnounceIDManager && m.Type != UpdateIDManager {
		return IDManagerMessage{}, ErrWrongMessageType
	}

	var msg IDManagerMessage
	err := json.NewDecoder(bytes.NewBuffer(m.Content)).Decode(&msg)
	if err != nil {
		return IDManagerMessage{}, fmt.Errorf("unexpected decoding error: %v", err)
	}

	if msg.Type != AnnounceIDManager && msg.Type != UpdateIDManager {
		return IDManagerMessage{}, ErrWrongMessageType
	}

	return msg, nil
}

func BuildIDManagerAnnounce(ID string, IP string, Port string, publicKey []byte) (BroadcastMessage, error) {
	return BuildIDManagerMessage(AnnounceIDManager, ID, IP, Port, publicKey)
}

func BuildIDManagerUpdate(ID string, IP string, Port string, publicKey []byte) (BroadcastMessage, error) {
	return BuildIDManagerMessage(UpdateIDManager, ID, IP, Port, publicKey)
}

func BuildIDManagerMessage(ty BroadcastType, ID string, IP string, Port string, publicKey []byte) (BroadcastMessage, error) {
	innerMsg := IDManagerMessage{
		ID:        ID,
		IP:        IP,
		Port:      Port,
		Type:      ty,
		PublicKey: publicKey,
	}

	innerMsgBytes, err := json.Marshal(innerMsg)
	if err != nil {
		log.Printf("failed to marshal id manager broadcast message: %v", err)
		return BroadcastMessage{}, err
	}

	return BroadcastMessage{
		Type:    ty,
		Content: innerMsgBytes,
	}, nil
}
