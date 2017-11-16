package logrus_env

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type LogrusEnvHook struct {
	Values map[string]string
}

func NewHook(keys []string) LogrusEnvHook {
	newHook := LogrusEnvHook{
		Values: map[string]string{},
	}
	for _, key := range keys {
		if value, found := os.LookupEnv(key); found {
			// convert environment variable key to from the usual
			// uppercase/snakecase to lowercase/camelcase
			// for structured logging
			fields := strings.Split(strings.ToLower(key), "_")
			for idx, f := range fields {
				if idx > 0 {
					fields[idx] = strings.Title(f)
				}
			}
			camelCaseKey := strings.Join(fields, "")
			newHook.Values[camelCaseKey] = value
		}
	}
	return newHook
}

// Levels provides the levels to run the hook on.
func (hook LogrusEnvHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire is called by logrus when something is logged.
func (hook LogrusEnvHook) Fire(entry *logrus.Entry) error {
	for key, value := range hook.Values {
		entry.Data[key] = value
	}
	return nil
}
