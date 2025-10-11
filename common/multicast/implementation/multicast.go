package implementation

import (
	multicast "eaglechat/common/multicast/interface"
	"net"
	"sync"
)

const (
	defaultAnnouncementsBuffer = 16
	DefaultUDPAddress          = "239.0.0.1:9999"
)

// udpMulticastNetwork implements the MulticastNetwork interface using UDP.
type udpMulticastNetwork struct {
	conn          *net.UDPConn
	addr          *net.UDPAddr
	announcements chan multicast.BroadcastMessage
	done          chan struct{}
	wg            sync.WaitGroup
}

// New creates and initializes a new MulticastNetwork.
func New(address string) (multicast.MulticastNetwork, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	net := &udpMulticastNetwork{
		conn:          conn,
		addr:          addr,
		announcements: make(chan multicast.BroadcastMessage, defaultAnnouncementsBuffer),
		done:          make(chan struct{}),
	}

	net.wg.Add(1)
	go net.listenLoop()

	return net, nil
}

// Announcements returns the channel for receiving broadcast messages.
func (n *udpMulticastNetwork) Announcements() <-chan multicast.BroadcastMessage {
	return n.announcements
}

// Close shuts down the multicast network connection and cleans up resources.
func (n *udpMulticastNetwork) Close() error {
	// Signal goroutines to stop
	close(n.done)

	// Closing the connection will make the listener's ReadFromUDP unblock
	err := n.conn.Close()

	// Wait for all goroutines to finish
	n.wg.Wait()

	// Close the announcements channel after all producers are done
	close(n.announcements)

	return err
}

// Done returns a channel that is closed when the network object is fully shut down.
func (n *udpMulticastNetwork) Done() <-chan struct{} {
	return n.done
}
