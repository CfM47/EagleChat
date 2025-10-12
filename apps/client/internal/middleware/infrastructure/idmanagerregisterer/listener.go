package idmanagerregisterer

import (
	"context"
	"eaglechat/apps/client/internal/middleware/infrastructure/environment"
	"eaglechat/common/ezlog"
	multicast "eaglechat/common/multicast/interface"
	"errors"
	"net"
	"strconv"
)

// listenForIDManager waits for an ID Manager announcement.
func (r *registererImpl) listenForIDManager(ctx context.Context, multicastNet multicast.MulticastNetwork, idManagerChan chan<- IDManagerData, errChan chan<- error) {
	ezlog.Log(ctx).Info("registration listener loop started")

	for {
		select {
		case <-ctx.Done():
			return
		case <-r.foundIDManager:
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

				ip := net.ParseIP(idManagerMsg.IP)
				if ip == nil {
					ezlog.Log(ctx).Debugf("Received ID Manager announcement with invalid IP: %s", idManagerMsg.IP)
					continue
				}

				port, err := strconv.ParseUint(idManagerMsg.Port, 10, 16)
				if err != nil {
					ezlog.Log(ctx).Debugf("Received ID Manager announcement with invalid port: %s", idManagerMsg.Port)
					continue
				}

				idManagerChan <- IDManagerData{
					IP:   net.IP(idManagerMsg.IP),
					Port: uint16(port),
				}
				return
			}
		}
	}
}

func (r *registererImpl) tryDefaultIDManager(ctx context.Context, idManagerChan chan<- IDManagerData) {
	data, err := environment.GetDefaultIDManagerData()
	if err != nil {
		ezlog.Log(ctx).Info("No default ID Manager configured")
		return
	}

	select {
	case <-ctx.Done():
		return
	case idManagerChan <- IDManagerData{
		IP:   data.IP,
		Port: data.Port,
	}:
		return
	}
}
