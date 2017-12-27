package pemutil

import (
	"crypto/x509"
)

// BlockType is a PEM block type.
type BlockType string

// String satisfies the string interface for a block type.
func (bt BlockType) String() string {
	return string(bt)
}

const (
	// PrivateKey is the "PRIVATE KEY" block type.
	PrivateKey BlockType = "PRIVATE KEY"

	// RSAPrivateKey is the "RSA PRIVATE KEY" block type.
	RSAPrivateKey BlockType = "RSA PRIVATE KEY"

	// ECPrivateKey is the "EC PRIVATE KEY" block type.
	ECPrivateKey BlockType = "EC PRIVATE KEY"

	// PublicKey is the "PUBLIC KEY" block type.
	PublicKey BlockType = "PUBLIC KEY"

	// Certificate is the "CERTIFICATE" block type.
	Certificate BlockType = "CERTIFICATE"
)

// ParsePKCSPrivateKey attempts to decode a RSA private key first using PKCS1
// encoding, and then PKCS8 encoding.
func ParsePKCSPrivateKey(buf []byte) (interface{}, error) {
	// attempt PKCS1 parsing
	key, err := x509.ParsePKCS1PrivateKey(buf)
	if err == nil {
		return key, nil
	}

	// attempt PKCS8 parsing
	return x509.ParsePKCS8PrivateKey(buf)
}
