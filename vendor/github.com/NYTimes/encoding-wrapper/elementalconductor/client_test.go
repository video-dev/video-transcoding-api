package elementalconductor

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"gopkg.in/check.v1"
)

func (s *S) TestNewClient(c *check.C) {
	expected := Client{
		Host:            "https://mycluster.cloud.elementaltechnologies.com",
		UserLogin:       "myuser",
		APIKey:          "elemental-secret-key",
		AuthExpires:     45,
		AccessKeyID:     "aws-access-key",
		SecretAccessKey: "aws-secret-key",
		Destination:     "destination",
	}
	got := NewClient("https://mycluster.cloud.elementaltechnologies.com", "myuser", "elemental-secret-key", 45, "aws-access-key", "aws-secret-key", "destination")
	c.Assert(*got, check.DeepEquals, expected)
}

func (s *S) TestCreateAuthKey(c *check.C) {
	path := "/jobs"
	userID := "myuser"
	APIKey := "api-key"
	expire := time.Unix(1, 0)
	expireTimestamp := getUnixTimestamp(expire)
	innerKeyMD5 := md5.Sum([]byte(path + userID + APIKey + expireTimestamp))
	innerKey2 := hex.EncodeToString(innerKeyMD5[:])
	value := md5.Sum([]byte(APIKey + innerKey2))
	expected := hex.EncodeToString(value[:])
	client := NewClient("https://mycluster.cloud.elementaltechnologies.com", userID, APIKey, 45, "aws-access-key", "aws-secret-key", "destination")
	got := client.createAuthKey(path, expire)
	c.Assert(got, check.Equals, expected)
}

func (s *S) TestDoRequiredParameters(c *check.C) {
	var req *http.Request
	var data []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req = r
		data, _ = ioutil.ReadAll(r.Body)
		w.Write([]byte(`<response>test</response>`))
	}))
	defer server.Close()
	client := NewClient(server.URL, "myuser", "elemental-secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

	var respObj interface{}
	myJob := Job{
		XMLName: xml.Name{
			Local: "job",
		},
		Input: Input{
			FileInput: Location{
				URI:      "http://another.non.existent/video.mp4",
				Username: "user",
				Password: "pass123",
			},
		},
	}
	err := client.do("POST", "/jobs", myJob, &respObj)

	c.Assert(err, check.IsNil)
	c.Assert(req, check.NotNil)
	c.Assert(req.Method, check.Equals, "POST")
	c.Assert(req.URL.Path, check.Equals, "/api/jobs")
	c.Assert(req.Header.Get("Accept"), check.Equals, "application/xml")
	c.Assert(req.Header.Get("Content-type"), check.Equals, "application/xml")
	c.Assert(req.Header.Get("X-Auth-User"), check.Equals, client.UserLogin)

	c.Assert(req.Header.Get("X-Auth-Expires"), check.NotNil)
	timestampInt, err := strconv.ParseInt(req.Header.Get("X-Auth-Expires"), 10, 64)
	c.Assert(err, check.IsNil)
	timestampTime := time.Unix(timestampInt, 0)

	c.Assert(
		req.Header.Get("X-Auth-Key"),
		check.Equals,
		client.createAuthKey("/jobs", timestampTime),
	)
	var reqJob Job

	err = xml.Unmarshal(data, &reqJob)

	c.Assert(err, check.IsNil)
	c.Assert(reqJob, check.DeepEquals, myJob)
}

func (s *S) TestInvalidAuth(c *check.C) {
	errorResponse := `<?xml version="1.0" encoding="UTF-8"?>
<errors>
  <error>You must be logged in to access this page.</error>
</errors>`
	server, _ := s.startServer(http.StatusUnauthorized, errorResponse)
	defer server.Close()
	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")

	getJobsResponse, err := client.GetJob("1")
	c.Assert(getJobsResponse, check.IsNil)
	c.Assert(err, check.DeepEquals, &APIError{
		Status: http.StatusUnauthorized,
		Errors: errorResponse,
	})
}

func (s *S) TestAPIErrorMarshalling(c *check.C) {
	err := &APIError{
		Status: http.StatusInternalServerError,
		Errors: "something went wrong",
	}
	expectedError := `Error returned by the Elemental Conductor REST Interface: {"status":500,"errors":"something went wrong"}`
	c.Assert(err.Error(), check.Equals, expectedError)
}
