// Package pemutil provides a simple, high-level API to load, parse, and decode
// standard crypto primitives (ie, rsa.PrivateKey, ecdsa.PrivateKey, etc) from
// PEM-encoded data.
//
// The pemutil package commonly used similar to the following:
//
//		store, err := pemutil.LoadFile("/path/to/file")
// 		if err != nil { /* ... */ }
//
//		if rsaPrivKey, ok := store.RSAPrivateKey(); !ok {
//			// PEM does not contain an RSA private key
//		}
//
package pemutil

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

// Decode parses and decodes PEM-encoded data from buf, storing any resulting
// crypto primitives encountered into the Store. The decoded PEM BlockType will
// be used as the map key for each primitive.
func Decode(s Store, buf []byte) error {
	var err error
	var block *pem.Block

	// loop over pem encoded data
	for len(buf) > 0 {
		block, buf = pem.Decode(buf)
		if block == nil {
			return errors.New("invalid PEM data")
		}

		err = s.DecodeBlock(block)
		if err != nil {
			return err
		}
	}

	if len(s) == 0 {
		return errors.New("could not decode any PEM blocks")
	}

	return nil
}

// DecodeBytes decodes the supplied buf into a store.
func DecodeBytes(buf []byte) (Store, error) {
	s := Store{}
	err := Decode(s, buf)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// EncodePrimitive encodes the crypto primitive p into PEM-encoded data.
func EncodePrimitive(p interface{}) ([]byte, error) {
	var err error
	var typ BlockType
	var buf []byte

	switch v := p.(type) {
	case []byte:
		typ, buf = PrivateKey, v

	case *rsa.PrivateKey:
		typ, buf = RSAPrivateKey, x509.MarshalPKCS1PrivateKey(v)

	case *ecdsa.PrivateKey:
		typ = ECPrivateKey
		buf, err = x509.MarshalECPrivateKey(v)
		if err != nil {
			return nil, err
		}

	case *rsa.PublicKey, *ecdsa.PublicKey:
		typ = PublicKey
		buf, err = x509.MarshalPKIXPublicKey(v)
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("unsupported crypto primitive")
	}

	// encode
	return pem.EncodeToMemory(&pem.Block{
		Type:  typ.String(),
		Bytes: buf,
	}), nil
}

// GenerateSymmetricKeySet generates a private key crypto primitive, returning
// it as a Store.
func GenerateSymmetricKeySet(keyLen int) (Store, error) {
	// generate random bytes
	buf := make([]byte, keyLen)
	c, err := rand.Read(buf)
	if err != nil {
		return nil, err
	} else if c != keyLen {
		return nil, fmt.Errorf("could not generate %d random key bits", keyLen)
	}

	return Store{
		PrivateKey: buf,
	}, nil
}

// GenerateRSAKeySet generates a RSA private and public key crypto primitives,
// returning them as a Store.
func GenerateRSAKeySet(bitLen int) (Store, error) {
	key, err := rsa.GenerateKey(rand.Reader, bitLen)
	if err != nil {
		return nil, err
	}

	return Store{
		RSAPrivateKey: key,
		PublicKey:     key.Public(),
	}, nil
}

// GenerateECKeySet generates a EC private and public key crypto primitives,
// returning them as a Store.
func GenerateECKeySet(curve elliptic.Curve) (Store, error) {
	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}

	return Store{
		ECPrivateKey: key,
		PublicKey:    key.Public(),
	}, nil
}
