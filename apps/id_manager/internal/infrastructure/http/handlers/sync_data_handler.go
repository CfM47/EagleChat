package handlers

import (
	"eaglechat/apps/id_manager/internal/application/usecases"
	"eaglechat/apps/id_manager/internal/application/usecases/gossip"
	"eaglechat/common/ezlog"
	"eaglechat/common/simplecrypto"
	"eaglechat/common/simplecrypto/rsa"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SyncDataHandler handles POST requests for /sync endpoint for secure gossip exchange.
type SyncDataHandler struct {
	syncDataUC *usecases.SyncDataUseCase
	myPrivKey  *rsa.PrivateKey
	caPubKey   *rsa.PublicKey
}

// NewSyncDataHandler creates a new SyncDataHandler.
func NewSyncDataHandler(
	syncDataUC *usecases.SyncDataUseCase,
	myPrivKey *rsa.PrivateKey,
	caPubKey *rsa.PublicKey,
) *SyncDataHandler {
	return &SyncDataHandler{
		syncDataUC: syncDataUC,
		myPrivKey:  myPrivKey,
		caPubKey:   caPubKey,
	}
}

// Handle implements the handlers.Handler interface for SyncDataHandler.
func (h *SyncDataHandler) Handle(c *gin.Context) {
	ctx := c.Request.Context()

	logCtx := ezlog.NewLoggerContext("Sync Data Handler")

	if c.Request.Method != http.MethodPost {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	var req gossip.GossipExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad Request: " + err.Error()})
		return
	}

	// 1. Open the secure envelope from the peer
	peerGossipBytes, peerPubKeyFromEnvelope, err := simplecrypto.Open(req.Envelope, h.myPrivKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}

	// 2. Verify the peer's signature and extract their public key
	err = rsa.Verify(req.PublicKey, req.Signature, h.caPubKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}

	peerPubKey, rsaErr := rsa.PublicKeyFromBytes(req.PublicKey)
	if rsaErr != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + rsaErr.Error()})
		return
	}

	// 3. Cross-check that the public key from the certificate matches the one from the envelope.
	// We can now compare the two keys directly as they are the same type.
	if peerPubKeyFromEnvelope.Key.N.Cmp(peerPubKey.Key.N) != 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Public key in envelope does not match public key in signature"})
		return
	}

	// 4. Unmarshal the decrypted gossip payload
	var peerGossipPayload usecases.SyncData
	if err := json.Unmarshal(peerGossipBytes, &peerGossipPayload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad Request: Invalid gossip payload"})
		return
	}

	ezlog.Log(logCtx).Infof("SyncDataHandler: Received and verified gossip from peer. Procceeding to merge data.")

	// 4.5 Merge the incoming data into local repositories
	if err := h.syncDataUC.MergeData(ctx, peerGossipPayload); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error merging data"})
		return
	}

	ezlog.Log(logCtx).Infof("SyncDataHandler: Successfully merged gossip data from peer.")

	// 5. Prepare this node's own gossip payload to send back.
	myGossipData, err := h.syncDataUC.GetAllDataForSync(ctx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error preparing sync data"})
		return
	}
	myGossipPayloadBytes, err := json.Marshal(myGossipData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error marshalling own gossip data"})
		return
	}

	// 6. Seal the response gossip in a new envelope for the peer
	responseEnvelope, err := simplecrypto.Seal(myGossipPayloadBytes, h.myPrivKey, peerPubKeyFromEnvelope)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error sealing response"})
		return
	}

	// 7. Send the response
	response := gossip.GossipExchangeResponse{Envelope: responseEnvelope}
	c.JSON(http.StatusOK, response)
}
