package idmanagerpool

import (
	"context"
	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/apps/client/internal/middleware/domain/services"
	"eaglechat/apps/client/internal/middleware/infrastructure/environment"
	"eaglechat/apps/client/internal/middleware/infrastructure/idmanagerpool/repositories"
	"eaglechat/common/ezlog"
	"eaglechat/common/multicast/implementation"
	multicast "eaglechat/common/multicast/interface"
	"eaglechat/common/simplecrypto/rsa"
	"errors"
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
)

type idManagerPoolImpl struct {
	repository         repositories.IDManagerRepository
	privateKey         rsa.PrivateKey
	connectionBuilder  middleware_entities.IDManagerConnBuilder
	ownID              entities.UserID
	multicastNet       multicast.MulticastNetwork
	defaultManagerData *middleware_entities.IDManagerData
}

// BuildIDManagerPool creates a new IDManagerPool, initializes the repository, and starts
// processing multicast announcements.
func BuildIDManagerPool(privateKey rsa.PrivateKey, connectionBuilder middleware_entities.IDManagerConnBuilder, ownID entities.UserID) (services.IDManagerPool, error) {
	ctx := ezlog.NewLoggerContext("id-manager-pool-build")

	multicastNet, err := implementation.New(MulticastAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize multicast network: %w", err)
	}

	repo := repositories.NewInMemoryIDManagerRepository(ExpirationTime)

	pool := &idManagerPoolImpl{
		repository:        repo,
		privateKey:        privateKey,
		connectionBuilder: connectionBuilder,
		ownID:             ownID,
		multicastNet:      multicastNet,
	}

	data, err := environment.GetDefaultIDManagerData()
	if err != nil {
		ezlog.Log(ctx).Warnf("No default ID manager configured: %v", err)
	} else {
		pool.defaultManagerData = &data
	}

	go pool.processAnnouncements()

	return pool, nil
}

func (p *idManagerPoolImpl) Close() error {
	return p.multicastNet.Close()
}

func (p *idManagerPoolImpl) Done() <-chan struct{} {
	return p.multicastNet.Done()
}

func (p *idManagerPoolImpl) GetAny() (middleware_entities.IDManagerConnection, error) {
	ctx := ezlog.NewLoggerContext("id-manager-pool-get-any")

	managers := p.repository.GetAll()
	if len(managers) == 0 {
		defaultManager, err := p.getDefault(ctx)
		if err == nil {
			return defaultManager, nil
		}
		return nil, fmt.Errorf("no available id managers")
	}

	randomManager := managers[rand.Intn(len(managers))]
	return p.connectionBuilder(randomManager, p.privateKey, p.ownID)
}

func (p *idManagerPoolImpl) GetAll() ([]middleware_entities.IDManagerConnection, error) {
	ctx := ezlog.NewLoggerContext("id-manager-pool-get-all")

	managers := p.repository.GetAll()
	connections := make([]middleware_entities.IDManagerConnection, 0, len(managers))

	for _, manager := range managers {
		conn, err := p.connectionBuilder(manager, p.privateKey, p.ownID)
		if err != nil {
			ezlog.Log(ctx).Warnf("Error connecting to ID Manager %s:%d: %v", manager.IP, manager.Port, err)
			continue
		}
		connections = append(connections, conn)
	}

	if len(connections) == 0 {
		defaultManager, err := p.getDefault(ctx)
		if err == nil {
			return []middleware_entities.IDManagerConnection{defaultManager}, nil
		}
		ezlog.Log(ctx).Warn("No available ID Managers found and no default configured")
		return nil, fmt.Errorf("no available id managers")
	}

	return connections, nil
}

func (p *idManagerPoolImpl) processAnnouncements() {
	// TODO: add context logger
	for msg := range p.multicastNet.Announcements() {
		if msg.Type == multicast.AnnounceIDManager {
			idManagerMsg, err := msg.AsIDManagerMessage()
			if err != nil {
				log.Printf("Error decoding ID manager message: %v", err)
				continue
			}

			port, err := strconv.ParseUint(idManagerMsg.Port, 10, 16)
			if err != nil {
				log.Printf("Error parsing port from broadcast message: %v", err)
				continue
			}

			pk, err := rsa.PublicKeyFromBytes(idManagerMsg.PublicKey)
			if err != nil {
				log.Printf("Error parsing public key from broadcast message: %v", err)
			}
			managerData := middleware_entities.NewIDManagerData(net.ParseIP(idManagerMsg.IP), uint16(port), *pk)

			p.repository.Add(idManagerMsg.ID, managerData)
		}
	}
}

func (p *idManagerPoolImpl) getDefault(ctx context.Context) (middleware_entities.IDManagerConnection, error) {
	data := p.defaultManagerData
	if data == nil {
		ezlog.Log(ctx).Warnf("No default ID manager configured")
		return nil, errors.New("no default id manager configured")
	}

	conn, err := p.connectionBuilder(*data, p.privateKey, p.ownID)
	if err != nil {
		ezlog.Log(ctx).Warnf("Failed to connect to default ID manager: %v", err)
		return nil, err
	}

	return conn, nil
}
