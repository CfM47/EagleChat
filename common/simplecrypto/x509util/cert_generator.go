package x509util

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"time"

	simple_rsa "eaglechat/common/simplecrypto/rsa"
)

var (
	ErrCACertificateInvalid = errors.New("x509util: CA certificate provided is not a valid CA")
	ErrCAPrivateKeyInvalid  = errors.New("x509util: failed to parse CA private key")
	ErrCertCreation         = errors.New("x509util: failed to create certificate")
	ErrCertEncoding         = errors.New("x509util: failed to encode certificate to PEM")
)

// CreateAndSignCertificate generates a new X.509 certificate for a subject,
// signed by the provided CA.
// caCertPEM and caKeyPEM are the PEM-encoded CA certificate and private key.
// subjectPubKey is the public key of the entity for whom the certificate is being issued.
// commonName is the Common Name for the subject of the new certificate.
func CreateAndSignCertificate(
	caCertPEM []byte,
	caKeyPEM []byte,
	subjectPubKey *simple_rsa.PublicKey,
	commonName string,
) ([]byte, error) {
	// 1. Parse CA Certificate
	pemBlock, _ := pem.Decode(caCertPEM)
	if pemBlock == nil || pemBlock.Type != "CERTIFICATE" {
		return nil, ErrCACertificateInvalid
	}
	caCert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return nil, errors.Join(ErrCACertificateInvalid, err)
	}
	if !caCert.IsCA {
		return nil, ErrCACertificateInvalid
	}

	// 2. Parse CA Private Key
	caPrivKey, err := simple_rsa.PrivateKeyFromBytes(caKeyPEM)
	if err != nil {
		return nil, errors.Join(ErrCAPrivateKeyInvalid, err)
	}

	// 3. Prepare Subject Certificate Template
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to generate serial number: %v", ErrCertCreation, err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"EagleChat"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour), // Valid for 1 year

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign, // CertSign for CA
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}
	// The BasicConstraintsValid implies that this is NOT a CA certificate itself,
	// unless IsCA is set to true. Here, we want it to be a regular server/client cert.
	// If the subject is also a CA, IsCA would be true here.

	// 4. Create the Certificate
	// x509.CreateCertificate expects the CA's private key to implement crypto.Signer,
	// which simple_rsa.PrivateKey's underlying *rsa.PrivateKey does.
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, caCert, subjectPubKey.Key, caPrivKey.Key)
	if err != nil {
		return nil, errors.Join(ErrCertCreation, err)
	}

	// 5. Encode to PEM
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if certPEM == nil {
		return nil, ErrCertEncoding
	}

	return certPEM, nil
}
