package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"io"

	"github.com/knq/pemutil"
)

const (
	// RSAMinimumBitLen is the minimum accepted RSA key length.
	RSAMinimumBitLen = 2048
)

// RSASignerVerifier provides a standardized interface to low level RSA signing
// implementation.
//
// This is used internally to provide a common interface to the RSA Sign/Verify
// implementations for PKCS1v15 and PSS.
type RSASignerVerifier interface {
	// Sign signs data in buf using rand, priv and hash.
	Sign(rand io.Reader, priv *rsa.PrivateKey, hash crypto.Hash, buf []byte) ([]byte, error)

	// Verify verifies the signature sig against using pub, hash, and the
	// hashed data.
	Verify(pub *rsa.PublicKey, hash crypto.Hash, hashed []byte, sig []byte) error
}

// rsaMethod provides a wrapper for rsa signing methods.
type rsaMethod struct {
	signFunc   func(io.Reader, *rsa.PrivateKey, crypto.Hash, []byte) ([]byte, error)
	verifyFunc func(*rsa.PublicKey, crypto.Hash, []byte, []byte) error
}

// Sign signs the data in buf using rand, priv and hash.
func (r rsaMethod) Sign(rand io.Reader, priv *rsa.PrivateKey, hash crypto.Hash, buf []byte) ([]byte, error) {
	return r.signFunc(rand, priv, hash, buf)
}

// Verify verifies the signature sig against using pub, hash, and the hashed
// data.
func (r rsaMethod) Verify(pub *rsa.PublicKey, hash crypto.Hash, hashed []byte, sig []byte) error {
	return r.verifyFunc(pub, hash, hashed, sig)
}

// rsaMethodPKCS1v15 provides a RSA method that signs and verifies with
// PKCS1v15.
var rsaMethodPKCS1v15 = rsaMethod{
	signFunc:   rsa.SignPKCS1v15,
	verifyFunc: rsa.VerifyPKCS1v15,
}

// rsaMethodPSS provides a RSA method that signs and verifies with PSS.
var rsaMethodPSS = rsaMethod{
	signFunc: func(rand io.Reader, priv *rsa.PrivateKey, hash crypto.Hash, hashed []byte) ([]byte, error) {
		return rsa.SignPSS(rand, priv, hash, hashed, &rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthAuto,
			Hash:       hash,
		})
	},
	verifyFunc: func(pub *rsa.PublicKey, hash crypto.Hash, hashed []byte, sig []byte) error {
		return rsa.VerifyPSS(pub, hash, hashed, sig, &rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthAuto,
			Hash:       hash,
		})
	},
}

// rsaSigner provides a RSA Signer.
type rsaSigner struct {
	alg    Algorithm
	method RSASignerVerifier
	hash   crypto.Hash
	priv   *rsa.PrivateKey
	pub    *rsa.PublicKey
}

// NewRSASigner creates an RSA Signer for the specified Algorithm and provided
// low level RSA implementation.
func NewRSASigner(alg Algorithm, method RSASignerVerifier) func(pemutil.Store, crypto.Hash) (Signer, error) {
	return func(store pemutil.Store, hash crypto.Hash) (Signer, error) {
		var ok bool
		var privRaw, pubRaw interface{}
		var priv *rsa.PrivateKey
		var pub *rsa.PublicKey

		// check private key
		if privRaw, ok = store[pemutil.RSAPrivateKey]; ok {
			if priv, ok = privRaw.(*rsa.PrivateKey); !ok {
				return nil, errors.New("NewRSASigner: private key must be a *rsa.PrivateKey")
			}

			// check private key length
			if priv.N.BitLen() < RSAMinimumBitLen {
				return nil, fmt.Errorf("NewRSASigner: private key has length %d, must have minimum length of %d", priv.N.BitLen(), RSAMinimumBitLen)
			}
		}

		// check public key
		if pubRaw, ok = store[pemutil.PublicKey]; ok {
			if pub, ok = pubRaw.(*rsa.PublicKey); !ok {
				return nil, errors.New("NewRSASigner: public key must be a *rsa.PublicKey")
			}

			// check public key length
			if pub.N.BitLen() < RSAMinimumBitLen {
				return nil, fmt.Errorf("NewRSASigner: public key has length %d, must have minimum length of %d", pub.N.BitLen(), RSAMinimumBitLen)
			}
		}

		// check that either a private or public key has been provided
		if priv == nil && pub == nil {
			return nil, errors.New("NewRSASigner: either a private key or a public key must be provided")
		}

		return &rsaSigner{
			alg:    alg,
			method: method,
			hash:   hash,
			priv:   priv,
			pub:    pub,
		}, nil
	}
}

// SignBytes creates a signature for buf.
func (rs *rsaSigner) SignBytes(buf []byte) ([]byte, error) {
	var err error

	// check rs.priv
	if rs.priv == nil {
		return nil, errors.New("rsaSigner.Sign: priv cannot be nil")
	}

	// hash
	h := rs.hash.New()
	_, err = h.Write(buf)
	if err != nil {
		return nil, err
	}

	// sign
	return rs.method.Sign(rand.Reader, rs.priv, rs.hash, h.Sum(nil))
}

// Sign creates a signature for buf, returning it as a URL-safe base64 encoded
// byte slice.
func (rs *rsaSigner) Sign(buf []byte) ([]byte, error) {
	sig, err := rs.SignBytes(buf)
	if err != nil {
		return nil, err
	}

	// encode
	enc := make([]byte, b64.EncodedLen(len(sig)))
	b64.Encode(enc, sig)

	return enc, nil
}

// VerifyBytes creates a signature for buf, comparing it against the raw sig.
// If the sig is invalid, then ErrInvalidSignature is returned.
func (rs *rsaSigner) VerifyBytes(buf, sig []byte) error {
	var err error

	// check rs.pub
	if rs.pub == nil {
		return errors.New("rsaSigner.VerifyBytes: pub cannot be nil")
	}

	// hash
	h := rs.hash.New()
	_, err = h.Write(buf)
	if err != nil {
		return err
	}

	// verify
	err = rs.method.Verify(rs.pub, rs.hash, h.Sum(nil), sig)
	if err != nil {
		return ErrInvalidSignature
	}

	return nil
}

// Verify creates a signature for buf, comparing it against the URL-safe base64
// encoded sig and returning the decoded signature. If the sig is invalid, then
// ErrInvalidSignature will be returned.
func (rs *rsaSigner) Verify(buf, sig []byte) ([]byte, error) {
	var err error

	// decode
	dec, err := b64.DecodeString(string(sig))
	if err != nil {
		return nil, err
	}

	// verify
	err = rs.VerifyBytes(buf, dec)
	if err != nil {
		return nil, err
	}

	return dec, nil
}

// Encode serializes the JSON marshalable obj data as a JWT.
func (rs *rsaSigner) Encode(obj interface{}) ([]byte, error) {
	return rs.alg.Encode(rs, obj)
}

// Decode decodes a serialized token, verifying the signature, storing the
// decoded data from the token in obj.
func (rs *rsaSigner) Decode(buf []byte, obj interface{}) error {
	return rs.alg.Decode(rs, buf, obj)
}
