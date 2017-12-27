package jwt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

// EccSigner provides an Elliptic Curve Signer.
type EccSigner struct {
	alg   Algorithm
	curve elliptic.Curve
	hash  crypto.Hash
	priv  *ecdsa.PrivateKey
	pub   *ecdsa.PublicKey

	keyLen int
}

// NewEllipticSigner creates an Elliptic Curve Signer for the specified curve.
func NewEllipticSigner(alg Algorithm, curve elliptic.Curve) func(Store, crypto.Hash) (Signer, error) {
	curveBitSize := curve.Params().BitSize

	// precompute curve key len
	keyLen := curveBitSize / 8
	if curveBitSize%8 > 0 {
		keyLen++
	}

	return func(store Store, hash crypto.Hash) (Signer, error) {
		var ok bool
		var privRaw, pubRaw interface{}
		var priv *ecdsa.PrivateKey
		var pub *ecdsa.PublicKey

		// check private key
		if privRaw, ok = store.PrivateKey(); ok {
			if priv, ok = privRaw.(*ecdsa.PrivateKey); !ok {
				return nil, errors.New("NewEllipticSigner: private key must be a *ecdsa.PrivateKey")
			}

			// check curve type matches private key curve type
			if curveBitSize != priv.Curve.Params().BitSize {
				return nil, fmt.Errorf("NewEllipticSigner: private key have bit size %d", curve.Params().BitSize)
			}
		}

		// check public key
		if pubRaw, ok = store.PublicKey(); ok {
			if pub, ok = pubRaw.(*ecdsa.PublicKey); !ok {
				return nil, errors.New("NewEllipticSigner: public key must be a *ecdsa.PublicKey")
			}
		}

		// check that either a private or public key has been provided
		if priv == nil && pub == nil {
			return nil, errors.New("NewEllipticSigner: either a private key or a public key must be provided")
		}

		return &EccSigner{
			alg:    alg,
			curve:  curve,
			hash:   hash,
			priv:   priv,
			pub:    pub,
			keyLen: keyLen,
		}, nil
	}
}

// Mksig creates a byte slice of length 2*keyLen, copying the bytes from r and
// s into the slice, left padding r and s to keyLen.
func (es *EccSigner) Mksig(r, s *big.Int) ([]byte, error) {
	var n int

	buf := make([]byte, 2*es.keyLen)

	// copy r into buf
	rb := r.Bytes()
	n = copy(buf[es.keyLen-len(rb):], rb)
	if n != len(rb) {
		return nil, fmt.Errorf("EccSigner.Mksig: could not copy r into sig, copied: %d", n)
	}

	// copy s into buf
	sb := s.Bytes()
	n = copy(buf[es.keyLen+(es.keyLen-(len(sb))):], sb)
	if n != len(sb) {
		return nil, fmt.Errorf("EccSigner.Mksig: could not copy s into sig, copied: %d", n)
	}

	return buf, nil
}

// SignBytes creates a signature for buf.
func (es *EccSigner) SignBytes(buf []byte) ([]byte, error) {
	var err error

	// check es.priv
	if es.priv == nil {
		return nil, errors.New("EccSigner.SignBytes: priv cannot be nil")
	}

	// hash
	h := es.hash.New()
	_, err = h.Write(buf)
	if err != nil {
		return nil, err
	}

	// sign
	r, s, err := ecdsa.Sign(rand.Reader, es.priv, h.Sum(nil))
	if err != nil {
		return nil, err
	}

	// make sig
	return es.Mksig(r, s)
}

// Sign creates a signature for buf, returning it as a URL-safe base64 encoded
// byte slice.
func (es *EccSigner) Sign(buf []byte) ([]byte, error) {
	sig, err := es.SignBytes(buf)
	if err != nil {
		return nil, err
	}

	enc := make([]byte, b64.EncodedLen(len(sig)))
	b64.Encode(enc, sig)

	return enc, nil
}

// VerifyBytes creates a signature for buf, comparing it against the raw sig.
// If the sig is invalid, then ErrInvalidSignature is returned.
func (es *EccSigner) VerifyBytes(buf, sig []byte) error {
	var err error

	// check es.pub
	if es.pub == nil {
		return errors.New("EccSigner.VerifyBytes: pub cannot be nil")
	}

	// hash
	h := es.hash.New()
	_, err = h.Write(buf)
	if err != nil {
		return err
	}

	// check decoded length
	if len(sig) != 2*es.keyLen {
		return ErrInvalidSignature
	}

	r := big.NewInt(0).SetBytes(sig[:es.keyLen])
	s := big.NewInt(0).SetBytes(sig[es.keyLen:])

	// verify
	if !ecdsa.Verify(es.pub, h.Sum(nil), r, s) {
		return ErrInvalidSignature
	}

	return nil
}

// Verify creates a signature for buf, comparing it against the URL-safe base64
// encoded sig and returning the decoded signature. If the sig is invalid, then
// ErrInvalidSignature will be returned.
func (es *EccSigner) Verify(buf, sig []byte) ([]byte, error) {
	var err error

	// decode
	dec, err := b64.DecodeString(string(sig))
	if err != nil {
		return nil, err
	}

	// verify
	err = es.VerifyBytes(buf, dec)
	if err != nil {
		return nil, err
	}

	return dec, nil
}

// Encode serializes the JSON marshalable obj data as a JWT.
func (es *EccSigner) Encode(obj interface{}) ([]byte, error) {
	return es.alg.Encode(es, obj)
}

// Decode decodes a serialized token, verifying the signature, storing the
// decoded data from the token in obj.
func (es *EccSigner) Decode(buf []byte, obj interface{}) error {
	return es.alg.Decode(es, buf, obj)
}
