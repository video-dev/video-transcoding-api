package jwt

import (
	"bytes"
	"encoding/json"
)

// Header is a JWT header.
type Header struct {
	Type      string    `json:"typ"`
	Algorithm Algorithm `json:"alg"`
}

// Claims contains the registered JWT claims.
//
// See: https://tools.ietf.org/html/rfc7519#section-4.1
type Claims struct {
	// Issuer ("iss") identifies the principal that issued the JWT.
	Issuer string `json:"iss,omitempty"`

	// Subject ("sub") identifies the principal that is the subject of the JWT.
	Subject string `json:"sub,omitempty"`

	// Audience ("aud") identifies the recipients that the JWT is intended for.
	Audience string `json:"aud,omitempty"`

	// Expiration ("exp") identifies the expiration time on or after which the
	// JWT MUST NOT be accepted for processing.
	Expiration json.Number `json:"exp,omitempty"`

	// NotBefore ("nbf") identifies the time before which the JWT MUST NOT be
	// accepted for processing.
	NotBefore json.Number `json:"nbf,omitempty"`

	// IssuedAt ("iat") identifies the time at which the JWT was issued.
	IssuedAt json.Number `json:"iat,omitempty"`

	// JwtID ("jti") provides a unique identifier for the JWT.
	JwtID string `json:"jti,omitempty"`
}

// Token is a JWT token, comprising header, payload (ie, claims), and
// signature.
type Token struct {
	Header    Header `jwt:"header"`
	Payload   Claims `jwt:"payload"`
	Signature []byte `jwt:"signature"`
}

// UnverifiedToken is a token split into its composite parts, but has not yet
// been decoded, nor verified.
type UnverifiedToken struct {
	Header, Payload, Signature []byte
}

// DecodeUnverifiedToken decodes a token into the provided UnverifiedToken.
func DecodeUnverifiedToken(buf []byte, ut *UnverifiedToken) error {
	b := bytes.Split(buf, tokenSep)
	if len(b) != 3 {
		return ErrInvalidToken
	}

	ut.Header = b[0]
	ut.Payload = b[1]
	ut.Signature = b[2]

	return nil
}
