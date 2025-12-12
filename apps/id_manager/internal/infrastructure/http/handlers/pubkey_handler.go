package handlers

import (
	"eaglechat/apps/id_manager/internal/application/usecases/gossip"
	"eaglechat/common/simplecrypto/rsa"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PubKeyHandler is the HTTP handler for the /pubkey endpoint.
// It implements the handlers.Handler interface for gin.
type PubKeyHandler struct {
	privKey   *rsa.PrivateKey
	signature []byte
}

// NewPubKeyHandler creates a new PubKeyHandler.
func NewPubKeyHandler(privKey *rsa.PrivateKey, signature []byte) *PubKeyHandler {
	return &PubKeyHandler{
		privKey:   privKey,
		signature: signature,
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

	response := gossip.PublicKeyResponse{
		PublicKey: pubKeyBytes,
		Signature: h.signature,
	}

	c.JSON(http.StatusOK, response)
}
