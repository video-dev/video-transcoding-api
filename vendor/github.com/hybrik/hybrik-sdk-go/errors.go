package hybrik

import "errors"

var (
	// ErrNoOAPIKey is the error for zero value config OAPIKey
	ErrNoOAPIKey = errors.New("no oapi key provided")
	// ErrNoOAPISecret is the error for zero value config OAPISecret
	ErrNoOAPISecret = errors.New("no oapi secret provided")
	// ErrNoAuthKey is the error for zero value config AuthKey
	ErrNoAuthKey = errors.New("no auth key provided")
	// ErrNoAuthSecret is the error for zero value config AuthSecret
	ErrNoAuthSecret = errors.New("no auth secret provided")
	// ErrNoAPIURL is the error for zero value config API URL
	ErrNoAPIURL = errors.New("no api url provided")
	// ErrNoComplianceDate is the error for zero value config ComplianceDate
	ErrNoComplianceDate = errors.New("no compliance date provided")

	// ErrInvalidComplianceDate is the error for an invalid compliance date provided
	ErrInvalidComplianceDate = errors.New("provided compliance date is invalid. Expected format: 'YYYYMMDD'")

	// ErrInvalidURL is the error for an invalid URL provided
	ErrInvalidURL = errors.New("invalid api url provided")
)
