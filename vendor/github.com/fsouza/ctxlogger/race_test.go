// Copyright 2016 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ctxlogger

import (
	"strings"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/tsuru/tsuru/safe"
)

type safeFakeHook struct {
	test.Hook
	mtx sync.Mutex
}

func (h *safeFakeHook) Fire(e *logrus.Entry) error {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	return h.Hook.Fire(e)
}

func TestVarsLoggerIsSafe(t *testing.T) {
	var fakeHook safeFakeHook
	const N = 32
	var b safe.Buffer
	logger := logrus.New()
	logger.Out = &b
	logger.Level = logrus.DebugLevel
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Hooks.Add(&fakeHook)
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i int) {
			innerLogger := varsLogger(map[string]string{"name": "gopher"}, logger)
			innerLogger.WithField("some", "thing").Info("be advised")
			wg.Done()
		}(i)
	}
	wg.Wait()
	logLines := strings.Split(strings.TrimSpace(b.String()), "\n")
	if len(logLines) != N {
		t.Errorf("wrong log lines returned, wanted %d log lines, got %d:\n%s", N, len(logLines), b.String())
	}
	if len(fakeHook.Entries) != N {
		t.Errorf("wrong number of entries in the hook. want %d, got %d", N, len(fakeHook.Entries))
	}
}
