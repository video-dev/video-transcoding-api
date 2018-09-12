// Package gserviceaccount provides a simple way to load Google service account
// credentials and create a corresponding oauth2.TokenSource from it.
package gserviceaccount

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/knq/jwt"
	"github.com/knq/jwt/bearer"
	"github.com/knq/pemutil"
)

const (
	// DefaultAlgorithm is the default jwt.Algothrithm to use with service
	// account tokens.
	DefaultAlgorithm = jwt.RS256

	// DefaultExpiration is the default token expiration duration to use with
	// service account tokens.
	DefaultExpiration = 1 * time.Hour
)

// GServiceAccount wraps Google Service Account parameters, and are the same
// values found in a standard JSON-encoded credentials file provided by Google.
type GServiceAccount struct {
	Type                    string `json:"type,omitempty"`
	ProjectID               string `json:"project_id,omitempty"`
	PrivateKeyID            string `json:"private_key_id,omitempty"`
	PrivateKey              string `json:"private_key,omitempty"`
	ClientEmail             string `json:"client_email,omitempty"`
	ClientID                string `json:"client_id,omitempty"`
	AuthURI                 string `json:"auth_uri,omitempty"`
	TokenURI                string `json:"token_uri,omitempty"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url,omitempty"`
	ClientX509CertURL       string `json:"client_x509_cert_url,omitempty"`

	expiration time.Duration     `json:"-"`
	signer     jwt.Signer        `json:"-"`
	transport  http.RoundTripper `json:"-"`
	mu         sync.Mutex        `json:"-"`
}

// FromJSON loads service account credentials from the JSON encoded buf.
func FromJSON(buf []byte, opts ...Option) (*GServiceAccount, error) {
	var err error

	// unmarshal
	gsa := new(GServiceAccount)
	if err = json.Unmarshal(buf, gsa); err != nil {
		return nil, err
	}

	// apply opts
	for _, o := range opts {
		if err = o(gsa); err != nil {
			return nil, err
		}
	}

	return gsa, nil
}

// FromReader loads Google service account credentials from a reader.
func FromReader(r io.Reader, opts ...Option) (*GServiceAccount, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return FromJSON(buf, opts...)
}

// FromFile loads Google service account credentials from a reader.
func FromFile(path string, opts ...Option) (*GServiceAccount, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return FromJSON(buf, opts...)
}

// Signer returns a jwt.Signer for use when signing tokens.
func (gsa *GServiceAccount) Signer() (jwt.Signer, error) {
	gsa.mu.Lock()
	defer gsa.mu.Unlock()

	if gsa.signer == nil {
		keyset, err := pemutil.DecodeBytes([]byte(gsa.PrivateKey))
		if err != nil {
			return nil, fmt.Errorf("jwt/gserviceaccount: could not decode private key: %v", err)
		}
		keyset.AddPublicKeys()

		s, err := DefaultAlgorithm.New(keyset)
		if err != nil {
			return nil, err
		}
		gsa.signer = s
	}

	return gsa.signer, nil
}

// TokenSource returns a oauth2.TokenSource for the Google Service Account
// using the provided context and scopes. The resulting token source should be
// wrapped with oauth2.ReusableTokenSource prior to being used elsewhere.
//
// If the supplied context is nil, context.Background() will be used.
//
// If additional claims need to be added to the TokenSource (ie, subject or the
// "sub" field), use jwt/bearer.Claim to add them before wrapping the
// TokenSource with oauth2.ReusableTokenSource.
func (gsa *GServiceAccount) TokenSource(ctxt context.Context, scopes ...string) (*bearer.Bearer, error) {
	var err error

	// simple check that required fields are present
	if gsa.ClientEmail == "" || gsa.TokenURI == "" {
		return nil, errors.New("jwt/gserviceaccount: ClientEmail and TokenURI cannot be empty")
	}

	// set up subject and context
	if ctxt == nil {
		ctxt = context.Background()
	}

	// get signer
	signer, err := gsa.Signer()
	if err != nil {
		return nil, err
	}

	// determine expiration
	expiration := gsa.expiration
	if expiration == 0 {
		expiration = DefaultExpiration
	}

	// bearer grant options
	opts := []bearer.Option{
		bearer.ExpiresIn(expiration),
		bearer.IssuedAt(true),
		bearer.Claim("iss", gsa.ClientEmail),
		bearer.Claim("aud", gsa.TokenURI),
		bearer.Scope(scopes...),
	}

	// add transport
	if gsa.transport != nil {
		opts = append(opts, bearer.Transport(gsa.transport))
	}

	// create token source
	b, err := bearer.NewTokenSource(
		signer,
		gsa.TokenURI,
		ctxt,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Client returns a HTTP client using the provided context and scopes for the
// service account as the underlying transport.
//
// When called with the appropriate scopes, the created client can be passed to
// any Google API for creating a service client:
//
// 		import (
// 			dns "google.golang.org/api/dns/v2beta1"
//      )
//      cl, err := gsa.Client(ctxt, dns.CloudPlatformScope, dns.NdevClouddnsReadwriteScope)
// 		if err != nil { /* ... */ }
//      dnsService, err := dns.New(cl)
// 		if err != nil { /* ... */ }
//
// Note: this is a convenience func only.
func (gsa *GServiceAccount) Client(ctxt context.Context, scopes ...string) (*http.Client, error) {
	b, err := gsa.TokenSource(ctxt, scopes...)
	if err != nil {
		return nil, err
	}
	return b.Client(), nil
}
