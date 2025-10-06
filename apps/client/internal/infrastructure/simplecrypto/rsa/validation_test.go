package rsa_test

import (
	"testing"

	"eaglechat/apps/client/internal/infrastructure/simplecrypto/rsa"
	"github.com/stretchr/testify/assert"
)

func TestSignatures(t *testing.T) {
	t.Run("Signature can be created and verified successfully", func(t *testing.T) {
		privKey, pubKey, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)

		message := []byte("this message needs to be signed")
		signature, err := rsa.Sign(message, privKey)
		assert.NoError(t, err)

		err = rsa.Verify(message, signature, pubKey)
		assert.NoError(t, err)
	})

	t.Run("Verification fails for a tampered message", func(t *testing.T) {
		privKey, pubKey, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)
		originalMessage := []byte("this is the original content")
		signature, err := rsa.Sign(originalMessage, privKey)
		assert.NoError(t, err)

		tamperedMessage := []byte("this is the tampered content")
		err = rsa.Verify(tamperedMessage, signature, pubKey)
		assert.Error(t, err)
	})

	t.Run("Verification fails with the wrong public key", func(t *testing.T) {
		privKeyA, _, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)
		_, pubKeyB, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)

		message := []byte("some content")
		signature, err := rsa.Sign(message, privKeyA)
		assert.NoError(t, err)

		err = rsa.Verify(message, signature, pubKeyB)
		assert.Error(t, err)
	})
}
