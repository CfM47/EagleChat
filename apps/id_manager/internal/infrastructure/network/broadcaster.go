package network

import (
	"encoding/json"
	"net"
	"time"
)

// BroadcastMessage defines the structure for multicast announcements.
type BroadcastMessage struct {
	Type string `json:"type"` // "ANNOUNCE" or "UPDATED"
	ID   string `json:"id"`
	IP   string `json:"ip"`
	Time int64  `json:"time"`
}

// Broadcaster periodically announces the ID Manager's presence via multicast.
type Broadcaster struct {
	id        string
	ip        string
	groupAddr string
	interval  time.Duration
	conn      *net.UDPConn
	stopChan  chan struct{}
}

// NewBroadcaster initializes a new Broadcaster.
func NewBroadcaster(id, ip, groupAddr string, interval time.Duration) (*Broadcaster, error) {
	addr, err := net.ResolveUDPAddr("udp", groupAddr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return &Broadcaster{
		id:        id,
		ip:        ip,
		groupAddr: groupAddr,
		interval:  interval,
		conn:      conn,
		stopChan:  make(chan struct{}),
	}, nil
}

// Start begins periodic announcements.
func (b *Broadcaster) Start() {
	ticker := time.NewTicker(b.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				msg := BroadcastMessage{
					Type: "ANNOUNCE",
					ID:   b.id,
					IP:   b.ip,
					Time: time.Now().Unix(),
				}
				data, _ := json.Marshal(msg)
				b.conn.Write(data)
			case <-b.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop halts the broadcaster.
func (b *Broadcaster) Stop() {
	close(b.stopChan)
	b.conn.Close()
}

// BroadcastUpdated sends an "UPDATED" message to the multicast group.
func (b *Broadcaster) BroadcastUpdated() {
	msg := BroadcastMessage{
		Type: "UPDATED",
		ID:   b.id,
		IP:   b.ip,
		Time: time.Now().Unix(),
	}
	data, _ := json.Marshal(msg)
	b.conn.Write(data)
}
