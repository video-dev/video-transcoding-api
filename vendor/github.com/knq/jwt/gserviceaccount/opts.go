package gserviceaccount

import (
	"net/http"
	"net/url"
	"time"
)

// Option is a GServiceAccount option.
type Option func(*GServiceAccount) error

// Transport is a GServiceAccount option to set the client transport used by
// the token source.
func Transport(transport http.RoundTripper) Option {
	return func(gsa *GServiceAccount) error {
		gsa.transport = transport
		return nil
	}
}

// Proxy is a GServiceAccount option to set a HTTP proxy used for by the token
// source.
func Proxy(proxy string) Option {
	return func(gsa *GServiceAccount) error {
		u, err := url.Parse(proxy)
		if err != nil {
			return err
		}
		return Transport(&http.Transport{
			Proxy: http.ProxyURL(u),
		})(gsa)
	}
}

// Expiration is a GServiceAccount option to set a expiration limit for tokens
// generated from the token source.
func Expiration(expiration time.Duration) Option {
	return func(gsa *GServiceAccount) error {
		gsa.expiration = expiration
		return nil
	}
}
