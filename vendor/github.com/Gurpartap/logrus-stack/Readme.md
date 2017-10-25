# logrus-stack 🎯
[![GoDoc](https://godoc.org/github.com/Gurpartap/logrus-stack?status.svg)](https://godoc.org/github.com/Gurpartap/logrus-stack)

logrus-stack provides [facebookgo/stack](https://github.com/facebookgo/stack) integration hook for [sirupsen/logrus](https://github.com/sirupsen/logrus).

Instead of setting file, line, and func name values individually, this hook sets "caller" and/or "stack" objects containing file, line and func name.

The values play well with `logrus.TextFormatter{}` as well as `logrus.JSONFormatter{}`. See example outputs below.

## Doc

There's not much to it. See usage and [GoDoc](https://godoc.org/github.com/Gurpartap/logrus-stack).

## Usage

```bash
$ go get github.com/Gurpartap/logrus-stack
```

```go
logrus.AddHook(logrus_stack.StandardHook())
```

Same as:
```go
callerLevels := logrus.AllLevels
stackLevels := []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel}

logrus.AddHook(logrus_stack.NewHook(callerLevels, stackLevels))
```

## Example

```go
package main

import (
	"errors"
	"os"

	"github.com/Gurpartap/logrus-stack"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	JobID string
}

func (w Worker) Perform() {
	logrus.WithField("jod_id", w.JobID).Infoln("Now working")

	err := errors.New("I don't know what to do yet")
	if err != nil {
		logrus.Errorln(err)
		return
	}

	// ...
}

func main() {
	// Setup logrus.
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stderr)

	// Add the stack hook.
	logrus.AddHook(logrus_stack.StandardHook())

	// Let's try it.
	Worker{"123"}.Perform()
}
```

```bash
$ go run example/main.go
```

```go
{"caller":{"File":"github.com/Gurpartap/logrus-stack/example/main.go","Line":16,"Name":"Worker.Perform"},"jod_id":"123","level":"info","msg":"Now working","time":"2016-10-10T01:17:40+05:30"}
{"caller":{"File":"github.com/Gurpartap/logrus-stack/example/main.go","Line":20,"Name":"Worker.Perform"},"level":"error","msg":"I don't know what to do yet","stack":[{"File":"github.com/Gurpartap/logrus-stack/example/main.go","Line":20,"Name":"Worker.Perform"},{"File":"github.com/Gurpartap/logrus-stack/example/main.go","Line":36,"Name":"main"},{"File":"/usr/local/Cellar/go/1.7.1/libexec/src/runtime/proc.go","Line":183,"Name":"main"},{"File":"/usr/local/Cellar/go/1.7.1/libexec/src/runtime/asm_amd64.s","Line":2086,"Name":"goexit"}],"time":"2016-10-10T01:17:40+05:30"}
```

###### Same as above but indented:

```json
{
	"caller": {
		"File": "github.com/Gurpartap/logrus-stack/example/main.go",
		"Line": 16,
		"Name": "Worker.Perform"
	},
	"jod_id": "123",
	"level": "info",
	"msg": "Now working",
	"time": "2016-10-10T09:41:00+05:30"
}
{
	"caller": {
		"File": "github.com/Gurpartap/logrus-stack/example/main.go",
		"Line": 20,
		"Name": "Worker.Perform"
	},
	"level": "error",
	"msg": "I don't know what to do yet",
	"stack": [
		{
			"File": "github.com/Gurpartap/logrus-stack/example/main.go",
			"Line": 20,
			"Name": "Worker.Perform"
		},
		{
			"File": "github.com/Gurpartap/logrus-stack/example/main.go",
			"Line": 36,
			"Name": "main"
		},
		{
			"File": "/usr/local/Cellar/go/1.7.1/libexec/src/runtime/proc.go",
			"Line": 183,
			"Name": "main"
		},
		{
			"File": "/usr/local/Cellar/go/1.7.1/libexec/src/runtime/asm_amd64.s",
			"Line": 2086,
			"Name": "goexit"
		}
	],
	"time": "2016-10-10T09:41:00+05:30"
}

```

If the same example was used with `logrus.SetFormatter(&logrus.TextFormatter{})` instead, the output would be:

```bash
INFO[0000] Now working                                   caller=github.com/Gurpartap/logrus-stack/example/main.go:16 Worker.Perform jod_id=123
ERRO[0000] I don't know what to do yet                   caller=github.com/Gurpartap/logrus-stack/example/main.go:20 Worker.Perform stack=github.com/Gurpartap/logrus-stack/example/main.go:20            Worker.Perform
github.com/Gurpartap/logrus-stack/example/main.go:36            main
/usr/local/Cellar/go/1.7.1/libexec/src/runtime/proc.go:183      main
/usr/local/Cellar/go/1.7.1/libexec/src/runtime/asm_amd64.s:2086 goexit
```

Hello 👋

Follow me on https://twitter.com/Gurpartap
