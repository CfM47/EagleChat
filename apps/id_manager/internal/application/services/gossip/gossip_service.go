package gossip

import (
	"bytes"
	"context"
	"eaglechat/apps/id_manager/internal/application/usecases"
	"eaglechat/apps/id_manager/internal/application/usecases/gossip"
	"eaglechat/common/ezlog"
	"eaglechat/common/simplecrypto"
	"eaglechat/common/simplecrypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"sync"
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

	// Crypto dependencies for secure communication
	myPubKey    *rsa.PublicKey
	myPrivKey   *rsa.PrivateKey
	mySignature []byte
	caPubKey    *rsa.PublicKey

	// Logger context
	logCtx context.Context
}

// NewGossipService creates and returns a new GossipService instance.
func NewGossipService(
	syncUseCase *usecases.SyncDataUseCase,
	peerProvider PeerProvider,
	ownAddress string,
	gossipPort string,
	interval time.Duration,
	myPubKey *rsa.PublicKey,
	myPrivKey *rsa.PrivateKey,
	mySignature []byte,
	caPubKey *rsa.PublicKey,
) *GossipService {
	return &GossipService{
		syncUseCase:  syncUseCase,
		peerProvider: peerProvider,
		ownAddress:   ownAddress,
		gossipPort:   gossipPort,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		ticker:      time.NewTicker(interval),
		stopChan:    make(chan bool),
		myPubKey:    myPubKey,
		myPrivKey:   myPrivKey,
		mySignature: mySignature,
		caPubKey:    caPubKey,
		logCtx:      ezlog.NewLoggerContext("gossip service"),
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
	ezlog.Log(s.logCtx).Info("GossipService started")
}

// Stop terminates the periodic gossip synchronization.
func (s *GossipService) Stop() {
	close(s.stopChan)
	ezlog.Log(s.logCtx).Info("GossipService stopped")
}

// NotifyPeersOfUpdate sends a notification to a random peer to trigger a data pull.
func (s *GossipService) NotifyPeersOfUpdate(ctx context.Context) error {
	peers, err := s.peerProvider.DiscoverIDManagerIPs(ctx)
	if err != nil {
		return fmt.Errorf("failed to discover peers for notification: %w", err)
	}

	var wg sync.WaitGroup
	for _, peerIP := range s.getAllValidPeers(peers) {
		wg.Add(1)
		go func(peerIP string) {
			defer wg.Done()
			if err := s.notifyPeer(ctx, peerIP); err != nil {
				ezlog.Log(s.logCtx).Errorf("Failed to notify peer %s: %v", peerIP, err)
			}
		}(peerIP)
	}
	wg.Wait()

	ezlog.Log(s.logCtx).Infof("Successfully notified of updates to peers %v", peers)
	return nil
}

// TriggerSyncFromPeer orchestrates the two-phase secure gossip exchange with a specific peer.
func (s *GossipService) TriggerSyncFromPeer(ctx context.Context, peerAddress string) error {
	// Phase 1: Get peer's public key and verify identity
	peerPubKey, err := s.fetchAndVerifyPeerPublicKey(ctx, peerAddress)
	if err != nil {
		return fmt.Errorf("phase 1 (public key discovery) failed with %s: %w", peerAddress, err)
	}

	// Phase 2: Perform the secure data exchange
	if err := s.performSecureGossipExchange(ctx, peerAddress, peerPubKey); err != nil {
		return fmt.Errorf("phase 2 (secure gossip exchange) failed with %s: %w", peerAddress, err)
	}

	ezlog.Log(s.logCtx).Infof("Successfully completed secure sync with %s", peerAddress)
	return nil
}

func (s *GossipService) notifyPeer(ctx context.Context, peerIP string) error {
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

	ezlog.Log(s.logCtx).Infof("Successfully notified peer %s of updates", peerIP)
	return nil
}

// fetchAndVerifyPeerPublicKey performs Phase 1 of the secure gossip exchange.
// It fetches the peer's public key and certificate, then verifies them.
// Returns the trusted peer's public key.
func (s *GossipService) fetchAndVerifyPeerPublicKey(ctx context.Context, peerAddress string) (*rsa.PublicKey, error) {
	pubKeyURL := fmt.Sprintf("http://%s/gossip/pubkey", peerAddress)
	req, err := http.NewRequestWithContext(ctx, "GET", pubKeyURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create public key request for %s: %w", pubKeyURL, err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform public key GET request to %s: %w", pubKeyURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("public key GET request to %s returned non-OK status: %d, body: %s", pubKeyURL, resp.StatusCode, string(bodyBytes))
	}

	var pubKeyResp gossip.PublicKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&pubKeyResp); err != nil {
		return nil, fmt.Errorf("failed to decode public key response from %s: %w", pubKeyURL, err)
	}

	// Verify the peer's certificate against the trusted CA
	err = rsa.Verify(pubKeyResp.PublicKey, pubKeyResp.Signature, s.caPubKey)
	if err != nil {
		return nil, fmt.Errorf("peer certificate verification failed for %s: %w", pubKeyURL, err)
	}

	peerPubKey, rsaErr := rsa.PublicKeyFromBytes(pubKeyResp.PublicKey)
	if rsaErr != nil {
		return nil, fmt.Errorf("failed to deserialize peer public key %s: %w", pubKeyURL, err)
	}

	return peerPubKey, nil
}

// performSecureGossipExchange performs Phase 2 of the secure gossip exchange.
// It sends this node's gossip encrypted to the peer and processes the peer's encrypted response.
func (s *GossipService) performSecureGossipExchange(ctx context.Context, peerAddress string, peerPubKey *rsa.PublicKey) error {
	// Get our own data to send (this node's current state)
	myGossipData, err := s.syncUseCase.GetAllDataForSync(ctx)
	if err != nil {
		return fmt.Errorf("failed to get own gossip data: %w", err)
	}
	myGossipBytes, err := json.Marshal(myGossipData)
	if err != nil {
		return fmt.Errorf("failed to marshal own gossip data: %w", err)
	}

	// Encrypt our data for the peer using the peer's trusted public key
	envelope, err := simplecrypto.Seal(myGossipBytes, s.myPrivKey, peerPubKey)
	if err != nil {
		return fmt.Errorf("failed to seal gossip envelope for %s: %w", peerAddress, err)
	}

	myPubKeyBytes, err := s.myPubKey.ToBytes()
	if err != nil {
		return fmt.Errorf("failed to get own public key bytes: %w", err)
	}

	// Prepare the POST request body
	exchangeReq := gossip.GossipExchangeRequest{
		Envelope:  envelope,
		Signature: s.mySignature, // Send our own certificate for the peer to verify
		PublicKey: myPubKeyBytes,
	}
	reqBodyBytes, err := json.Marshal(exchangeReq)
	if err != nil {
		return fmt.Errorf("failed to marshal gossip exchange request for %s: %w", peerAddress, err)
	}
	reqBodyReader := bytes.NewReader(reqBodyBytes)

	// Make the POST request to the modified /sync endpoint
	syncURL := fmt.Sprintf("http://%s/sync", peerAddress)
	req, err := http.NewRequestWithContext(ctx, "POST", syncURL, reqBodyReader)
	if err != nil {
		return fmt.Errorf("failed to create secure sync request for %s: %w", syncURL, err)
	}
	req.Header.Set("Content-Type", "application/json") // Ensure content type is set

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform secure sync POST request to %s: %w", syncURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("secure sync POST request to %s returned non-OK status: %d, body: %s", syncURL, resp.StatusCode, string(bodyBytes))
	}

	// Decrypt and process the response from the peer
	var exchangeResp gossip.GossipExchangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&exchangeResp); err != nil {
		return fmt.Errorf("failed to decode secure sync response from %s: %w", syncURL, err)
	}

	// Open the peer's response envelope. This decrypts the data and verifies their signature.
	peerResponseDataBytes, _, err := simplecrypto.Open(exchangeResp.Envelope, s.myPrivKey)
	if err != nil {
		return fmt.Errorf("failed to open peer's response envelope from %s: %w", syncURL, err)
	}

	// Unmarshal the peer's data and merge it
	var peerResponsePayload usecases.SyncData // Assuming SyncDataUseCase expects usecases.SyncData
	if err := json.Unmarshal(peerResponseDataBytes, &peerResponsePayload); err != nil {
		return fmt.Errorf("failed to unmarshal peer's response payload from %s: %w", syncURL, err)
	}

	if err := s.syncUseCase.MergeData(ctx, peerResponsePayload); err != nil {
		return fmt.Errorf("failed to merge data from %s: %w", syncURL, err)
	}

	return nil
}

func (s *GossipService) performPeriodicSync() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	peers, err := s.peerProvider.DiscoverIDManagerIPs(ctx)
	if err != nil {
		ezlog.Log(s.logCtx).Errorf("GossipService: Failed to discover peers: %v", err)
		return
	}

	peerIP, ok := s.selectRandomPeer(peers)
	if !ok {
		ezlog.Log(s.logCtx).Warn("GossipService: No valid peers to sync with.")
		return
	}

	peerAddress := fmt.Sprintf("%s:%s", peerIP, s.gossipPort)
	if err := s.TriggerSyncFromPeer(ctx, peerAddress); err != nil {
		ezlog.Log(s.logCtx).Errorf("GossipService: Periodic sync failed: %v", err)
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

func (s *GossipService) getAllValidPeers(peers []net.IP) []string {
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
	return validPeers
}
