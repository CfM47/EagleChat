package handlers

import (
	"eaglechat/apps/id_manager/internal/application/usecases/gossip"
	"eaglechat/common/simplecrypto/rsa"
	"github.com/gin-gonic/gin"
	"net/http"
)

// PubKeyHandler is the HTTP handler for the /gossip/pubkey endpoint.
// It implements the handlers.Handler interface for gin.
type PubKeyHandler struct {
	privKey *rsa.PrivateKey
	cert    []byte
}

// NewPubKeyHandler creates a new PubKeyHandler.
func NewPubKeyHandler(privKey *rsa.PrivateKey, cert []byte) *PubKeyHandler {
	return &PubKeyHandler{
		privKey: privKey,
		cert:    cert,
	}
}

// Handle handles the GET request for the public key using gin.Context.
func (h *PubKeyHandler) Handle(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.String(http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	pubKey := h.privKey.PublicKey()

	pubKeyBytes, err := pubKey.ToBytes()
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal server error: failed to serialize public key")
		return
	}

	// Concatenate public key bytes and certificate bytes for signing.
	// The rsa.Sign function handles the hashing internally.
	dataToSign := append(pubKeyBytes, h.cert...)

	signature, err := rsa.Sign(dataToSign, h.privKey)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal server error: failed to sign response")
		return
	}

	response := gossip.PublicKeyResponse{
		PublicKey:   pubKeyBytes,
		Certificate: h.cert,
		Signature:   signature,
	}

	c.JSON(http.StatusOK, response)
}
