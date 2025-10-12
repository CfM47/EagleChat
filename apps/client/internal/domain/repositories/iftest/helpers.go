package iftest

import (
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/apps/client/internal/domain/repositories"
	"eaglechat/common/simplecrypto/rsa"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const keyPoolSize = 5

var (
	keyPairs []*rsa.PrivateKey
	initOnce sync.Once
)

// initKeys generates a pool of RSA key pairs to be used by the test suite.
// This is done once to avoid the expensive key generation operation in each test.
func initKeys() {
	initOnce.Do(func() {
		keyPairs = make([]*rsa.PrivateKey, keyPoolSize)
		for i := 0; i < keyPoolSize; i++ {
			privKey, _, _ := rsa.GenerateKeyPair()
			keyPairs[i] = privKey
		}
	})
}

// getTestKeys returns a pre-generated private key from the pool.
// It panics if the index is out of bounds.
func getTestKeys(index int) *rsa.PrivateKey {
	initKeys() // Ensures keys are generated
	if index < 0 || index >= keyPoolSize {
		panic("test key index out of bounds")
	}
	return keyPairs[index]
}

// --- Chat Generation Helpers ---

type sortByCreatedDate []entities.Message

func (a sortByCreatedDate) Len() int           { return len(a) }
func (a sortByCreatedDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByCreatedDate) Less(i, j int) bool { return a[i].CreatedTime.Before(a[j].CreatedTime) }

// generateChat creates a realistic, randomly interspersed but chronologically ordered chat history.
func generateChat(self, peer entities.User, numSelf, numPeer int) []entities.Message {
	var messages []entities.Message

	// 1. Generate all messages without timestamps
	for i := 0; i < numSelf; i++ {
		messages = append(messages, entities.NewMessage(self, peer, fmt.Sprintf("Message from self %d", i)))
	}
	for i := 0; i < numPeer; i++ {
		messages = append(messages, entities.NewMessage(peer, self, fmt.Sprintf("Message from peer %d", i)))
	}

	// 2. Shuffle them to randomize who sent when
	rand.Shuffle(len(messages), func(i, j int) {
		messages[i], messages[j] = messages[j], messages[i]
	})

	// 3. Assign deterministic, incrementing timestamps
	startTime := time.Now().Add(-1 * time.Hour).Truncate(time.Second)
	for i := range messages {
		messages[i].CreatedTime = startTime.Add(time.Duration(i) * time.Minute)
	}

	// 4. Shuffle again so the returned slice isn't pre-sorted
	rand.Shuffle(len(messages), func(i, j int) {
		messages[i], messages[j] = messages[j], messages[i]
	})

	return messages
}

// insertChat is a test helper to save a slice of messages to the repository.
func insertChat(t *testing.T, repo repositories.ClientRepository, messages []entities.Message) {
	t.Helper()
	for _, msg := range messages {
		require.NoError(t, repo.SaveMessage(msg))
	}
}

// incrementUnreadCount is a test helper to call the repository's IncrementUnreadCount method multiple times.
func incrementUnreadCount(t *testing.T, repo repositories.ClientRepository, partnerID entities.UserID, times int) {
	t.Helper()
	for i := 0; i < times; i++ {
		require.NoError(t, repo.IncrementUnreadCount(partnerID))
	}
}
