package sdhook

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"cloud.google.com/go/compute/metadata"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/knq/jwt/gserviceaccount"

	"github.com/sirupsen/logrus"

	errorReporting "google.golang.org/api/clouderrorreporting/v1beta1"
	logging "google.golang.org/api/logging/v2"
)

// Option represents an option that modifies the Stackdriver hook settings.
type Option func(*StackdriverHook) error

// Levels is an option that sets the logrus levels that the StackdriverHook
// will create log entries for.
func Levels(levels ...logrus.Level) Option {
	return func(sh *StackdriverHook) error {
		sh.levels = levels
		return nil
	}
}

// ProjectID is an option that sets the project ID which is needed for the log
// name.
func ProjectID(projectID string) Option {
	return func(sh *StackdriverHook) error {
		sh.projectID = projectID
		return nil
	}
}

// EntriesService is an option that sets the Google API entry service to use
// with Stackdriver.
func EntriesService(service *logging.EntriesService) Option {
	return func(sh *StackdriverHook) error {
		sh.service = service
		return nil
	}
}

// LoggingService is an option that sets the Google API logging service to use.
func LoggingService(service *logging.Service) Option {
	return func(sh *StackdriverHook) error {
		sh.service = service.Entries
		return nil
	}
}

// ErrorService is an option that sets the Google API error reporting service to use.
func ErrorService(errorService *errorReporting.Service) Option {
	return func(sh *StackdriverHook) error {
		sh.errorService = errorService
		return nil
	}
}

// HTTPClient is an option that sets the http.Client to be used when creating
// the Stackdriver service.
func HTTPClient(client *http.Client) Option {
	return func(sh *StackdriverHook) error {
		// create logging service
		l, err := logging.New(client)
		if err != nil {
			return err
		}
		// create error reporting service
		e, err := errorReporting.New(client)
		if err != nil {
			return err
		} else {
			ErrorService(e)
		}

		return LoggingService(l)(sh)
	}
}

// MonitoredResource is an option that sets the monitored resource to send with
// each log entry.
func MonitoredResource(resource *logging.MonitoredResource) Option {
	return func(sh *StackdriverHook) error {
		sh.resource = resource
		return nil
	}
}

// Resource is an option that sets the resource information to send with each
// log entry.
//
// Please see https://cloud.google.com/logging/docs/api/v2/resource-list for
// the list of labels required per ResType.
func Resource(typ ResType, labels map[string]string) Option {
	return func(sh *StackdriverHook) error {
		return MonitoredResource(&logging.MonitoredResource{
			Type:   string(typ),
			Labels: labels,
		})(sh)
	}
}

// LogName is an option that sets the log name to send with each log entry.
//
// Log names are specified as "projects/{projectID}/logs/{logName}"
// if the projectID is set. Otherwise, it's just "{logName}"
func LogName(name string) Option {
	return func(sh *StackdriverHook) error {
		if sh.projectID == "" {
			sh.logName = name
		} else {
			sh.logName = fmt.Sprintf("projects/%s/logs/%s", sh.projectID, name)
		}
		return nil
	}
}

// ErrorReportingLogName is an option that sets the log name to send
// with each error message for error reporting.
// Only used when ErrorReportingService has been set.
func ErrorReportingLogName(name string) Option {
	return func(sh *StackdriverHook) error {
		sh.errorReportingLogName = name
		return nil
	}
}

// Labels is an option that sets the labels to send with each log entry.
func Labels(labels map[string]string) Option {
	return func(sh *StackdriverHook) error {
		sh.labels = labels
		return nil
	}
}

// PartialSuccess is an option that toggles whether or not to write partial log
// entries.
func PartialSuccess(enabled bool) Option {
	return func(sh *StackdriverHook) error {
		sh.partialSuccess = enabled
		return nil
	}
}

// ErrorReportingService is an option that defines the name of the service
// being tracked for Stackdriver error reporting.
// See:
// https://cloud.google.com/error-reporting/docs/formatting-error-messages
func ErrorReportingService(service string) Option {
	return func(sh *StackdriverHook) error {
		sh.errorReportingServiceName = service
		return nil
	}
}

// requiredScopes are the oauth2 scopes required for stackdriver logging.
var requiredScopes = []string{
	logging.CloudPlatformScope,
}

// GoogleServiceAccountCredentialsJSON is an option that creates the
// Stackdriver logging service using the supplied Google service account
// credentials.
//
// Google Service Account credentials can be downloaded from the Google Cloud
// console: https://console.cloud.google.com/iam-admin/serviceaccounts/
func GoogleServiceAccountCredentialsJSON(buf []byte) Option {
	return func(sh *StackdriverHook) error {
		var err error

		// load credentials
		gsa, err := gserviceaccount.FromJSON(buf)
		if err != nil {
			return err
		}

		// check project id
		if gsa.ProjectID == "" {
			return errors.New("google service account credentials missing project_id")
		}

		// set project id
		err = ProjectID(gsa.ProjectID)(sh)
		if err != nil {
			return err
		}

		// set resource type
		err = Resource(ResTypeProject, map[string]string{
			"project_id": gsa.ProjectID,
		})(sh)
		if err != nil {
			return err
		}

		// create token source
		ts, err := gsa.TokenSource(nil, requiredScopes...)
		if err != nil {
			return err
		}

		// set client
		return HTTPClient(&http.Client{
			Transport: &oauth2.Transport{
				Source: oauth2.ReuseTokenSource(nil, ts),
			},
		})(sh)
	}
}

// GoogleServiceAccountCredentialsFile is an option that loads Google Service
// Account credentials for use with the StackdriverHook from the specified
// file.
//
// Google Service Account credentials can be downloaded from the Google Cloud
// console: https://console.cloud.google.com/iam-admin/serviceaccounts/
func GoogleServiceAccountCredentialsFile(path string) Option {
	return func(sh *StackdriverHook) error {
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		return GoogleServiceAccountCredentialsJSON(buf)(sh)
	}
}

// GoogleComputeCredentials is an option that loads the Google Service Account
// credentials from the GCE metadata associated with the GCE compute instance.
// If serviceAccount is empty, then the default service account credentials
// associated with the GCE instance will be used.
func GoogleComputeCredentials(serviceAccount string) Option {
	return func(sh *StackdriverHook) error {
		var err error

		// get compute metadata scopes associated with the service account
		scopes, err := metadata.Scopes(serviceAccount)
		if err != nil {
			return err
		}

		// check if all the necessary scopes are provided
		for _, s := range requiredScopes {
			if !sliceContains(scopes, s) {
				// NOTE: if you are seeing this error, you probably need to
				// recreate your compute instance with the correct scope
				//
				// as of August 2016, there is not a way to add a scope to an
				// existing compute instance
				return fmt.Errorf("missing required scope %s in compute metadata", s)
			}
		}

		return HTTPClient(&http.Client{
			Transport: &oauth2.Transport{
				Source: google.ComputeTokenSource(serviceAccount),
			},
		})(sh)
	}
}

// sliceContains returns true if haystack contains needle.
func sliceContains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}

	return false
}

func GoogleLoggingAgent() Option {
	return func(sh *StackdriverHook) error {
		var err error
		// set agent client. It expects that the forward input fluentd plugin
		// is properly configured by the Google logging agent, which is by default.
		// See more at:
		// https://cloud.google.com/error-reporting/docs/setup/ec2
		sh.agentClient, err = fluent.New(fluent.Config{
			AsyncConnect: true,
			MaxRetry:     -1,
		})
		if err != nil {
			return fmt.Errorf("could not find fluentd agent on 127.0.0.1:24224")
		}
		return nil
	}
}
