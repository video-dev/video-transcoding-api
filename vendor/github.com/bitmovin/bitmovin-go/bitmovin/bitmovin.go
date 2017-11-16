package bitmovin

import (
	"net/http"
	"time"
)

type Bitmovin struct {
	HTTPClient *http.Client
	APIKey     *string
	APIBaseURL *string
}

func NewBitmovinDefaultTimeout(apiKey string, baseURL string) *Bitmovin {
	return NewBitmovin(apiKey, baseURL, 5)
}

func NewBitmovinDefault(apiKey string) *Bitmovin {
	return NewBitmovinDefaultTimeout(apiKey, `https://api.bitmovin.com/v1/`)
}

func NewBitmovin(apiKey string, baseURL string, timeout int64) *Bitmovin {
	return &Bitmovin{
		HTTPClient: &http.Client{
			Timeout: time.Second * time.Duration(timeout),
		},
		APIKey:     &apiKey,
		APIBaseURL: &baseURL,
	}
}
