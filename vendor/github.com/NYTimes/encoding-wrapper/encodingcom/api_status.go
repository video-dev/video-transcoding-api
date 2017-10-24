package encodingcom

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"
)

// APIStatusResponse is the response returned by the APIStatus function.
//
// It describes the current status of the Encoding.com API.
type APIStatusResponse struct {
	Status     string `json:"status"`
	StatusCode string `json:"status_code"`
	Incident   string `json:"incident"`
}

// OK returns whether the given status represents no problem in the
// Encoding.com API.
func (s *APIStatusResponse) OK() bool {
	return s.StatusCode == "ok"
}

// APIStatus queries the current status of the Encoding.com API.
//
// The host parameter is optional, and when omitted, will default to
// "http://status.encoding.com".
//
// See http://goo.gl/3JKSxy for more details.
func APIStatus(endpoint string) (*APIStatusResponse, error) {
	client := http.Client{
		Transport: &http.Transport{
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			ResponseHeaderTimeout: 2 * time.Second,
		},
	}
	url := strings.TrimRight(endpoint, "/") + "/status.php?format=json"
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var apiResp APIStatusResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return nil, err
	}
	return &apiResp, nil
}
