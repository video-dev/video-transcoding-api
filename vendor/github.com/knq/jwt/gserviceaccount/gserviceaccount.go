// Package gserviceaccount provides a simple way to load Google service account
// credentials and create a corresponding oauth2.TokenSource from it.
package gserviceaccount

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"time"

	"golang.org/x/net/context"

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

	signer jwt.Signer `json:"-"`
	mu     sync.Mutex `json:"-"`
}

// FromJSON loads service account credentials from the JSON encoded buf.
func FromJSON(buf []byte) (*GServiceAccount, error) {
	var gsa GServiceAccount

	// unmarshal
	err := json.Unmarshal(buf, &gsa)
	if err != nil {
		return nil, err
	}

	return &gsa, nil
}

// FromReader loads Google service account credentials from a reader.
func FromReader(r io.Reader) (*GServiceAccount, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return FromJSON(buf)
}

// FromFile loads Google service account credentials from a reader.
func FromFile(path string) (*GServiceAccount, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return FromJSON(buf)
}

// Signer returns a jwt.Signer for use when signing tokens.
func (gsa *GServiceAccount) Signer() (jwt.Signer, error) {
	gsa.mu.Lock()
	defer gsa.mu.Unlock()

	if gsa.signer == nil {
		keyset := pemutil.Store{}
		err := keyset.Decode([]byte(gsa.PrivateKey))
		if err != nil {
			return nil, fmt.Errorf("jwt/gserviceaccount: could not decode private key: %v", err)
		}

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

	// bearer grant options
	opts := []bearer.Option{
		bearer.ExpiresIn(DefaultExpiration),
		bearer.IssuedAt(true),
		bearer.Claim("iss", gsa.ClientEmail),
		bearer.Claim("aud", gsa.TokenURI),
		bearer.Scope(scopes...),
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
