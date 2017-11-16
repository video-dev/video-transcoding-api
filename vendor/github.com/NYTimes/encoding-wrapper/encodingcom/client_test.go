package encodingcom

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"gopkg.in/check.v1"
)

func (s *S) mockMediaResponseObject(message string, errors string) interface{} {
	return map[string]interface{}{
		"response": map[string]interface{}{
			"message": message,
			"errors":  map[string]string{"error": errors},
		},
	}
}

func (s *S) TestNewClient(c *check.C) {
	expected := Client{
		Endpoint: "https://manage.encoding.com",
		UserID:   "myuser",
		UserKey:  "secret-key",
	}
	got, err := NewClient("https://manage.encoding.com", "myuser", "secret-key")
	c.Assert(err, check.IsNil)
	c.Assert(*got, check.DeepEquals, expected)
}

func (s *S) TestYesNoBooleanMarshal(c *check.C) {
	bTrue := YesNoBoolean(true)
	bFalse := YesNoBoolean(false)
	data, err := json.Marshal(bTrue)
	c.Assert(err, check.IsNil)
	c.Assert(string(data), check.Equals, `"yes"`)
	data, err = json.Marshal(bFalse)
	c.Assert(err, check.IsNil)
	c.Assert(string(data), check.Equals, `"no"`)
}

func (s *S) TestYesNoBooleanUnmarshal(c *check.C) {
	data := []byte(`{"true":"yes", "false":"no"}`)
	var m map[string]YesNoBoolean
	err := json.Unmarshal(data, &m)
	c.Assert(err, check.IsNil)
	c.Assert(m, check.DeepEquals, map[string]YesNoBoolean{
		"true":  YesNoBoolean(true),
		"false": YesNoBoolean(false),
	})

	invalidData := []byte(`{"true":"true"}`)
	err = json.Unmarshal(invalidData, &m)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, `invalid value: "true"`)
}

func (s *S) TestZeroOneBooleanMarshal(c *check.C) {
	bTrue := ZeroOneBoolean(true)
	bFalse := ZeroOneBoolean(false)
	data, err := json.Marshal(bTrue)
	c.Assert(err, check.IsNil)
	c.Assert(string(data), check.Equals, `"1"`)
	data, err = json.Marshal(bFalse)
	c.Assert(err, check.IsNil)
	c.Assert(string(data), check.Equals, `"0"`)
}

func (s *S) TestZeroOneBooleanUnmarshal(c *check.C) {
	data := []byte(`{"true":"1", "false":"0"}`)
	var m map[string]ZeroOneBoolean
	err := json.Unmarshal(data, &m)
	c.Assert(err, check.IsNil)
	c.Assert(m, check.DeepEquals, map[string]ZeroOneBoolean{
		"true":  ZeroOneBoolean(true),
		"false": ZeroOneBoolean(false),
	})

	invalidData := []byte(`{"true":"true"}`)
	err = json.Unmarshal(invalidData, &m)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, `invalid value: "true"`)
}

func (s *S) TestDoMediaAction(c *check.C) {
	server, requests := s.startServer(`{"response": {"message": "Deleted"}}`)
	defer server.Close()
	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	cancelMediaResponse, err := client.doMediaAction("12345", "CancelMedia")
	c.Assert(err, check.IsNil)
	c.Assert(cancelMediaResponse, check.DeepEquals, &Response{
		Message: "Deleted",
	})
	req := <-requests
	c.Assert(req.query["action"], check.Equals, "CancelMedia")
}

func (s *S) TestDoMediaActionFailure(c *check.C) {
	server, requests := s.startServer(`{"response": {"message": "Deleted", "errors": {"error": "something went wrong"}}}`)
	defer server.Close()
	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	cancelMediaResponse, err := client.doMediaAction("12345", "CancelMedia")
	c.Assert(err, check.NotNil)
	c.Assert(cancelMediaResponse, check.IsNil)
	req := <-requests
	c.Assert(req.query["action"], check.Equals, "CancelMedia")
}

func (s *S) TestDoMissingRequiredParameters(c *check.C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		byteResponse, _ := json.Marshal(s.mockMediaResponseObject("", "Wrong user id or key!"))
		w.Write(byteResponse)
	}))
	defer server.Close()
	client := Client{Endpoint: server.URL}
	err := client.do(&request{
		Action:  "AddMedia",
		MediaID: "123456",
		Source:  []string{"http://some.non.existent/video.mp4"},
	}, nil)
	c.Assert(err, check.NotNil)
	apiErr, ok := err.(*APIError)
	c.Assert(ok, check.Equals, true)
	c.Assert(apiErr.Message, check.Equals, "")
	c.Assert(apiErr.Errors, check.DeepEquals, []string{"Wrong user id or key!"})
}

func (s *S) TestDoMediaResponse(c *check.C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		byteResponse, _ := json.Marshal(s.mockMediaResponseObject("it worked!", ""))
		w.Write(byteResponse)
	}))
	defer server.Close()
	client := Client{Endpoint: server.URL}
	var result map[string]*Response
	err := client.do(&request{
		Action:  "GetStatus",
		MediaID: "123456",
	}, &result)
	c.Assert(err, check.IsNil)
	c.Assert(result["response"].Message, check.Equals, "it worked!")
}

func (s *S) TestDoRequiredParameters(c *check.C) {
	var req *http.Request
	var data string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req = r
		data = req.FormValue("json")
		w.Write([]byte(`{"response": {"status": "added"}}`))
	}))
	defer server.Close()
	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	var respObj map[string]interface{}
	err := client.do(&request{Action: "GetStatus"}, &respObj)
	c.Assert(err, check.IsNil)
	c.Assert(req, check.NotNil)
	c.Assert(req.Method, check.Equals, "POST")
	c.Assert(req.URL.Path, check.Equals, "/")
	c.Assert(req.Header.Get("Content-Type"), check.Equals, "application/x-www-form-urlencoded")
	var m map[string]interface{}
	err = json.Unmarshal([]byte(data), &m)
	c.Assert(err, check.IsNil)
	c.Assert(m, check.DeepEquals, map[string]interface{}{
		"query": map[string]interface{}{
			"userid":  "myuser",
			"userkey": "123",
			"action":  "GetStatus",
		},
	})
	c.Assert(respObj, check.DeepEquals, map[string]interface{}{
		"response": map[string]interface{}{
			"status": "added",
		},
	})
}

func (s *S) TestDoInvalidResponse(c *check.C) {
	server, _ := s.startServer(`{invalid json}`)
	defer server.Close()
	client := Client{Endpoint: server.URL, UserID: "myuser", UserKey: "123"}
	var resp Response
	err := client.do(&request{Action: "GetStatus"}, &resp)
	c.Assert(err, check.NotNil)
}

func (s *S) TestAPIErrorRepresentation(c *check.C) {
	err := &APIError{
		Message: "something went wrong",
		Errors:  []string{"error 1", "error 2"},
	}
	expectedMsg := `Error returned by the Encoding.com API: {"Message":"something went wrong","Errors":["error 1","error 2"]}`
	c.Assert(err.Error(), check.Equals, expectedMsg)
}
