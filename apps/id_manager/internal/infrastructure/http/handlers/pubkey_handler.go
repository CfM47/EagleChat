package handlers

import (
	"eaglechat/apps/id_manager/internal/application/usecases/gossip"
	"eaglechat/common/simplecrypto/rsa"
	"github.com/gin-gonic/gin"
	"net/http" // Keep for http.Status... constants
)

// PublicKeyProvider is an interface that abstracts the source of the public key and certificate.
// This makes the handler more testable.
type PublicKeyProvider interface {
	GetPublicKey() *rsa.PublicKey
	GetCertificate() []byte
	GetPrivateKey() *rsa.PrivateKey
}

// PubKeyHandler is the HTTP handler for the /gossip/pubkey endpoint.
// It implements the handlers.Handler interface for gin.
type PubKeyHandler struct {
	keyProvider PublicKeyProvider
}

// NewPubKeyHandler creates a new PubKeyHandler.
func NewPubKeyHandler(provider PublicKeyProvider) *PubKeyHandler {
	return &PubKeyHandler{keyProvider: provider}
}

// Handle handles the GET request for the public key using gin.Context.
func (h *PubKeyHandler) Handle(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.String(http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	pubKey := h.keyProvider.GetPublicKey()
	cert := h.keyProvider.GetCertificate()
	privKey := h.keyProvider.GetPrivateKey()

	pubKeyBytes, err := pubKey.ToBytes()
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal server error: failed to serialize public key")
		return
	}

	// Concatenate public key bytes and certificate bytes for signing.
	// The rsa.Sign function handles the hashing internally.
	dataToSign := append(pubKeyBytes, cert...)

	signature, err := rsa.Sign(dataToSign, privKey)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal server error: failed to sign response")
		return
	}

	response := gossip.PublicKeyResponse{
		PublicKey:   pubKeyBytes,
		Certificate: cert,
		Signature:   signature,
	}

	c.JSON(http.StatusOK, response)
}
