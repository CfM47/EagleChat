package idmanagerregisterer

import (
	"context"
	multicast "eaglechat/common/multicast/interface"
	"errors"
	"log"
)

// listenForIDManager waits for an ID Manager announcement.
func (r *registererImpl) listenForIDManager(ctx context.Context, multicastNet multicast.MulticastNetwork, idManagerChan chan<- multicast.IDManagerMessage, errChan chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-multicastNet.Announcements():
			if !ok {
				errChan <- errors.New("multicast announcements channel closed unexpectedly")
				return
			}

			if msg.Type == multicast.AnnounceIDManager {
				idManagerMsg, err := msg.AsIDManagerMessage()
				if err != nil {
					log.Printf("Received malformed ID Manager announcement: %v", err)
					continue
				}
				idManagerChan <- idManagerMsg
				return
			}
		}
	}
}
