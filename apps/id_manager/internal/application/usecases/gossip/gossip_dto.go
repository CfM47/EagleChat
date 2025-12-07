package gossip

import "eaglechat/common/simplecrypto"

// PublicKeyResponse is the structure returned by the /gossip/pubkey endpoint.
// It includes the peer's public key, its certificate, and a signature of the
// marshaled PublicKeyAndCert struct to prove ownership.
type PublicKeyResponse struct {
	PublicKey   []byte `json:"public_key"`
	Certificate []byte `json:"certificate"`
	Signature   []byte `json:"signature"` // Signature of SHA256(Marshal(PublicKey + Certificate))
}

// GossipExchangeRequest is the payload sent to the /sync endpoint for a secure exchange.
// It contains the encrypted gossip data and the sender's certificate for verification.
type GossipExchangeRequest struct {
	Envelope    *simplecrypto.SecureEnvelope `json:"envelope"`
	Certificate []byte                       `json:"certificate"`
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
