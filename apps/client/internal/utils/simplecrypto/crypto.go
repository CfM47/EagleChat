package simplecrypto

import (
	"eaglechat/apps/client/internal/utils/simplecrypto/aes"
	"eaglechat/apps/client/internal/utils/simplecrypto/rsa"
	"encoding/json"
	"errors"
)

var (
	ErrInvalidSignature  = errors.New("crypto: signature verification failed")
	ErrDecryptKey        = errors.New("crypto: failed to decrypt AES key")
	ErrDecryptCiphertext = errors.New("crypto: failed to decrypt message ciphertext")
)

// InnerEnvelope contains the actual application data and the sender's public key
// needed for signature verification. This struct is what gets encrypted.
type InnerEnvelope struct {
	SenderPubKeyBytes []byte `json:"sender_pub_key"`
	Message           []byte `json:"message"`
}

// SecureEnvelope is the data structure for a fully encrypted and authenticated message.
// This is the packet that should be sent over the network.
type SecureEnvelope struct {
	// The AES key, encrypted with the recipient's RSA public key.
	EncryptedAESKey []byte `json:"encrypted_aes_key"`
	// The encrypted InnerEnvelope.
	Ciphertext []byte `json:"ciphertext"`
	// A signature of the Ciphertext to ensure authenticity and integrity.
	Signature []byte `json:"signature"`
}

// Seal encrypts and signs a message to create a secure envelope.
func Seal(message []byte, senderPrivKey *rsa.PrivateKey, recipientPubKey *rsa.PublicKey) (*SecureEnvelope, error) {
	// 1. Marshal the sender's public key.
	senderPubKeyBytes, err := senderPrivKey.PublicKey().ToBytes()
	if err != nil {
		return nil, err
	}

	// 2. Create and marshal the inner envelope.
	innerEnvelope := &InnerEnvelope{
		SenderPubKeyBytes: senderPubKeyBytes,
		Message:           message,
	}
	innerEnvelopeBytes, err := json.Marshal(innerEnvelope)
	if err != nil {
		return nil, err
	}

	// 3. Generate a new AES key for this message.
	aesKey, err := aes.GenerateKey()
	if err != nil {
		return nil, err
	}

	// 4. Encrypt the inner envelope with the AES key.
	ciphertext, err := aes.Encrypt(innerEnvelopeBytes, aesKey)
	if err != nil {
		return nil, err
	}

	// 5. Encrypt the AES key with the recipient's RSA public key.
	encryptedAESKey, err := rsa.Encrypt(aesKey, recipientPubKey)
	if err != nil {
		return nil, err
	}

	// 6. Sign the ciphertext directly (the underlying Sign func will hash it).
	signature, err := rsa.Sign(ciphertext, senderPrivKey)
	if err != nil {
		return nil, err
	}

	return &SecureEnvelope{
		EncryptedAESKey: encryptedAESKey,
		Ciphertext:      ciphertext,
		Signature:       signature,
	}, nil
}

// Open decrypts and verifies a secure envelope to retrieve the original message
// and the sender's public key.
func Open(envelope *SecureEnvelope, recipientPrivKey *rsa.PrivateKey) ([]byte, *rsa.PublicKey, error) {
	// 1. Decrypt the AES key with our private key.
	aesKey, err := rsa.Decrypt(envelope.EncryptedAESKey, recipientPrivKey)
	if err != nil {
		return nil, nil, ErrDecryptKey
	}

	// 2. Decrypt the ciphertext with the revealed AES key.
	innerEnvelopeBytes, err := aes.Decrypt(envelope.Ciphertext, aesKey)
	if err != nil {
		return nil, nil, ErrDecryptCiphertext
	}

	// 3. Unmarshal the inner envelope to get the sender's public key.
	var innerEnvelope InnerEnvelope
	if err := json.Unmarshal(innerEnvelopeBytes, &innerEnvelope); err != nil {
		return nil, nil, err
	}

	// 4. Reconstruct the sender's public key.
	senderPubKey, err := rsa.PublicKeyFromBytes(innerEnvelope.SenderPubKeyBytes)
	if err != nil {
		return nil, nil, err
	}

	// 5. Verify the signature against the ciphertext using the sender's public key.
	if err := rsa.Verify(envelope.Ciphertext, envelope.Signature, senderPubKey); err != nil {
		return nil, nil, ErrInvalidSignature
	}

	// 6. If all checks pass, return the original message and the verified public key.
	return innerEnvelope.Message, senderPubKey, nil
}
