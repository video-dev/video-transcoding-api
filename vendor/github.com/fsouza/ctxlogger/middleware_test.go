// Copyright 2016 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ctxlogger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/NYTimes/gizmo/web"
	"github.com/sirupsen/logrus"
)

func TestContextLoggerMiddleware(t *testing.T) {
	var tests = []struct {
		testCase    string
		header      http.Header
		vars        map[string]string
		wantLogLine map[string]string
	}{
		{
			"request with vars",
			http.Header{},
			map[string]string{"jobId": "some-job", "limit": "3"},
			map[string]string{
				"level": "info",
				"msg":   "received a nice request!",
				"jobId": "some-job",
				"limit": "3",
			},
		},
		{
			"request with vars and id",
			http.Header{"X-Request-Id": []string{"request-123"}},
			map[string]string{"jobId": "some-job", "limit": "3"},
			map[string]string{
				"level":     "info",
				"msg":       "received a nice request!",
				"jobId":     "some-job",
				"limit":     "3",
				"requestId": "request-123",
			},
		},
		{
			"request with vars and id - override",
			http.Header{"X-Request-Id": []string{"request-123"}},
			map[string]string{"jobId": "some-job", "limit": "3", "requestId": "request-1234"},
			map[string]string{
				"level":     "info",
				"msg":       "received a nice request!",
				"jobId":     "some-job",
				"limit":     "3",
				"requestId": "request-123",
			},
		},
		{
			"request with no vars",
			http.Header{},
			nil,
			map[string]string{
				"level": "info",
				"msg":   "received a nice request!",
			},
		},
		{
			"request with no vars and id",
			http.Header{"X-Request-Id": []string{"request-123"}},
			nil,
			map[string]string{
				"level":     "info",
				"msg":       "received a nice request!",
				"requestId": "request-123",
			},
		},
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := r.Context().Value(ContextKey).(*logrus.Logger)
		logger.Info("received a nice request!")
		w.WriteHeader(http.StatusOK)
	})
	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/something", nil)
			req.Header = test.header
			web.SetRouteVars(req, test.vars)
			rec := httptest.NewRecorder()
			var b bytes.Buffer
			logger := logrus.New()
			logger.Out = &b
			logger.Level = logrus.DebugLevel
			logger.Formatter = &logrus.JSONFormatter{}
			middleware := ContextLogger(logger)
			h := middleware(handler)
			h.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Errorf("wrong status code\ngot  %d\nwant %d", rec.Code, http.StatusOK)
			}
			var logLine map[string]string
			err := json.Unmarshal(bytes.TrimSpace(b.Bytes()), &logLine)
			if err != nil {
				t.Fatal(err)
			}
			test.wantLogLine["time"] = logLine["time"]
			if !reflect.DeepEqual(logLine, test.wantLogLine) {
				t.Errorf("wrong log line returned\ngot  %#v\nwant %#v", logLine, test.wantLogLine)
			}
		})
	}
}

func TestVarsLogger(t *testing.T) {
	var b bytes.Buffer
	logger := logrus.New()
	logger.Out = &b
	logger.Level = logrus.DebugLevel
	logger.Formatter = &logrus.JSONFormatter{}
	vl := varsLogger(map[string]string{"jobId": "some-job", "whatever": "yes!"}, logger)
	if vl == nil {
		t.Fatal("unexpected <nil> logger")
	}
	vl.Error("something went wrong")
	vl.Warning("something will go wrong")
	vl.Debug("what a lovely day")
	vl.WithField("whatever", "no!").Error("something went wrong again")
	lines := strings.Split(strings.TrimSpace(b.String()), "\n")
	if len(lines) != 4 {
		t.Fatalf("wrong number of line. Got %d. Want %d", len(lines), 4)
	}
	rawData := fmt.Sprintf("[%s]", strings.Join(lines, ","))
	var logMsgs []map[string]string
	err := json.Unmarshal([]byte(rawData), &logMsgs)
	if err != nil {
		t.Fatal(err)
	}
	for _, msg := range logMsgs {
		if msg["time"] == "" {
			t.Errorf("unexpected empty time in log entry: %#v", msg)
		}
		delete(msg, "time")
	}
	expectedMsgs := []map[string]string{
		{"level": "error", "jobId": "some-job", "msg": "something went wrong", "whatever": "yes!"},
		{"level": "warning", "jobId": "some-job", "msg": "something will go wrong", "whatever": "yes!"},
		{"level": "debug", "jobId": "some-job", "msg": "what a lovely day", "whatever": "yes!"},
		{"level": "error", "jobId": "some-job", "msg": "something went wrong again", "whatever": "no!"},
	}
	if !reflect.DeepEqual(logMsgs, expectedMsgs) {
		t.Errorf("wrong messages returned\nwant %#v\ngot  %#v", expectedMsgs, logMsgs)
	}
}
