package simplecrypto_test

import (
	"eaglechat/apps/client/internal/utils/simplecrypto"
	"eaglechat/apps/client/internal/utils/simplecrypto/rsa"
	"testing"

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
		decryptedMessage, senderPubKey, err := simplecrypto.Open(envelope, privBob)
		assert.NoError(t, err)
		assert.Equal(t, originalMessage, decryptedMessage)

		// 4. Verify the returned public key matches the original sender's public key
		assert.Equal(t, pubAlice.Key, senderPubKey.Key)
	})

	t.Run("Open Fails With Wrong Recipient Key", func(t *testing.T) {
		// 1. Generate keys for Alice, Bob, and a malicious actor (Charlie)
		privAlice, _, err := rsa.GenerateKeyPair()
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
		_, _, err = simplecrypto.Open(envelope, privCharlie)
		assert.Error(t, err)
		assert.ErrorIs(t, err, simplecrypto.ErrDecryptKey, "Error should be of type ErrDecryptKey")
	})

	t.Run("Open Fails With Tampered Signature", func(t *testing.T) {
		// 1. Generate keys for Alice and Bob
		privAlice, _, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)
		privBob, pubBob, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)

		// 2. Alice seals a message for Bob
		originalMessage := []byte("secret message")
		envelope, err := simplecrypto.Seal(originalMessage, privAlice, pubBob)
		assert.NoError(t, err)

		// 3. Tamper with the signature
		envelope.Signature[0] ^= 0xff

		// 4. Bob tries to open it, which should fail signature verification
		_, _, err = simplecrypto.Open(envelope, privBob)
		assert.Error(t, err)
		assert.ErrorIs(t, err, simplecrypto.ErrInvalidSignature, "Error should be of type ErrInvalidSignature")
	})
}
