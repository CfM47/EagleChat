package network

import (
	"encoding/json"
	"log"
	"math/rand"
	"net"

	"eaglechat/apps/id_manager/internal/application/usecases"
)

// Listener receives multicast messages from other ID Managers.
type Listener struct {
	groupAddr string
	handler   *usecases.SyncUsersUseCase
	localID   string
	stopChan  chan struct{}
}

// NewListener initializes a multicast listener.
func NewListener(groupAddr, localID string, handler *usecases.SyncUsersUseCase) *Listener {
	return &Listener{
		groupAddr: groupAddr,
		handler:   handler,
		localID:   localID,
		stopChan:  make(chan struct{}),
	}
}

// Start listens for multicast messages and triggers sync operations.
func (l *Listener) Start() {
	addr, err := net.ResolveUDPAddr("udp", l.groupAddr)
	if err != nil {
		log.Fatalf("failed to resolve multicast addr: %v", err)
	}
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("failed to listen multicast: %v", err)
	}
	conn.SetReadBuffer(4096)

	go func() {
		buf := make([]byte, 4096)
		for {
			select {
			case <-l.stopChan:
				conn.Close()
				return
			default:
				n, _, err := conn.ReadFromUDP(buf)
				if err != nil {
					continue
				}
				var msg BroadcastMessage
				if err := json.Unmarshal(buf[:n], &msg); err != nil {
					continue
				}
				if msg.ID == l.localID {
					continue // ignore self
				}
				if msg.Type == "ANNOUNCE" {
					// Randomly decide to initiate sync
					if rand.Float64() < 0.2 {
						go l.handler.SyncWithPeer(msg.IP)
					}
				}
			}
		}
	}()
}

// Stop halts the listener.
func (l *Listener) Stop() {
	close(l.stopChan)
}
