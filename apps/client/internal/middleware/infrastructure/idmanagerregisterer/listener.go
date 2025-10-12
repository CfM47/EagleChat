package idmanagerregisterer

import (
	"context"
	"eaglechat/common/ezlog"
	multicast "eaglechat/common/multicast/interface"
	"errors"
)

// listenForIDManager waits for an ID Manager announcement.
func (r *registererImpl) listenForIDManager(ctx context.Context, multicastNet multicast.MulticastNetwork, idManagerChan chan<- multicast.IDManagerMessage, errChan chan<- error) {
	ezlog.Log(ctx).Info("registration listener loop started")

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-multicastNet.Announcements():
			if !ok {
				msg := "registration multicast announcements channel closed unexpectedly"
				ezlog.Log(ctx).Error(msg)
				errChan <- errors.New(msg)
				return
			}

			if msg.Type == multicast.AnnounceIDManager {
				idManagerMsg, err := msg.AsIDManagerMessage()
				if err != nil {
					ezlog.Log(ctx).Debugf("Received malformed ID Manager announcement: %v", err)
					continue
				}
				idManagerChan <- idManagerMsg
				return
			}
		}
	}
}
