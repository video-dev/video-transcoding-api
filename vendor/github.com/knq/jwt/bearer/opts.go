package bearer

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

// Option represents a Bearer option.
type Option func(*Bearer) error

// ExpiresIn is an option that will set the expiration duration generated for
// tokens to the specified duration.
func ExpiresIn(d time.Duration) Option {
	return func(tok *Bearer) error {
		if d != 0 {
			tok.addExpiration = true
			tok.expiresIn = d
		} else {
			tok.addExpiration = false
			tok.expiresIn = 0
		}

		return nil
	}
}

// IssuedAt is an option that toggles whether or not the Issued At ("iat")
// field is generated for the token.
func IssuedAt(enable bool) Option {
	return func(tok *Bearer) error {
		tok.addIssuedAt = enable
		return nil
	}
}

// NotBefore is an option that toggles whether or not the Not Before ("nbf")
// field is generated for the token.
func NotBefore(enable bool) Option {
	return func(tok *Bearer) error {
		tok.addNotBefore = enable
		return nil
	}
}

// Claim is an option that adds an additional claim that is generated with the
// token.
func Claim(name string, v interface{}) Option {
	return func(tok *Bearer) error {
		if tok.claims == nil {
			return errors.New("attempting to add claim to improperly created token")
		}

		tok.claims[name] = v

		return nil
	}
}

// Scope is an option that adds a Scope ("scope") field to generated tokens.
// Scopes are joined with a space (" ") separator.
func Scope(scopes ...string) Option {
	return func(tok *Bearer) error {
		if len(scopes) > 0 {
			return Claim("scope", strings.Join(scopes, " "))(tok)
		}
		return nil
	}
}

// Transport is an option that sets an underlying client transport to the
// exchange process.
func Transport(transport http.RoundTripper) Option {
	return func(tok *Bearer) error {
		tok.transport = transport
		return nil
	}
}
