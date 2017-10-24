package pemutil

import (
	"bytes"
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
func (s Store) AddPublicKeys() error {
	// generate rsa public key
	if key, ok := s[RSAPrivateKey]; ok {
		rsaPrivKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return errors.New("block type RSAPrivateKey does not contain *rsa.PrivateKey")
		}
		if _, ok = s[PublicKey]; !ok {
			return s.add(PublicKey, rsaPrivKey.Public())
		}
	}

	// generate ecdsa public key
	if key, ok := s[ECPrivateKey]; ok {
		ecdsaPrivKey, ok := key.(*ecdsa.PrivateKey)
		if !ok {
			return errors.New("block type ECPrivateKey does not contain *ecdsa.PrivateKey")
		}
		if _, ok = s[PublicKey]; !ok {
			return s.add(PublicKey, ecdsaPrivKey.Public())
		}
	}

	return nil
}

// Decode parses and decodes PEM-encoded data from buf, storing any resulting
// crypto primitives encountered into the Store. The decoded PEM BlockType will
// be used as the map key for each primitive.
func (s Store) Decode(buf []byte) error {
	var err error
	var block *pem.Block

	// loop over pem encoded data
	for len(buf) > 0 {
		block, buf = pem.Decode(buf)
		if block == nil {
			return errors.New("invalid PEM data")
		}

		err = s.addBlock(block)
		if err != nil {
			return err
		}
	}

	if len(s) == 0 {
		return errors.New("could not decode any PEM blocks")
	}

	return nil
}

// addBlock decodes PEM block data, adding any crypto primitive to the store.
func (s Store) addBlock(block *pem.Block) error {
	switch BlockType(block.Type) {
	case PrivateKey:
		// try pkcs1 and then pkcs8 decoding
		key, err := parsePKCSPrivateKey(block.Bytes)
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
		key, err := parsePKCSPrivateKey(block.Bytes)
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

// LoadFile loads crypto primitives from PEM encoded data stored in filename.
func (s Store) LoadFile(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return s.Decode(buf)
}
