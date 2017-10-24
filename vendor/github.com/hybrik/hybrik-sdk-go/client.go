package hybrik

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

// APIInterface is interface for the underlying client object
type APIInterface interface {
	connect() error
	isExpired() bool
	CallAPI(method string, apiPath string, params url.Values, body io.Reader) (string, error)
}

// ClientInterface is an interface for commonly used API calls
type ClientInterface interface {
	// Generic
	CallAPI(method string, apiPath string, params url.Values, body io.Reader) (string, error)

	// Jobs API
	QueueJob(string) (string, error)
	GetJobInfo(string) (JobInfo, error)
	StopJob(string) error

	// Presets API
	GetPreset(string) (Preset, error)
	CreatePreset(Preset) (Preset, error)
	DeletePreset(string) error
}

// Client implements the ClientInterface
type Client struct {
	client APIInterface
}

// Config represents the configuration params that are necessary for the API calls to Hybrik to work
type Config struct {
	URL            string
	ComplianceDate string
	OAPIKey        string
	OAPISecret     string
	AuthKey        string
	AuthSecret     string
	OAPIURL        string
}

// API is the implementation of the HybrikAPI methods
type API struct {
	Config     Config
	token      string
	expiration string
}

type connectResponse struct {
	Token          string `json:"token"`
	ExpirationTime string `json:"expiration_time"`
}

// NewClient creates an instance of the HybrikAPI client
func NewClient(config Config) (*Client, error) {

	switch {
	case config.URL == "":
		return &Client{}, ErrNoAPIURL
	case config.OAPIKey == "":
		return &Client{}, ErrNoOAPIKey
	case config.OAPISecret == "":
		return &Client{}, ErrNoOAPISecret
	case config.AuthKey == "":
		return &Client{}, ErrNoAuthKey
	case config.AuthSecret == "":
		return &Client{}, ErrNoAuthSecret
	case config.ComplianceDate == "":
		return &Client{}, ErrNoComplianceDate
	case !regexp.MustCompile(`^\d{8}$`).MatchString(config.ComplianceDate):
		return &Client{}, ErrNoComplianceDate
	}

	_, err := url.ParseRequestURI(config.URL)
	if err != nil {
		return &Client{}, ErrInvalidURL
	}

	parts := strings.Split(config.URL, "//")
	if len(parts) < 2 {
		return &Client{}, ErrInvalidURL
	}

	config.OAPIURL = fmt.Sprintf("%s//%s:%s@%s",
		parts[0], config.OAPIKey, config.OAPISecret, parts[1],
	)

	return &Client{
		client: &API{
			Config: config,
		},
	}, nil
}

func (a *API) connect() error {

	data := make(map[string]string)
	data["auth_key"] = a.Config.AuthKey
	data["auth_secret"] = a.Config.AuthSecret

	jsonBody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/login", a.Config.OAPIURL), bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("X-Hybrik-Compliance", a.Config.ComplianceDate)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	callResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Login failed w/ %d : %s", resp.StatusCode, string(callResp))
	}

	var cr connectResponse
	err = json.Unmarshal(callResp, &cr)
	if err != nil {
		return err
	}

	a.token = cr.Token
	a.expiration = cr.ExpirationTime

	return nil
}

func (a *API) isExpired() bool {
	if a.expiration == "" {
		return true
	}

	t, err := time.Parse(`2006-01-02T15:04:05.999Z`, a.expiration)
	if err != nil {
		return true
	}

	return time.Now().After(t)
}

var connectLock sync.Mutex

// CallAPI is the general method to call for access to the API
func (a *API) CallAPI(method string, apiPath string, params url.Values, body io.Reader) (string, error) {
	// Retrieves the 'token' and 'expiration_time'
	connectLock.Lock()
	if a.isExpired() {
		err := a.connect()
		if err != nil {
			connectLock.Unlock()
			return "", err
		}
	}
	connectLock.Unlock()

	// Does the necessary http call here
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", a.Config.OAPIURL, apiPath), body)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Hybrik-Compliance", a.Config.ComplianceDate)
	req.Header.Set("X-Hybrik-Sapiauth", a.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.URL.RawQuery = params.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	callResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%d - %s %s: %s", resp.StatusCode, method, apiPath, string(callResp))
	}

	return string(callResp), nil
}

// CallAPI can be used to make any GET/POST/PUT/DELETE API call to the Hybrik service
func (c *Client) CallAPI(method string, apiPath string, params url.Values, body io.Reader) (string, error) {
	return c.client.CallAPI(method, apiPath, params, body)
}
