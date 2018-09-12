package jwt

//go:generate stringer -type Algorithm -output alg_string.go .

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
)

// Algorithm is the type for signing algorithms implemented in this package.
type Algorithm uint

// Signer is the shared interface for an Algorithm's encoding, decoding,
// signing, and verify to handle the crypto primitives and lower-level API
// calls.
type Signer interface {
	// SignBytes creates a signature for buf.
	SignBytes(buf []byte) ([]byte, error)

	// Sign creates a signature for buf, returning it as a URL-safe base64
	// encoded byte slice.
	Sign(buf []byte) ([]byte, error)

	// VerifyBytes creates a signature for buf, comparing it against the raw
	// sig. If the sig is invalid, then ErrInvalidSignature is returned.
	VerifyBytes(buf, dec []byte) error

	// Verify creates a signature for buf, comparing it against the URL-safe
	// base64 encoded sig and returning the decoded signature. If the sig is
	// invalid, then ErrInvalidSignature will be returned.
	Verify(buf, sig []byte) ([]byte, error)

	// Encode encodes obj as a serialized JWT.
	Encode(obj interface{}) ([]byte, error)

	// Decode decodes a serialized JWT in buf into obj, and verifying the JWT
	// signature in the process.
	Decode(buf []byte, obj interface{}) error
}

const (
	// NONE provides a JWT signing method for NONE.
	//
	// NOTE: This is not implemented for security reasons.
	NONE Algorithm = iota

	// HS256 provides a JWT signing method for HMAC using SHA-256.
	//
	// See http://tools.ietf.org/html/rfc7518#section-3.2
	HS256

	// HS384 provides a JWT signing method for HMAC using SHA-384.
	//
	// See http://tools.ietf.org/html/rfc7518#section-3.2
	HS384

	// HS512 provides a JWT signing method for HMAC using SHA-512.
	//
	// See http://tools.ietf.org/html/rfc7518#section-3.2
	HS512

	// RS256 provides a JWT signing method for RSASSA-PKCS1-V1_5 using SHA-256.
	//
	// See http://tools.ietf.org/html/rfc7518#section-3.3
	RS256

	// RS384 provides a JWT signing method for RSASSA-PKCS1-V1_5 using SHA-384.
	//
	// See http://tools.ietf.org/html/rfc7518#section-3.3
	RS384

	// RS512 provides a JWT signing method for RSASSA-PKCS1-V1_5 using SHA-512.
	//
	// See http://tools.ietf.org/html/rfc7518#section-3.3
	RS512

	// ES256 provides a JWT signing method for ECDSA using P-256 and SHA-256.
	//
	// See http://tools.ietf.org/html/rfc7518#section-3.4
	ES256

	// ES384 provides a JWT signing method for ECDSA using P-384 and SHA-384.
	//
	// See http://tools.ietf.org/html/rfc7518#section-3.4
	ES384

	// ES512 provides a JWT signing method for ECDSA using P-521 and SHA-512.
	//
	// See http://tools.ietf.org/html/rfc7518#section-3.4
	ES512

	// PS256 provides a JWT signing method for RSASSA-PSS using SHA-256 and
	// MGF1 mask generation function with SHA-256.
	//
	// See http://tools.ietf.org/html/rfc7518#section-3.5
	PS256

	// PS384 provides a JWT signing method for RSASSA-PSS using SHA-384 hash
	// algorithm and MGF1 mask generation function with SHA-384.
	//
	// See http://tools.ietf.org/html/rfc7518#section-3.5
	PS384

	// PS512 provides a JWT signing method for RSASSA-PSS using SHA-512 hash
	// algorithm and MGF1 mask generation function with SHA-512.
	//
	// See http://tools.ietf.org/html/rfc7518#section-3.5
	PS512
)

// algSet is the set of Algorithm implementations.
var algSet = []struct {
	init func(Store, crypto.Hash) (Signer, error)
	hash crypto.Hash
}{
	// none
	NONE: {func(Store, crypto.Hash) (Signer, error) {
		return nil, ErrInvalidAlgorithm
	}, crypto.SHA256},

	// HS256 is HMAC + SHA-256
	HS256: {NewHMACSigner(HS256), crypto.SHA256},

	// HS384 is HMAC + SHA-384
	HS384: {NewHMACSigner(HS384), crypto.SHA384},

	// HS512 is HMAC + SHA-512
	HS512: {NewHMACSigner(HS512), crypto.SHA512},

	// RS256 is RSASSA-PKCS1-V1_5 + SHA-256
	RS256: {NewRSASigner(RS256, RSAMethodPKCS1v15), crypto.SHA256},

	// RS384 is RSASSA-PKCS1-V1_5 + SHA-384
	RS384: {NewRSASigner(RS384, RSAMethodPKCS1v15), crypto.SHA384},

	// RS512 is RSASSA-PKCS1-V1_5 + SHA-512
	RS512: {NewRSASigner(RS512, RSAMethodPKCS1v15), crypto.SHA512},

	// ES256 is ECDSA P-256 + SHA-256
	ES256: {NewEllipticSigner(ES256, elliptic.P256()), crypto.SHA256},

	// ES384 is ECDSA P-384 + SHA-384
	ES384: {NewEllipticSigner(ES384, elliptic.P384()), crypto.SHA384},

	// ES512 is ECDSA P-521 + SHA-512
	ES512: {NewEllipticSigner(ES512, elliptic.P521()), crypto.SHA512},

	// PS256 is RSASSA-PSS + SHA-256
	PS256: {NewRSASigner(PS256, RSAMethodPSS), crypto.SHA256},

	// PS384 is RSASSA-PSS + SHA-384
	PS384: {NewRSASigner(PS384, RSAMethodPSS), crypto.SHA384},

	// PS512 is RSASSA-PSS + SHA-512
	PS512: {NewRSASigner(PS512, RSAMethodPSS), crypto.SHA512},
}

// New creates a Signer using the supplied keyset.
//
// The keyset can be of type []byte, *rsa.{PrivateKey,PublicKey},
// *ecdsa.{PrivateKey,PublicKey}, or compatible with the Store interface.
//
// If a private key is not provided, tokens cannot be Encode'd.  Public keys
// will be automatically generated for RSA and ECC private keys if none were
// provided in the keyset.
func (alg Algorithm) New(keyset interface{}) (Signer, error) {
	a := algSet[alg]

	// check hash
	if !a.hash.Available() {
		return nil, ErrInvalidHash
	}

	var s Store

	// load the data
	switch p := keyset.(type) {
	// regular store
	case Store:
		s = p

	// raw key
	case []byte:
		s = &Keystore{Key: p}

	// rsa keys
	case *rsa.PrivateKey:
		s = &Keystore{Key: p}
	case *rsa.PublicKey:
		s = &Keystore{PubKey: p}

	// ecc keys
	case *ecdsa.PrivateKey:
		s = &Keystore{Key: p}
	case *ecdsa.PublicKey:
		s = &Keystore{PubKey: p}

	default:
		return nil, ErrInvalidKeyset
	}

	return a.init(s, a.hash)
}

// Header builds the JWT header for the algorithm.
func (alg Algorithm) Header() Header {
	return Header{
		Type:      "JWT",
		Algorithm: alg,
	}
}

// Encode serializes a JWT using the Algorithm and Signer.
func (alg Algorithm) Encode(signer Signer, obj interface{}) ([]byte, error) {
	return Encode(alg, signer, obj)
}

// Decode decodes a serialized JWT in buf into obj, and verifies the JWT
// signature using the Algorithm and Signer.
//
// If the token or signature is invalid, ErrInvalidToken or ErrInvalidSignature
// will be returned, respectively. Otherwise, any other errors encountered
// during token decoding will be returned.
func (alg Algorithm) Decode(signer Signer, buf []byte, obj interface{}) error {
	return Decode(alg, signer, buf, obj)
}

// MarshalText marshals Algorithm into a serialized form.
func (alg Algorithm) MarshalText() ([]byte, error) {
	return []byte(alg.String()), nil
}

// UnmarshalText attempts to unmarshal buf into an Algorithm.
func (alg *Algorithm) UnmarshalText(buf []byte) error {
	switch string(buf) {
	// hmac
	case "HS256":
		*alg = HS256
	case "HS384":
		*alg = HS384
	case "HS512":
		*alg = HS512

	// rsa-pkcs1v15
	case "RS256":
		*alg = RS256
	case "RS384":
		*alg = RS384
	case "RS512":
		*alg = RS512

	// ecc
	case "ES256":
		*alg = ES256
	case "ES384":
		*alg = ES384
	case "ES512":
		*alg = ES512

	// rsa-pss
	case "PS256":
		*alg = PS256
	case "PS384":
		*alg = PS384
	case "PS512":
		*alg = PS512

	// error
	default:
		return ErrInvalidAlgorithm
	}

	return nil
}
