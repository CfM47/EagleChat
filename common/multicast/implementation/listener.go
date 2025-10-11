package implementation

import (
	"encoding/json"
	"log"

	multicast "eaglechat/common/multicast/interface"
)

const maxDatagramSize = 8192

// listenLoop runs in a background goroutine to read and process incoming multicast packets.
func (n *udpMulticastNetwork) listenLoop() {
	defer n.wg.Done()

	buf := make([]byte, maxDatagramSize)
	for {
		select {
		case <-n.done:
			return
		default:
			// Read from the UDP connection
			nBytes, _, err := n.conn.ReadFromUDP(buf)
			if err != nil {
				// Check if the error is due to the connection being closed, which is expected on shutdown.
				select {
				case <-n.done:
					// This error is expected because Close() was called.
					return
				default:
					log.Printf("Multicast read error: %v", err)
					continue
				}
			}

			// Decode the message
			var msg multicast.BroadcastMessage
			if err := json.Unmarshal(buf[:nBytes], &msg); err != nil {
				log.Printf("Failed to decode multicast message: %v", err)
				continue
			}

			// Send the message to the announcements channel, non-blocking.
			select {
			case n.announcements <- msg:
			default:
				log.Println("Multicast announcements channel full; dropping message.")
			}
		}
	}
}
