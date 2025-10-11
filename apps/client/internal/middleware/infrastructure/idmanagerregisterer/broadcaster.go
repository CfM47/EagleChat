package idmanagerregisterer

import (
	"context"
	multicast "eaglechat/common/multicast/interface"
	"fmt"
	"time"
)

// broadcastLoop periodically sends out a client registration message.
func (r *registererImpl) broadcastLoop(ctx context.Context, multicastNet multicast.MulticastNetwork, errChan chan<- error) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	registerMsg := multicast.NewClientRegisterMessage()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := multicastNet.Broadcast(registerMsg); err != nil {
				errChan <- fmt.Errorf("failed to broadcast registration message: %w", err)
				return
			}
		}
	}
}
