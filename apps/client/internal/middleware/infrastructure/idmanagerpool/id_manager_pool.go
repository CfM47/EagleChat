package idmanagerpool

import (
	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/apps/client/internal/middleware/domain/services"
	"eaglechat/apps/client/internal/middleware/infrastructure/idmanagerpool/repositories"
	"eaglechat/apps/client/internal/utils/simplecrypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

const (
	MulticastAddress = "239.0.0.1:9999"
	ExpirationTime   = 30 * time.Second
	maxDatagramSize  = 8192
)

// BroadcastMessage defines the structure of the announcement message sent via multicast.
type BroadcastMessage struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	IP   string `json:"ip"`
	Port string `json:"port"`
	Time int64  `json:"time"`
}

type idManagerPoolImpl struct {
	repository        repositories.IDManagerRepository
	privateKey        rsa.PrivateKey
	connectionBuilder middleware_entities.IDManagerConnBuilder
	ownID             entities.UserID
}

// BuildIDManagerPool creates a new IDManagerPool, initializes the repository, and starts
// listening for multicast announcements.
func BuildIDManagerPool(privateKey rsa.PrivateKey, connectionBuilder middleware_entities.IDManagerConnBuilder, ownID entities.UserID) (services.IDManagerPool, error) {
	repo := repositories.NewInMemoryIDManagerRepository(ExpirationTime)

	pool := &idManagerPoolImpl{
		repository:        repo,
		privateKey:        privateKey,
		connectionBuilder: connectionBuilder,
	}

	go pool.listenForAnnouncements()

	return pool, nil
}

func (p *idManagerPoolImpl) GetAny() (middleware_entities.IDManagerConnection, error) {
	managers := p.repository.GetAll()
	if len(managers) == 0 {
		return nil, fmt.Errorf("no available id managers")
	}

	// For now, we don't have the client's private key, so we pass nil.
	// This will need to be addressed when authentication is implemented.
	randomManager := managers[rand.Intn(len(managers))]
	return p.connectionBuilder(randomManager, p.privateKey, p.ownID)
}

func (p *idManagerPoolImpl) GetAll() ([]middleware_entities.IDManagerConnection, error) {
	managers := p.repository.GetAll()
	connections := make([]middleware_entities.IDManagerConnection, 0, len(managers))

	for _, manager := range managers {
		// For now, we don't have the client's private key, so we pass nil.
		conn, err := p.connectionBuilder(manager, p.privateKey, p.ownID)
		if err != nil {
			// Log the error and continue to the next manager.
			log.Printf("Error connecting to ID Manager %s:%d: %v", manager.IP, manager.Port, err)
			continue
		}
		connections = append(connections, conn)
	}

	if len(connections) == 0 {
		return nil, fmt.Errorf("no available id managers")
	}

	return connections, nil
}

func (p *idManagerPoolImpl) listenForAnnouncements() {
	addr, err := net.ResolveUDPAddr("udp", MulticastAddress)
	if err != nil {
		log.Fatalf("Error resolving multicast address: %v", err)
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Error listening to multicast address: %v", err)
	}
	defer conn.Close()

	conn.SetReadBuffer(maxDatagramSize)

	buffer := make([]byte, maxDatagramSize)
	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from multicast UDP: %v", err)
			continue
		}

		var msg BroadcastMessage
		if err := json.Unmarshal(buffer[:n], &msg); err != nil {
			log.Printf("Error unmarshalling broadcast message: %v", err)
			continue
		}

		if msg.Type == "ANNOUNCE" {
			port, err := strconv.ParseUint(msg.Port, 10, 16)
			if err != nil {
				log.Printf("Error parsing port from broadcast message: %v", err)
				continue
			}

			managerData := middleware_entities.IDManagerData{
				IP:   net.ParseIP(msg.IP),
				Port: uint16(port),
			}

			p.repository.Add(msg.ID, managerData)
		}
	}
}
