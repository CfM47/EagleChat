package p2pconn

import (
	"bufio"
	"eaglechat/apps/client/internal/middleware/domain/entities"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	maxMessageSize = 10 * 1024 * 1024 // 10 MB
	ackTimeout     = 10 * time.Second
)

// p2pConnection is the concrete implementation of the entities.P2PConnection interface.
type p2pConnection struct {
	conn     net.Conn
	wg       sync.WaitGroup
	outgoing chan []byte
	incoming chan []byte
	quit     chan struct{}
	done     chan struct{}

	pendingAcks   map[string]chan error
	pendingAcksMu sync.Mutex
}

// Dial creates and returns a new P2P connection to the given address.
func Dial(ip string, port uint16) (entities.P2PConnection, error) {
	addr := net.JoinHostPort(ip, fmt.Sprintf("%d", port))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return newP2PConnection(conn), nil
}

// newP2PConnection is the internal constructor for a connection.
func newP2PConnection(conn net.Conn) *p2pConnection {
	p := &p2pConnection{
		conn:        conn,
		outgoing:    make(chan []byte, 64),
		incoming:    make(chan []byte),
		quit:        make(chan struct{}),
		done:        make(chan struct{}),
		pendingAcks: make(map[string]chan error),
	}

	p.wg.Add(2)
	go p.readLoop()
	go p.writeLoop()

	// This goroutine waits for the loops to finish and then closes the done channel.
	go func() {
		p.wg.Wait()
		close(p.done)
	}()

	return p
}

func (p *p2pConnection) readLoop() {
	defer p.wg.Done()
	defer p.Close()
	reader := bufio.NewReader(p.conn)

	for {
		var length uint32
		if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
			return // Error or EOF
		}

		if length > maxMessageSize {
			return // Message too large, disconnect
		}

		buf := make([]byte, length)
		if _, err := io.ReadFull(reader, buf); err != nil {
			return // Error or EOF
		}

		var msg WireMessage
		if err := json.Unmarshal(buf, &msg); err != nil {
			log.Printf("p2pconn: failed to unmarshal message: %v", err)
			continue
		}

		switch msg.Type {
		case MsgTypeData:
			// First, forward the payload to the application. This blocks until the
			// application consumes the message from the channel.
			select {
			case p.incoming <- msg.Payload:
				// After consumption, send the ACK.
				ackMsg := WireMessage{Type: MsgTypeAck, ID: msg.ID}
				ackBytes, err := json.Marshal(ackMsg)
				if err == nil {
					// Use a non-blocking select to avoid deadlocking the readLoop
					// if the outgoing channel buffer is full.
					select {
					case p.outgoing <- ackBytes:
					default:
						log.Printf("p2pconn: failed to send ACK for message %s: outgoing buffer full", msg.ID)
					}
				}

			case <-p.quit:
				return
			}

		case MsgTypeAck:
			p.pendingAcksMu.Lock()
			if ackChan, ok := p.pendingAcks[msg.ID]; ok {
				ackChan <- nil // Signal that the ACK was received
			}
			p.pendingAcksMu.Unlock()
		}
	}
}

func (p *p2pConnection) writeLoop() {
	defer p.wg.Done()
	defer p.conn.Close()

	for {
		select {
		case msg := <-p.outgoing:
			// Prepend message with its length
			length := uint32(len(msg))
			if err := binary.Write(p.conn, binary.BigEndian, length); err != nil {
				return
			}
			if _, err := p.conn.Write(msg); err != nil {
				return
			}
		case <-p.quit:
			return
		}
	}
}

func (p *p2pConnection) Send(data []byte) error {
	msgID := uuid.New().String()
	ackChan := make(chan error, 1)

	p.pendingAcksMu.Lock()
	p.pendingAcks[msgID] = ackChan
	p.pendingAcksMu.Unlock()

	defer func() {
		p.pendingAcksMu.Lock()
		delete(p.pendingAcks, msgID)
		p.pendingAcksMu.Unlock()
	}()

	msg := WireMessage{
		Type:    MsgTypeData,
		ID:      msgID,
		Payload: data,
	}

	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	select {
	case p.outgoing <- jsonBytes:
		// Message is on its way, now wait for the ACK.
	case <-p.quit:
		return errors.New("connection is closed")
	}

	select {
	case err := <-ackChan:
		return err // Will be nil on success
	case <-time.After(ackTimeout):
		return errors.New("send confirmation timed out")
	}
}

func (p *p2pConnection) Receive() <-chan []byte {
	return p.incoming
}

func (p *p2pConnection) Close() {
	// Closing the quit channel is the idiomatic way to signal all loops to exit.
	// We use a sync.Once to prevent a panic from closing it multiple times.
	var once sync.Once
	once.Do(func() {
		close(p.quit)
	})
}

func (p *p2pConnection) Done() <-chan struct{} {
	return p.done
}
