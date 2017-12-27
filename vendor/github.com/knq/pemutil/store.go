package pemutil

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
)

// Store is a store containing crypto primitives.
//
// A store can contain any of the following crypto primitives:
//     []byte 								-- raw key
//     *rsa.PrivateKey, *ecdsa.PrivateKey   -- rsa / ecdsa private key
//     *rsa.PublicKey, *ecdsa.PublicKey     -- rsa / ecdsa public key
//     *x509.Certificate                    -- x509 certificate
type Store map[BlockType]interface{}

// encOrder is the standard encode order for a Store.
var encOrder = []BlockType{
	PrivateKey,
	RSAPrivateKey,
	ECPrivateKey,
	PublicKey,
	Certificate,
}

// Bytes returns all crypto primitives in the store as a single byte slice
// containing the PEM-encoded versions of the crypto primitives.
func (s Store) Bytes() ([]byte, error) {
	if len(s) == 0 {
		return nil, errors.New("store is empty")
	}

	// encode
	var res bytes.Buffer
	for _, k := range encOrder {
		if p, ok := s[k]; ok {
			buf, err := EncodePrimitive(p)
			if err != nil {
				return nil, err
			}

			_, err = res.Write(buf)
			if err != nil {
				return nil, err
			}
		}
	}

	return res.Bytes(), nil
}

// AddPublicKeys adds the public keys for a RSAPrivateKey or ECPrivateKey block
// type generating and storing the corresponding *PublicKey block if not
// already present.
//
// Useful when a Store is missing the public key for a private key.
func (s Store) AddPublicKeys() {
	if _, ok := s[PublicKey]; ok {
		return
	}

	for _, typ := range []BlockType{PrivateKey, RSAPrivateKey, ECPrivateKey} {
		if key, ok := s[typ]; ok {
			if v, ok := key.(interface {
				Public() crypto.PublicKey
			}); ok {
				s[PublicKey] = v.Public()
			}
		}
	}
}

// Decode parses and decodes PEM-encoded data from buf, storing any resulting
// crypto primitives encountered into the Store. The decoded PEM BlockType will
// be used as the map key for each primitive.
func (s Store) Decode(buf []byte) error {
	return Decode(s, buf)
}

// DecodeBlock decodes PEM block data, adding any crypto primitive encountered the
// store.
func (s Store) DecodeBlock(block *pem.Block) error {
	switch BlockType(block.Type) {
	case PrivateKey:
		// try pkcs1 and then pkcs8 decoding
		key, err := ParsePKCSPrivateKey(block.Bytes)
		if err == nil {
			return s.add(RSAPrivateKey, key)
		}

		// must be a raw key (ie, use decoded b64 value as key)
		return s.add(PrivateKey, block.Bytes)

	case PublicKey:
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			// use the raw b64 decoded bytes
			key = block.Bytes
		}
		return s.add(PublicKey, key)

	case RSAPrivateKey:
		// try pkcs1 then pkcs8 decoding
		key, err := ParsePKCSPrivateKey(block.Bytes)
		if err != nil {
			return err
		}
		return s.add(RSAPrivateKey, key)

	case ECPrivateKey:
		key, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return err
		}
		return s.add(ECPrivateKey, key)

	case Certificate:
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return err
		}
		return s.add(Certificate, cert)
	}

	return fmt.Errorf("unknown block type %s", block.Type)
}

// add adds a crypto primitive to the store, returning an error if the defined
// block is already present.
func (s Store) add(typ BlockType, v interface{}) error {
	if _, ok := s[typ]; ok {
		return fmt.Errorf("block type %s already present", typ)
	}
	s[typ] = v
	return nil
}

// PublicKey returns the public key contained within the store.
func (s Store) PublicKey() (crypto.PublicKey, bool) {
	v, ok := s[PublicKey]
	if !ok {
		return nil, false
	}
	z, ok := v.(crypto.PublicKey)
	return z, ok
}

// PrivateKey returns the private key contained within the store.
func (s Store) PrivateKey() (crypto.PrivateKey, bool) {
	for _, typ := range []BlockType{PrivateKey, RSAPrivateKey, ECPrivateKey} {
		v, ok := s[typ]
		if ok {
			z, ok := v.(crypto.PrivateKey)
			return z, ok
		}
	}

	return nil, false
}

// RSAPublicKey returns the RSA public key contained within the store.
func (s Store) RSAPublicKey() (*rsa.PublicKey, bool) {
	v, ok := s[PublicKey]
	if !ok {
		return nil, false
	}
	z, ok := v.(*rsa.PublicKey)
	return z, ok
}

// RSAPrivateKey returns the RSA private key contained within the store.
func (s Store) RSAPrivateKey() (*rsa.PrivateKey, bool) {
	v, ok := s[RSAPrivateKey]
	if !ok {
		return nil, false
	}
	z, ok := v.(*rsa.PrivateKey)
	return z, ok
}

// ECPublicKey returns the ECDSA public key contained within the store.
func (s Store) ECPublicKey() (*ecdsa.PublicKey, bool) {
	v, ok := s[PublicKey]
	if !ok {
		return nil, false
	}
	z, ok := v.(*ecdsa.PublicKey)
	return z, ok
}

// ECPrivateKey returns the ECDSA private key contained within the store.
func (s Store) ECPrivateKey() (*ecdsa.PrivateKey, bool) {
	v, ok := s[ECPrivateKey]
	if !ok {
		return nil, false
	}
	z, ok := v.(*ecdsa.PrivateKey)
	return z, ok
}

// Certificate returns the X509 certificate contained within the store.
func (s Store) Certificate() (*x509.Certificate, bool) {
	v, ok := s[Certificate]
	if !ok {
		return nil, false
	}
	z, ok := v.(*x509.Certificate)
	return z, ok
}

// LoadFile loads crypto primitives from PEM encoded data stored in filename.
func (s Store) LoadFile(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return Decode(s, buf)
}

// LoadFile creates a store and loads any crypto primitives in the PEM encoded
// data stored in filename.
//
// Note: calls AddPublicKeys() after successfully loading a file. If that
// behavior is not desired, please manually create the store and call Decode,
// or DecodeBlock.
func LoadFile(filename string) (Store, error) {
	s := make(Store)
	if err := s.LoadFile(filename); err != nil {
		return nil, err
	}
	s.AddPublicKeys()
	return s, nil
}

// WriteFile writes the crypto primitives in the store to filename with mode
// 0600.
func (s Store) WriteFile(filename string) error {
	buf, err := s.Bytes()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, buf, 0600)
}
