package rsa_test

import (
	"testing"

	"eaglechat/apps/client/internal/utils/simplecrypto/rsa"
	"github.com/stretchr/testify/assert"
)

func TestEncryption(t *testing.T) {
	t.Run("Message can be encrypted and decrypted successfully", func(t *testing.T) {
		privKey, pubKey, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)

		originalMessage := []byte("this is a secret message")
		ciphertext, err := rsa.Encrypt(originalMessage, pubKey)
		assert.NoError(t, err)

		decryptedMessage, err := rsa.Decrypt(ciphertext, privKey)
		assert.NoError(t, err)
		assert.Equal(t, originalMessage, decryptedMessage)
	})

	t.Run("Decrypting with the wrong key fails", func(t *testing.T) {
		_, pubKeyA, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)
		privKeyB, _, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)

		originalMessage := []byte("this is a secret message")
		ciphertext, err := rsa.Encrypt(originalMessage, pubKeyA)
		assert.NoError(t, err)

		_, err = rsa.Decrypt(ciphertext, privKeyB)
		assert.Error(t, err)
	})
}
