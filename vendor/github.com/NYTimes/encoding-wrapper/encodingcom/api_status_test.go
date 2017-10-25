package encodingcom

import (
	"net/http"
	"net/http/httptest"

	"gopkg.in/check.v1"
)

func (s *S) TestAPIStatus(c *check.C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"Encoding Queue Processing Delays","status_code":"queue_slow","incident":"Our encoding queue is processing slower than normal.  Check back for updates."}`))
	}))
	defer server.Close()
	resp, err := APIStatus(server.URL)
	c.Assert(err, check.IsNil)
	c.Assert(*resp, check.DeepEquals, APIStatusResponse{
		Status:     "Encoding Queue Processing Delays",
		StatusCode: "queue_slow",
		Incident:   "Our encoding queue is processing slower than normal.  Check back for updates.",
	})
}

func (s *S) TestAPIStatusFailToConnect(c *check.C) {
	_, err := APIStatus("http://192.0.2.13:8080")
	c.Assert(err, check.NotNil)
}

func (s *S) TestAPIStatusInvalidResponse(c *check.C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{not a valid json}`))
	}))
	defer server.Close()
	_, err := APIStatus(server.URL)
	c.Assert(err, check.NotNil)
}

func (s *S) TestAPIStatusOK(c *check.C) {
	var tests = []struct {
		input string
		want  bool
	}{
		{"ok", true},
		{"encoding_delay", false},
		{"api_out", false},
		{"maintenance", false},
		{"pc_queue_slow", false},
	}
	for _, test := range tests {
		status := APIStatusResponse{StatusCode: test.input}
		got := status.OK()
		c.Check(got, check.Equals, test.want)
	}
}
