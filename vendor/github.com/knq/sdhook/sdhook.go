// Package sdhook provides a logrus compatible logging hook for Google
// Stackdriver logging.
package sdhook

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/facebookgo/stack"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/sirupsen/logrus"

	errorReporting "google.golang.org/api/clouderrorreporting/v1beta1"
	logging "google.golang.org/api/logging/v2"
)

const (
	// DefaultName is the default name passed to LogName when using service
	// account credentials.
	DefaultName = "default"
)

// StackdriverHook provides a logrus hook to Google Stackdriver logging.
type StackdriverHook struct {
	// levels are the levels that logrus will hook to.
	levels []logrus.Level

	// projectID is the projectID
	projectID string

	// service is the logging service.
	service *logging.EntriesService

	// service is the error reporting service.
	errorService *errorReporting.Service

	// resource is the monitored resource.
	resource *logging.MonitoredResource

	// logName is the name of the log.
	logName string

	// labels are the labels to send with each log entry.
	labels map[string]string

	// partialSuccess allows partial writes of log entries if there is a badly
	// formatted log.
	partialSuccess bool

	// agentClient defines the fluentd logger object that can send data to
	// to the Google logging agent.
	agentClient *fluent.Fluent

	// errorReportingServiceName defines the value of the field <service>,
	// required for a valid error reporting payload. If this value is set,
	// messages where level/severity is higher than or equal to "error" will
	// be sent to Stackdriver error reporting.
	// See more at:
	// https://cloud.google.com/error-reporting/docs/formatting-error-messages
	errorReportingServiceName string

	// errorReportingLogName is the name of the log for error reporting.
	// It must contain the string "error"
	// If not given, the string "<logName>_error" is used.
	errorReportingLogName string
}

// New creates a StackdriverHook using the provided options that is suitible
// for using with logrus for logging to Google Stackdriver.
func New(opts ...Option) (*StackdriverHook, error) {
	var err error

	sh := &StackdriverHook{
		levels: logrus.AllLevels,
	}

	// apply opts
	for _, o := range opts {
		err = o(sh)
		if err != nil {
			return nil, err
		}
	}

	// check service, resource, logName set
	if sh.service == nil && sh.agentClient == nil {
		return nil, errors.New("no stackdriver service was provided")
	}
	if sh.resource == nil && sh.agentClient == nil {
		return nil, errors.New("the monitored resource was not provided")
	}
	if sh.projectID == "" && sh.agentClient == nil {
		return nil, errors.New("the project id was not provided")
	}

	// set default project name
	if sh.logName == "" {
		err = LogName(DefaultName)(sh)
		if err != nil {
			return nil, err
		}
	}

	// If error reporting log name not set, set it to log name
	// plus string suffix
	if sh.errorReportingLogName == "" {
		sh.errorReportingLogName = sh.logName + "_errors"
	}

	return sh, nil
}

func isError(entry *logrus.Entry) bool {
	if entry != nil {
		switch entry.Level {
		case logrus.ErrorLevel:
			return true
		case logrus.FatalLevel:
			return true
		case logrus.PanicLevel:
			return true
		}
	}
	return false
}

// Levels returns the logrus levels that this hook is applied to. This can be
// set using the Levels Option.
func (sh *StackdriverHook) Levels() []logrus.Level {
	return sh.levels
}

// Fire writes the message to the Stackdriver entry service.
func (sh *StackdriverHook) Fire(entry *logrus.Entry) error {
	go func(entry *logrus.Entry) {
		var httpReq *logging.HttpRequest

		// convert entry data to labels
		labels := make(map[string]string, len(entry.Data))
		for k, v := range entry.Data {
			switch x := v.(type) {
			case string:
				labels[k] = x

			case *http.Request:
				httpReq = &logging.HttpRequest{
					Referer:       x.Referer(),
					RemoteIp:      x.RemoteAddr,
					RequestMethod: x.Method,
					RequestUrl:    x.URL.String(),
					UserAgent:     x.UserAgent(),
				}

			case *logging.HttpRequest:
				httpReq = x

			default:
				labels[k] = fmt.Sprintf("%v", v)
			}
		}

		// write log entry
		if sh.agentClient != nil {
			sh.sendLogMessageViaAgent(entry, labels, httpReq)
		} else {
			sh.sendLogMessageViaAPI(entry, labels, httpReq)
		}
	}(sh.copyEntry(entry))

	return nil
}

func (sh *StackdriverHook) copyEntry(entry *logrus.Entry) *logrus.Entry {
	e := *entry
	e.Data = make(logrus.Fields, len(entry.Data))
	for k, v := range entry.Data {
		e.Data[k] = v
	}
	return &e
}

func (sh *StackdriverHook) sendLogMessageViaAgent(entry *logrus.Entry, labels map[string]string, httpReq *logging.HttpRequest) {
	// The log entry payload schema is defined by the Google fluentd
	// logging agent. See more at:
	// https://github.com/GoogleCloudPlatform/fluent-plugin-google-cloud
	logEntry := map[string]interface{}{
		"severity":         strings.ToUpper(entry.Level.String()),
		"timestampSeconds": strconv.FormatInt(entry.Time.Unix(), 10),
		"timestampNanos":   strconv.FormatInt(entry.Time.UnixNano()-entry.Time.Unix()*1000000000, 10),
		"message":          entry.Message,
	}
	for k, v := range labels {
		logEntry[k] = v
	}
	if httpReq != nil {
		logEntry["httpRequest"] = httpReq
	}
	// The error reporting payload JSON schema is defined in:
	// https://cloud.google.com/error-reporting/docs/formatting-error-messages
	// Which reflects the structure of the ErrorEvent type in:
	// https://godoc.org/google.golang.org/api/clouderrorreporting/v1beta1
	if sh.errorReportingServiceName != "" && isError(entry) {
		errorEvent := sh.buildErrorReportingEvent(entry, labels, httpReq)
		errorStructPayload, err := json.Marshal(errorEvent)
		if err != nil {
			log.Printf("error marshaling error reporting data: %s", err.Error())
		}
		var errorJSONPayload map[string]interface{}
		err = json.Unmarshal(errorStructPayload, &errorJSONPayload)
		if err != nil {
			log.Printf("error parsing error reporting data: %s", err.Error())
		}
		for k, v := range logEntry {
			errorJSONPayload[k] = v
		}
		if err := sh.agentClient.Post(sh.errorReportingLogName, errorJSONPayload); err != nil {
			log.Printf("error posting error reporting entries to logging agent: %s", err.Error())
		}
	} else {
		if err := sh.agentClient.Post(sh.logName, logEntry); err != nil {
			log.Printf("error posting log entries to logging agent: %s", err.Error())
		}
	}
}

func (sh *StackdriverHook) sendLogMessageViaAPI(entry *logrus.Entry, labels map[string]string, httpReq *logging.HttpRequest) {
	if sh.errorReportingServiceName != "" && isError(entry) {
		errorEvent := sh.buildErrorReportingEvent(entry, labels, httpReq)
		sh.errorService.Projects.Events.Report(sh.projectID, &errorEvent)
	} else {
		logName := sh.logName
		if sh.errorReportingLogName != "" && isError(entry) {
			logName = sh.errorReportingLogName
		}
		_, _ = sh.service.Write(&logging.WriteLogEntriesRequest{
			LogName:        logName,
			Resource:       sh.resource,
			Labels:         sh.labels,
			PartialSuccess: sh.partialSuccess,
			Entries: []*logging.LogEntry{
				{
					Severity:    strings.ToUpper(entry.Level.String()),
					Timestamp:   entry.Time.Format(time.RFC3339),
					TextPayload: entry.Message,
					Labels:      labels,
					HttpRequest: httpReq,
				},
			},
		}).Do()
	}
}

func (sh *StackdriverHook) buildErrorReportingEvent(entry *logrus.Entry, labels map[string]string, httpReq *logging.HttpRequest) errorReporting.ReportedErrorEvent {
	errorEvent := errorReporting.ReportedErrorEvent{
		EventTime: entry.Time.Format(time.RFC3339),
		Message:   entry.Message,
		ServiceContext: &errorReporting.ServiceContext{
			Service: sh.errorReportingServiceName,
			Version: labels["version"],
		},
		Context: &errorReporting.ErrorContext{
			User: labels["user"],
		},
	}
	// Assumes that caller stack frame information of type
	// github.com/facebookgo/stack.Frame has been added.
	// Possibly via a library like github.com/Gurpartap/logrus-stack
	if entry.Data["caller"] != nil {
		caller := entry.Data["caller"].(stack.Frame)
		errorEvent.Context.ReportLocation = &errorReporting.SourceLocation{
			FilePath:     caller.File,
			FunctionName: caller.Name,
			LineNumber:   int64(caller.Line),
		}
	}
	if httpReq != nil {
		errRepHttpRequest := &errorReporting.HttpRequestContext{
			Method:    httpReq.RequestMethod,
			Referrer:  httpReq.Referer,
			RemoteIp:  httpReq.RemoteIp,
			Url:       httpReq.RequestUrl,
			UserAgent: httpReq.UserAgent,
		}
		errorEvent.Context.HttpRequest = errRepHttpRequest
	}
	return errorEvent
}
