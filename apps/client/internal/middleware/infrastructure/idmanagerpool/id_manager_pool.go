package idmanagerpool

import (
	"context"
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/middleware/domain/services"
	"eaglechat/apps/client/internal/middleware/infrastructure/idmanagerpool/repositories"
	"eaglechat/common/ezlog"
	"eaglechat/common/ns"
	"eaglechat/common/simplecrypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

const (
	// How often to poll DNS for new ID Managers.
	DNSPOllInterval = 15 * time.Second
	// How long to consider an ID Manager valid without a successful poll.
	ExpirationTime = 30 * time.Second
)

type idManagerPoolImpl struct {
	repository        repositories.IDManagerRepository
	privateKey        rsa.PrivateKey
	connectionBuilder middleware_entities.IDManagerConnBuilder
	ownID             entities.UserID
	quitChan          chan struct{}
	doneChan          chan struct{}
}

//  TODO: Backlog

// // MetadataResponse defines the structure of the JSON response from the /metadata endpoint.
// type MetadataResponse struct {
// 	ID        string `json:"id"`
// 	PublicKey []byte `json:"public_key"`
// }

type StatusResponse struct {
	Status string `json:"status"`
}

// BuildIDManagerPool creates a new IDManagerPool, initializes the repository, and starts
// polling for ID Manager instances via DNS.
func BuildIDManagerPool(privateKey rsa.PrivateKey, connectionBuilder middleware_entities.IDManagerConnBuilder, ownID entities.UserID) (services.IDManagerPool, error) {
	repo := repositories.NewInMemoryIDManagerRepository(ExpirationTime)

	pool := &idManagerPoolImpl{
		repository:        repo,
		privateKey:        privateKey,
		connectionBuilder: connectionBuilder,
		ownID:             ownID,
		quitChan:          make(chan struct{}),
		doneChan:          make(chan struct{}),
	}

	ctx := ezlog.NewLoggerContext("id-manager-pool-poll-loop")
	go pool.pollDNSLoop(ctx)

	return pool, nil
}

func (p *idManagerPoolImpl) Close() error {
	close(p.quitChan)
	<-p.doneChan
	return nil
}

func (p *idManagerPoolImpl) Done() <-chan struct{} {
	return p.doneChan
}

func (p *idManagerPoolImpl) GetAny() (middleware_entities.IDManagerConnection, error) {
	ctx := ezlog.NewLoggerContext("id-manager-pool-get-any")
	ezlog.Log(ctx).Info("Fetching any available ID Manager connection")

	managers := p.repository.GetAll()
	if len(managers) == 0 {
		ezlog.Log(ctx).Warn("No available ID Managers found in repository, attempting to fetch default")
		return nil, fmt.Errorf("no available id managers")
	}

	randomManager := managers[rand.Intn(len(managers))]
	ezlog.Log(ctx).Infof("Selected ID Manager at %s:%d", randomManager.IP, randomManager.Port)

	return p.connectionBuilder(randomManager, p.privateKey, p.ownID)
}

func (p *idManagerPoolImpl) GetAll() ([]middleware_entities.IDManagerConnection, error) {
	ctx := ezlog.NewLoggerContext("id-manager-pool-get-all")
	ezlog.Log(ctx).Info("Fetching all available ID Manager connections")

	managers := p.repository.GetAll()

	ezlog.Log(ctx).Infof("Found %d available ID Managers in repository", len(managers))

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
		ezlog.Log(ctx).Warn("No available ID Managers found and no default configured")
		return nil, fmt.Errorf("no available id managers")
	}

	ezlog.Log(ctx).Infof("Successfully connected to %d ID Managers", len(connections))

	return connections, nil
}

func (p *idManagerPoolImpl) pollDNSLoop(ctx context.Context) {
	ezlog.Log(ctx).Info("Starting DNS polling loop for ID Managers...")

	defer close(p.doneChan)
	defer ezlog.Log(ctx).Info("Stopped DNS polling loop for ID Managers.")

	ticker := time.NewTicker(DNSPOllInterval)
	defer ticker.Stop()

	// Poll once immediately on startup
	pollCtx := ezlog.NewLoggerContext("id-manager-pool-poll")
	p.pollDNS(pollCtx)

	for {
		select {
		case <-ticker.C:
			pollCtx = ezlog.NewLoggerContext("id-manager-pool-poll")
			p.pollDNS(pollCtx)
		case <-p.quitChan:
			return
		}
	}
}

func (p *idManagerPoolImpl) pollDNS(ctx context.Context) {
	ips, err := ns.NewDNSDiscovery().DiscoverIDManagerIPs(ctx)
	if err != nil {
		ezlog.Log(ctx).Warnf("DNS lookup for ID managers failed: %v", err)
		return
	}

	for _, ip := range ips {
		checkHealth(ctx, ip, middleware_entities.DefaultIDManagerPort)

		// TODO: Backlog
		//
		// pk, err := rsa.PublicKeyFromBytes(metadata.PublicKey)
		// if err != nil {
		// 	log.Printf("Error parsing public key from metadata: %v", err)
		// 	continue
		// }

		// port, _ := strconv.ParseUint(IDManagerMetadataPort, 10, 16)
		// managerData := middleware_entities.NewIDManagerData(ip, uint16(port), *pk)

		p.repository.Add(middleware_entities.NewIDManagerData(ip, middleware_entities.DefaultIDManagerPort))
	}
}

func checkHealth(ctx context.Context, ip net.IP, port uint16) bool {
	url := fmt.Sprintf("http://%s:%s/status", ip.String(), middleware_entities.DefaultIDManagerPort)
	ezlog.Log(ctx).Infof("Fetching status from %s", url)

	resp, err := http.Get(url)
	if err != nil {
		ezlog.Log(ctx).Warnf("Error fetching status from %s: %v", url, err)
		return false
	} else {
		ezlog.Log(ctx).Infof("Received message from %s: %v", url, resp)
	}

	defer resp.Body.Close()

	var status StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		log.Printf("Error decoding metadata from %s: %v", url, err)
		return false
	}
	if status.Status != "ok" {
		ezlog.Log(ctx).Warnf("ID Manager at %s returned non-ok status: %s", url, status.Status)
		return false
	}

	return true
}

