package simplecrypto_test

import (
	"testing"

	"eaglechat/common/simplecrypto"
	"eaglechat/common/simplecrypto/rsa"

	"github.com/stretchr/testify/assert"
)

func TestSignedEnvelope(t *testing.T) {
	t.Run("Seal and Open Success", func(t *testing.T) {
		privAlice, pubAlice, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)

		originalMessage := []byte("this is a top secret message")
		envelope, err := simplecrypto.SignEnvelope(originalMessage, privAlice)
		assert.NoError(t, err)
		assert.NotNil(t, envelope)

		msg, err := simplecrypto.OpenSignedEnvelope(envelope, pubAlice)
		assert.NoError(t, err)
		assert.Equal(t, originalMessage, msg)
	})

	t.Run("Open Fails With Wrong Key", func(t *testing.T) {
		privAlice, _, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)
		_, pubBob, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)

		originalMessage := []byte("secret message")
		envelope, err := simplecrypto.SignEnvelope(originalMessage, privAlice)
		assert.NoError(t, err)

		_, err = simplecrypto.OpenSignedEnvelope(envelope, pubBob)
		assert.Error(t, err)
	})

	t.Run("Open Fails With Tampered Signature", func(t *testing.T) {
		privAlice, pubAlice, err := rsa.GenerateKeyPair()
		assert.NoError(t, err)

		originalMessage := []byte("secret message")
		envelope, err := simplecrypto.SignEnvelope(originalMessage, privAlice)
		assert.NoError(t, err)

		envelope.Signature[0] ^= 0xff

		_, err = simplecrypto.OpenSignedEnvelope(envelope, pubAlice)
		assert.Error(t, err)
	})
}
