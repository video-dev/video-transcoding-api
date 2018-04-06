// Copyright 2018 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logging

import (
	"os"

	"github.com/Gurpartap/logrus-stack"
	"github.com/knq/sdhook"
	"github.com/marzagao/logrus-env"
	"github.com/sirupsen/logrus"
)

// Config contains configuration for logging level and services integration.
type Config struct {
	Level string `envconfig:"LOGGING_LEVEL" default:"info"`

	// List of environment variables that should be included in all log
	// lines.
	EnvironmentVariables []string `envconfig:"LOGGING_ENVIRONMENT_VARIABLES"`

	// Send logs to StackDriver?
	SendToStackDriver bool `envconfig:"LOGGING_SEND_TO_STACKDRIVER"`

	// StackDriver error reporting options. When present, error logs are
	// going to be reported as errors on StackDriver.
	StackDriverErrorServiceName string `envconfig:"LOGGING_STACKDRIVER_ERROR_SERVICE_NAME"`
	StackDriverErrorLogName     string `envconfig:"LOGGING_STACKDRIVER_ERROR_LOG_NAME" default:"error_log"`
}

// Logger returns a logrus logger with the features defined in the config.
func (c *Config) Logger() (*logrus.Logger, error) {
	level, err := logrus.ParseLevel(c.Level)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.Out = os.Stdout
	logger.Level = level
	logger.Hooks.Add(logrus_stack.StandardHook())
	logger.Hooks.Add(logrus_env.NewHook(c.EnvironmentVariables))
	if c.SendToStackDriver {
		opts := []sdhook.Option{sdhook.GoogleLoggingAgent()}
		if c.StackDriverErrorServiceName != "" {
			opts = append(opts,
				sdhook.ErrorReportingService(c.StackDriverErrorServiceName),
				sdhook.ErrorReportingLogName(c.StackDriverErrorLogName),
			)
		}
		gcpLoggingHook, err := sdhook.New(opts...)
		if err != nil {
			return nil, err
		}
		logger.Hooks.Add(gcpLoggingHook)
	}
	return logger, nil
}
