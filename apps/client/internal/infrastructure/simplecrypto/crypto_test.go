package simplecrypto_test

import (
	"testing"

	"eaglechat/apps/client/internal/infrastructure/simplecrypto"
	"eaglechat/apps/client/internal/infrastructure/simplecrypto/rsa"
	"github.com/stretchr/testify/assert"
)

func TestHybridEncryption(t *testing.T) {
	t.Run("Seal and Open Success", func(t *testing.T) {
		// 1. Generate key pairs for sender (Alice) and recipient (Bob)
		privAlice, pubAlice, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)
		privBob, pubBob, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)

		// 2. Seal the message
		originalMessage := []byte("this is a top secret message for bob")
		envelope, err := simplecrypto.Seal(originalMessage, privAlice, pubBob)
		assert.NoError(t, err)
		assert.NotNil(t, envelope)

		// 3. Open the envelope
		decryptedMessage, err := simplecrypto.Open(envelope, privBob, pubAlice)
		assert.NoError(t, err)
		assert.Equal(t, originalMessage, decryptedMessage)
	})

	t.Run("Open Fails With Wrong Recipient Key", func(t *testing.T) {
		// 1. Generate keys for Alice, Bob, and a malicious actor (Charlie)
		privAlice, pubAlice, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)
		_, pubBob, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)
		privCharlie, _, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)

		// 2. Alice seals a message for Bob
		originalMessage := []byte("secret message")
		envelope, err := simplecrypto.Seal(originalMessage, privAlice, pubBob)
		assert.NoError(t, err)

		// 3. Charlie intercepts and tries to open it, which should fail
		_, err = simplecrypto.Open(envelope, privCharlie, pubAlice)
		assert.Error(t, err)
		assert.ErrorIs(t, err, simplecrypto.ErrDecryptKey, "Error should be of type ErrDecryptKey")
	})

	t.Run("Open Fails With Invalid Signature", func(t *testing.T) {
		// 1. Generate keys for Alice, Bob, and Charlie
		privAlice, _, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)
		privBob, pubBob, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)
		_, pubCharlie, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)

		// 2. Alice seals a message for Bob
		originalMessage := []byte("secret message")
		envelope, err := simplecrypto.Seal(originalMessage, privAlice, pubBob)
		assert.NoError(t, err)

		// 3. Bob tries to open it, but mistakenly uses Charlie's public key for verification
		_, err = simplecrypto.Open(envelope, privBob, pubCharlie)
		assert.Error(t, err)
		assert.ErrorIs(t, err, simplecrypto.ErrInvalidSignature, "Error should be of type ErrInvalidSignature")
	})
}
