package zencoder

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Zencoder struct {
	BaseUrl string
	Header  http.Header
	Client  *http.Client
}

func NewZencoder(apiKey string) *Zencoder {
	return &Zencoder{
		Client:  http.DefaultClient,
		BaseUrl: "https://app.zencoder.com/api/v2/",
		Header: http.Header{
			"Content-Type":     []string{"application/json"},
			"Accept":           []string{"application/json"},
			"Zencoder-Api-Key": []string{apiKey},
			"User-Agent":       []string{"gozencoder v1"},
		},
	}
}

func (z *Zencoder) call(method, path string, request interface{}, expectedStatus []int) (*http.Response, error) {
	var buffer io.Reader
	if request != nil {
		b, err := json.Marshal(request)
		if err != nil {
			return nil, err
		}

		buffer = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", z.BaseUrl, path), buffer)
	if err != nil {
		return nil, err
	}

	req.Header = z.Header

	resp, err := z.Client.Do(req)
	if err != nil {
		return resp, err
	}

	for _, status := range expectedStatus {
		if resp.StatusCode == status {
			return resp, err
		}
	}

	return nil, errors.New(resp.Status)
}

func (z *Zencoder) post(path string, request interface{}, response interface{}) error {
	resp, err := z.call("POST", path, request, []int{http.StatusCreated, http.StatusOK})
	if err != nil {
		return err
	}

	if err := UnmarshalBody(resp.Body, response); err != nil {
		return err
	}

	return nil
}

func (z *Zencoder) putNoContent(path string) error {
	_, err := z.call("PUT", path, nil, []int{http.StatusNoContent})
	if err != nil {
		return err
	}

	return nil
}

func (z *Zencoder) getBody(path string, response interface{}) error {
	resp, err := z.call("GET", path, nil, []int{http.StatusOK})
	if err != nil {
		return err
	}

	if err := UnmarshalBody(resp.Body, response); err != nil {
		return err
	}

	return nil
}

func UnmarshalBody(body io.ReadCloser, result interface{}) error {
	defer body.Close()

	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, result)
}
