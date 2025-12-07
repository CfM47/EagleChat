package x509util

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem" // New import
	"errors"
	"os"

	simple_rsa "eaglechat/common/simplecrypto/rsa"
)

var (
	ErrCertificateParsing = errors.New("x509util: failed to parse certificate")
	ErrCACertificate      = errors.New("x509util: failed to load or parse CA certificate")
	ErrCertificateNotRSA  = errors.New("x509util: certificate does not contain an RSA public key")
	ErrVerificationFailed = errors.New("x509util: certificate verification failed")
)

// Verifier is responsible for verifying peer certificates against a trusted CA.
type Verifier struct {
	caCertPool *x509.CertPool
}

// NewVerifier creates a new verifier.
// The caCertPath should be the path to the trusted root CA's certificate file in PEM format.
func NewVerifier(caCertPath string) (*Verifier, error) {
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, ErrCACertificate
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		return nil, ErrCACertificate
	}

	return &Verifier{
		caCertPool: certPool,
	}, nil
}

// verifyCertificate checks if a given certificate is valid and has been signed by the trusted CA.
// The certificate is expected to be in PEM format. It returns the parsed certificate on success.
func (v *Verifier) verifyCertificate(certPEMBytes []byte) (*x509.Certificate, error) {
	// First, decode the PEM block to get the raw DER bytes
	pemBlock, _ := pem.Decode(certPEMBytes)
	if pemBlock == nil || pemBlock.Type != "CERTIFICATE" {
		return nil, errors.Join(ErrCertificateParsing, errors.New("invalid PEM block or not a CERTIFICATE type"))
	}
	certDERBytes := pemBlock.Bytes

	// Now parse the DER-encoded certificate
	cert, err := x509.ParseCertificate(certDERBytes)
	if err != nil {
		return nil, errors.Join(ErrCertificateParsing, err)
	}

	// Set up verification options
	opts := x509.VerifyOptions{
		Roots: v.caCertPool,
	}

	if _, err := cert.Verify(opts); err != nil {
		return nil, errors.Join(ErrVerificationFailed, err)
	}

	return cert, nil
}

// VerifyAndExtractPublicKey verifies the certificate and, on success, extracts the public key,
// returning it as our abstracted simple_rsa.PublicKey type.
func (v *Verifier) VerifyAndExtractPublicKey(certBytes []byte) (*simple_rsa.PublicKey, error) {
	cert, err := v.verifyCertificate(certBytes)
	if err != nil {
		return nil, err
	}

	// Type assert the public key to the standard library's RSA public key type.
	// This is the only place in our code base that should do this.
	stdPubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, ErrCertificateNotRSA
	}

	// Wrap the standard key in our own abstracted type.
	return &simple_rsa.PublicKey{Key: stdPubKey}, nil
}
