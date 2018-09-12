package jwt

// Error is a jwt error.
type Error string

// Error satisfies the error interface.
func (err Error) Error() string {
	return string(err)
}

// Error values.
const (
	// ErrInvalidSignature is the invalid signature error.
	ErrInvalidSignature Error = "invalid signature"

	// ErrInvalidAlgorithm is the invalid algorithm error.
	ErrInvalidAlgorithm Error = "invalid algorithm"

	// ErrInvalidToken is the invalid token error.
	ErrInvalidToken Error = "invalid token"

	// ErrInvalidKeyset is the invalid keyset error.
	ErrInvalidKeyset Error = "invalid keyset"

	// ErrInvalidHash is the invalid hash error.
	ErrInvalidHash Error = "invalid hash"

	// ErrMissingPrivateKey is the missing private key error.
	ErrMissingPrivateKey Error = "missing private key"

	// ErrMissingPublicKey is the missing public key error.
	ErrMissingPublicKey Error = "missing public key"

	// ErrMissingPrivateOrPublicKey is the missing private or public key error.
	ErrMissingPrivateOrPublicKey Error = "missing private or public key"

	// ErrInvalidPrivateKey is the invalid private key error.
	ErrInvalidPrivateKey Error = "invalid private key"

	// ErrInvalidPublicKey is the invalid public key error.
	ErrInvalidPublicKey Error = "invalid public key"

	// ErrInvalidPrivateKeySize is the invalid private key size error.
	ErrInvalidPrivateKeySize Error = "invalid private key size"

	// ErrMismatchedBytesCopied is the mismatched bytes copied error.
	ErrMismatchedBytesCopied Error = "mismatched bytes copied"

	// ErrInvalidPublicKeySize is the invalid public key size error.
	ErrInvalidPublicKeySize Error = "invalid public key size"
)
