package gossip

import "eaglechat/common/simplecrypto"

// PublicKeyResponse is the structure returned by the /pubkey endpoint.
// It includes the peer's public key and a signature of the
// marshaled PublicKeyAndCert struct to prove ownership.
type PublicKeyResponse struct {
	PublicKey []byte `json:"public_key"`
	Signature []byte `json:"signature"`
}

// GossipExchangeRequest is the payload sent to the /sync endpoint for a secure exchange.
// It contains the encrypted gossip data and the sender's signature for verification.
type GossipExchangeRequest struct {
	Envelope  *simplecrypto.SecureEnvelope `json:"envelope"`
	Signature []byte                       `json:"signature"`
	PublicKey []byte                       `json:"public_key"`
}

// GossipExchangeResponse is the response from the /sync endpoint for a secure exchange.
// It contains the peer's consolidated gossip data, also encrypted.
type GossipExchangeResponse struct {
	Envelope *simplecrypto.SecureEnvelope `json:"envelope"`
}

// GossipPayload is the internal data structure that is actually exchanged.
// This struct will be marshaled to JSON and then encrypted inside the SecureEnvelope.
type GossipPayload struct {
	KnownPeers  map[string]PeerInfo `json:"known_peers"`
	LastUpdated map[string]int64    `json:"last_updated"`
}

type PeerInfo struct {
	Address string `json:"address"`
	// other metadata
}
