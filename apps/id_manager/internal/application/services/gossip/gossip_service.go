package gossip

import (
	"bytes"
	"context"
	"eaglechat/apps/id_manager/internal/application/usecases"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"
)

// GossipService is the core component for managing the gossip protocol.
type GossipService struct {
	syncUseCase  *usecases.SyncDataUseCase
	peerProvider PeerProvider
	ownAddress   string // IP:port of this instance
	gossipPort   string // Port for HTTP communication
	httpClient   *http.Client
	ticker       *time.Ticker
	stopChan     chan bool
}

// NewGossipService creates and returns a new GossipService instance.
func NewGossipService(
	syncUseCase *usecases.SyncDataUseCase,
	peerProvider PeerProvider,
	ownAddress string,
	gossipPort string,
	interval time.Duration,
) *GossipService {
	return &GossipService{
		syncUseCase:  syncUseCase,
		peerProvider: peerProvider,
		ownAddress:   ownAddress,
		gossipPort:   gossipPort,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		ticker:   time.NewTicker(interval),
		stopChan: make(chan bool),
	}
}

// Start initiates the periodic gossip synchronization.
func (s *GossipService) Start() {
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.performPeriodicSync()
			case <-s.stopChan:
				s.ticker.Stop()
				return
			}
		}
	}()
	log.Println("GossipService started")
}

// Stop terminates the periodic gossip synchronization.
func (s *GossipService) Stop() {
	close(s.stopChan)
	log.Println("GossipService stopped")
}

// NotifyPeersOfUpdate sends a notification to a random peer to trigger a data pull.
func (s *GossipService) NotifyPeersOfUpdate(ctx context.Context) error {
	peers, err := s.peerProvider.DiscoverIDManagerIPs(ctx)
	if err != nil {
		return fmt.Errorf("failed to discover peers for notification: %w", err)
	}

	peerIP, ok := s.selectRandomPeer(peers)
	if !ok {
		return fmt.Errorf("no valid peers to notify")
	}

	notificationURL := fmt.Sprintf("http://%s:%s/notify-update", peerIP, s.gossipPort)
	requestBody, err := json.Marshal(map[string]string{"source_address": s.ownAddress})
	if err != nil {
		return fmt.Errorf("failed to marshal notification request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", notificationURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create notification request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send notification to %s: %w", notificationURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("notification to %s returned non-OK status: %d", notificationURL, resp.StatusCode)
	}

	log.Printf("Successfully notified peer %s of updates", peerIP)
	return nil
}

// TriggerSyncFromPeer initiates a data pull from a specific peer.
func (s *GossipService) TriggerSyncFromPeer(ctx context.Context, peerAddress string) error {
	targetURL := fmt.Sprintf("http://%s/sync", peerAddress)
	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create sync request for %s: %w", targetURL, err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform sync GET request to %s: %w", targetURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sync GET request to %s returned non-OK status: %d, body: %s", targetURL, resp.StatusCode, string(bodyBytes))
	}

	var peerData usecases.SyncData
	if err := json.NewDecoder(resp.Body).Decode(&peerData); err != nil {
		return fmt.Errorf("failed to decode sync response from %s: %w", targetURL, err)
	}

	if err := s.syncUseCase.MergeData(ctx, peerData); err != nil {
		return fmt.Errorf("failed to merge data from %s: %w", targetURL, err)
	}

	log.Printf("Successfully triggered sync from %s and merged data.", targetURL)
	return nil
}

func (s *GossipService) performPeriodicSync() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	peers, err := s.peerProvider.DiscoverIDManagerIPs(ctx)
	if err != nil {
		log.Printf("GossipService: Failed to discover peers: %v", err)
		return
	}

	peerIP, ok := s.selectRandomPeer(peers)
	if !ok {
		log.Println("GossipService: No valid peers to sync with.")
		return
	}

	peerAddress := fmt.Sprintf("%s:%s", peerIP, s.gossipPort)
	if err := s.TriggerSyncFromPeer(ctx, peerAddress); err != nil {
		log.Printf("GossipService: Periodic sync failed: %v", err)
	}
}

func (s *GossipService) selectRandomPeer(peers []net.IP) (string, bool) {
	var validPeers []string
	ownHost, _, err := net.SplitHostPort(s.ownAddress)
	if err != nil {
		ownHost = s.ownAddress
	}

	for _, peerIP := range peers {
		peerAddrStr := peerIP.String()
		if peerAddrStr != ownHost {
			validPeers = append(validPeers, peerAddrStr)
		}
	}

	if len(validPeers) == 0 {
		return "", false
	}
	return validPeers[rand.Intn(len(validPeers))], true
}
