package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

const KeySize = 4096

// PublicKey is a wrapper around the standard rsa.PublicKey for API isolation.
type PublicKey struct {
	Key *rsa.PublicKey
}

// PrivateKey is a wrapper around the standard rsa.PrivateKey for API isolation.
type PrivateKey struct {
	Key *rsa.PrivateKey
}

// GenerateKeyPair creates a new RSA key pair.
func GenerateKeyPair() (*PrivateKey, *PublicKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, KeySize)
	if err != nil {
		return nil, nil, err
	}
	return &PrivateKey{Key: priv}, &PublicKey{Key: &priv.PublicKey}, nil
}

// --- Key Serialization ---

func (priv *PrivateKey) ToBytes() []byte {
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv.Key),
	}
	return pem.EncodeToMemory(pemBlock)
}

func (pub *PublicKey) ToBytes() ([]byte, error) {
	pubBytes, err := x509.MarshalPKIXPublicKey(pub.Key)
	if err != nil {
		return nil, err
	}
	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}
	return pem.EncodeToMemory(pemBlock), nil
}

// --- Key Deserialization ---

func PrivateKeyFromBytes(pemBytes []byte) (*PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{Key: key}, nil
}

func PublicKeyFromBytes(pemBytes []byte) (*PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("key in PEM block is not an RSA public key")
	}
	return &PublicKey{Key: rsaPub}, nil
}

// --- Cryptographic Operations ---

func Encrypt(msg []byte, pub *PublicKey) ([]byte, error) {
	hash := sha256.New()
	return rsa.EncryptOAEP(hash, rand.Reader, pub.Key, msg, nil)
}

func Decrypt(ciphertext []byte, priv *PrivateKey) ([]byte, error) {
	hash := sha256.New()
	return rsa.DecryptOAEP(hash, rand.Reader, priv.Key, ciphertext, nil)
}

func Sign(msg []byte, priv *PrivateKey) ([]byte, error) {
	hashed := sha256.Sum256(msg)
	return rsa.SignPSS(rand.Reader, priv.Key, crypto.SHA256, hashed[:], nil)
}

func Verify(msg []byte, signature []byte, pub *PublicKey) error {
	hashed := sha256.Sum256(msg)
	return rsa.VerifyPSS(pub.Key, crypto.SHA256, hashed[:], signature, nil)
}
