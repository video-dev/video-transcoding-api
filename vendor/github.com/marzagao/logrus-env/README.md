# logrus-env: Environment variables hook for Logrus

Pass-thru hook that sets the values of a given set of environment variables into logrus entries.

## Usage

Initialize the hook by passing the keys of all environment variables you wish to include in your logrus entries. Example:

```go
package main

import (
	"github.com/sirupsen/logrus"
	"github.com/marzagao/logrus-env"
)

func main() {
	logger := logrus.New()
	hook := logrus_env.NewHook([]string{"VARIABLE", "ANOTHER_VARIABLE"})
	logger.Hooks.Add(hook)
}
```

## Result

The result is that if you have a logging statement like this:

```go
	logger.WithFields(logrus.Fields{
		"someKey": "someValue",
	}).Info("something happened here")
```

And you initialize the hook like in the example from the Usage section, the end result will be the same as if you had done:

```go
	logger.WithFields(logrus.Fields{
		"someKey":         "some value",
		"variable":        "value of env variable",
		"anotherVariable": "value of another env variable",
	}).Info("something happened here")
```
