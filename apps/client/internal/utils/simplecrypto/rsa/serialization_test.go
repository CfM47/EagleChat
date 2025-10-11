package rsa_test

import (
	"os"
	"path/filepath"
	"testing"

	"eaglechat/apps/client/internal/utils/simplecrypto/rsa"
	"github.com/stretchr/testify/assert"
)

func TestKeySerialization(t *testing.T) {
	t.Run("Private Key can be serialized and deserialized", func(t *testing.T) {
		privKey, _, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)

		privBytes := privKey.ToBytes()
		restoredFromBytes, err := rsa.PrivateKeyFromBytes(privBytes)
		assert.NoError(t, err)
		assert.True(t, restoredFromBytes.Key.Equal(privKey.Key))
	})

	t.Run("Public Key can be serialized and deserialized", func(t *testing.T) {
		_, pubKey, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)

		pubBytes, err := pubKey.ToBytes()
		assert.NoError(t, err)
		restoredFromBytes, err := rsa.PublicKeyFromBytes(pubBytes)
		assert.NoError(t, err)
		assert.True(t, restoredFromBytes.Key.Equal(pubKey.Key))
	})
}

func TestFileStorage(t *testing.T) {
	t.Run("Private Key can be written to and read from a file", func(t *testing.T) {
		tempDir := t.TempDir()
		privKey, _, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)
		keyPath := filepath.Join(tempDir, "key.priv")

		err = os.WriteFile(keyPath, privKey.ToBytes(), 0600)
		assert.NoError(t, err)

		readBytes, err := os.ReadFile(keyPath)
		assert.NoError(t, err)
		restoredKey, err := rsa.PrivateKeyFromBytes(readBytes)
		assert.NoError(t, err)
		assert.True(t, restoredKey.Key.Equal(privKey.Key))
	})

	t.Run("Public Key can be written to and read from a file", func(t *testing.T) {
		tempDir := t.TempDir()
		_, pubKey, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)
		keyPath := filepath.Join(tempDir, "key.pub")

		pubBytes, err := pubKey.ToBytes()
		assert.NoError(t, err)
		err = os.WriteFile(keyPath, pubBytes, 0644)
		assert.NoError(t, err)

		readBytes, err := os.ReadFile(keyPath)
		assert.NoError(t, err)
		restoredKey, err := rsa.PublicKeyFromBytes(readBytes)
		assert.NoError(t, err)
		assert.True(t, restoredKey.Key.Equal(pubKey.Key))
	})
}
