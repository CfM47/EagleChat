package implementation

import (
	"encoding/json"

	"eaglechat/common/multicast/interface"
)

// Broadcast sends a message to the multicast group.
func (n *udpMulticastNetwork) Broadcast(message multicast.BroadcastMessage) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = n.conn.WriteToUDP(bytes, n.addr)
	return err
}
