package simplecrypto

import (
	"crypto/sha256"
	"errors"

	"eaglechat/apps/client/internal/utils/simplecrypto/aes"
	"eaglechat/apps/client/internal/utils/simplecrypto/rsa"
)

var (
	ErrInvalidSignature  = errors.New("crypto: signature verification failed")
	ErrDecryptKey        = errors.New("crypto: failed to decrypt AES key")
	ErrDecryptCiphertext = errors.New("crypto: failed to decrypt message ciphertext")
)

// SecureEnvelope is the data structure for a fully encrypted and authenticated message.
// This is the packet that should be sent over the network.
type SecureEnvelope struct {
	// The AES key, encrypted with the recipient's RSA public key.
	EncryptedAESKey []byte
	// The message content, encrypted with the AES key.
	Ciphertext []byte
	// A signature of the other two fields to ensure authenticity and integrity.
	Signature []byte
}

// Seal encrypts and signs a message to create a secure envelope.
func Seal(message []byte, senderPrivKey *rsa.PrivateKey, recipientPubKey *rsa.PublicKey) (*SecureEnvelope, error) {
	// 1. Generate a new AES key for this message.
	 aesKey, err := aes.GenerateKey()
	if err != nil {
		return nil, err
	}

	// 2. Encrypt the message with the AES key.
	ciphertext, err := aes.Encrypt(message, aesKey)
	if err != nil {
		return nil, err
	}

	// 3. Encrypt the AES key with the recipient's RSA public key.
	encryptedAESKey, err := rsa.Encrypt(aesKey, recipientPubKey)
	if err != nil {
		return nil, err
	}

	// 4. Sign the hashes of the encrypted parts.
	hasher := sha256.New()
	hasher.Write(encryptedAESKey)
	hasher.Write(ciphertext)
	signature, err := rsa.Sign(hasher.Sum(nil), senderPrivKey)
	if err != nil {
		return nil, err
	}

	return &SecureEnvelope{
		EncryptedAESKey: encryptedAESKey,
		Ciphertext:      ciphertext,
		Signature:       signature,
	}, nil
}

// Open decrypts and verifies a secure envelope to retrieve the original message.
func Open(envelope *SecureEnvelope, recipientPrivKey *rsa.PrivateKey, senderPubKey *rsa.PublicKey) ([]byte, error) {
	// 1. Verify the signature first.
	hasher := sha256.New()
	hasher.Write(envelope.EncryptedAESKey)
	hasher.Write(envelope.Ciphertext)
	if err := rsa.Verify(hasher.Sum(nil), envelope.Signature, senderPubKey); err != nil {
		return nil, ErrInvalidSignature
	}

	// 2. Decrypt the AES key with our private key.
	aesKey, err := rsa.Decrypt(envelope.EncryptedAESKey, recipientPrivKey)
	if err != nil {
		return nil, ErrDecryptKey
	}

	// 3. Decrypt the message with the revealed AES key.
	plaintext, err := aes.Decrypt(envelope.Ciphertext, aesKey)
	if err != nil {
		return nil, ErrDecryptCiphertext
	}

	return plaintext, nil
}
