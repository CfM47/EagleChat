package aes_test

import (
	"testing"

	"eaglechat/apps/client/internal/infrastructure/simplecrypto/aes"
	"github.com/stretchr/testify/assert"
)

func TestAESEncryption(t *testing.T) {
	t.Run("Encrypt and Decrypt Success", func(t *testing.T) {
		key, err := aes.GenerateKey()
		assert.NoError(t, err)
		assert.Len(t, key, aes.KeySize)

		originalText := []byte("this is a very secret message that needs to be encrypted with aes")

		ciphertext, err := aes.Encrypt(originalText, key)
		assert.NoError(t, err)
		assert.NotEmpty(t, ciphertext)
		assert.NotEqual(t, originalText, ciphertext)

		plaintext, err := aes.Decrypt(ciphertext, key)
		assert.NoError(t, err)
		assert.Equal(t, originalText, plaintext)
	})

	t.Run("Decrypt Fails With Wrong Key", func(t *testing.T) {
		keyA, err := aes.GenerateKey()
		assert.NoError(t, err)
		keyB, err := aes.GenerateKey()
		assert.NoError(t, err)

		originalText := []byte("secret message")
		ciphertext, err := aes.Encrypt(originalText, keyA)
		assert.NoError(t, err)

		_, err = aes.Decrypt(ciphertext, keyB)
		assert.Error(t, err, "Decryption with wrong key should fail")
	})

	t.Run("Decrypt Fails With Corrupted Ciphertext", func(t *testing.T) {
		key, err := aes.GenerateKey()
		assert.NoError(t, err)

		originalText := []byte("another secret")
		ciphertext, err := aes.Encrypt(originalText, key)
		assert.NoError(t, err)

		// Tamper with the ciphertext (flip a bit)
		ciphertext[len(ciphertext)-1] ^= 0x01

		_, err = aes.Decrypt(ciphertext, key)
		assert.Error(t, err, "Decryption of corrupted data should fail")
	})
}
