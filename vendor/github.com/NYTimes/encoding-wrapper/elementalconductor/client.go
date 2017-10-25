// Package elementalconductor provides types and methods for interacting with the
// Elemental Conductor API.
//
// You can get more details on the API at https://<elemental_server>/help/rest_api.
package elementalconductor

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Client is the basic type for interacting with the API. It provides methods
// matching the available actions in the API.
type Client struct {
	Host            string
	UserLogin       string
	APIKey          string
	AuthExpires     int
	AccessKeyID     string
	SecretAccessKey string
	Destination     string
}

// APIError represents an error returned by the Elemental Cloud REST API.
//
// See https://<elemental_server>/help/rest_api#rest_basics_errors_and_warnings
// for more details.
type APIError struct {
	Status int    `json:"status,omitempty"`
	Errors string `json:"errors,omitempty"`
}

// Error converts the whole interlying information to a representative string.
//
// It encodes the list of errors in JSON format.
func (apiErr *APIError) Error() string {
	data, _ := json.Marshal(apiErr)
	return fmt.Sprintf("Error returned by the Elemental Conductor REST Interface: %s", data)
}

// NewClient creates a instance of the client type.
func NewClient(host, userLogin, apiKey string, authExpires int, accessKeyID string, secretAccessKey string, destination string) *Client {
	return &Client{
		Host:            host,
		UserLogin:       userLogin,
		APIKey:          apiKey,
		AuthExpires:     authExpires,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		Destination:     destination,
	}
}

func getUnixTimestamp(givenTime time.Time) string {
	return strconv.FormatInt(givenTime.UTC().Unix(), 10)
}

func (c *Client) do(method string, path string, body interface{}, out interface{}) error {
	apiPath := "/api" + path
	xmlRequest, err := xml.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(method, c.Host+apiPath, strings.NewReader(string(xmlRequest)))
	if err != nil {
		return err
	}
	expiresTime := time.Now().Add(time.Duration(c.AuthExpires) * time.Second)
	expiresTimestamp := getUnixTimestamp(expiresTime)
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-type", "application/xml")
	req.Header.Set("X-Auth-User", c.UserLogin)
	req.Header.Set("X-Auth-Expires", expiresTimestamp)
	req.Header.Set("X-Auth-Key", c.createAuthKey(path, expiresTime))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return &APIError{
			Status: resp.StatusCode,
			Errors: string(respData),
		}
	}
	if out != nil && len(respData) > 1 {
		return xml.Unmarshal(respData, out)
	}
	return nil
}

func (c *Client) createAuthKey(URL string, expire time.Time) string {
	expireString := getUnixTimestamp(expire)
	hasher := md5.New()
	hasher.Write([]byte(URL))
	hasher.Write([]byte(c.UserLogin))
	hasher.Write([]byte(c.APIKey))
	hasher.Write([]byte(expireString))
	innerKey := hex.EncodeToString(hasher.Sum(nil))
	hasher = md5.New()
	hasher.Write([]byte(c.APIKey))
	hasher.Write([]byte(innerKey))
	return hex.EncodeToString(hasher.Sum(nil))
}
