package config

import "github.com/NYTimes/gizmo/config"

// Config is a struct to contain all the needed configuration for the
// Transcoding API.
type Config struct {
	*config.Server
	*config.S3

	EncodingComUserID  string `envconfig:"ENCODINGCOM_USER_ID"`
	EncodingComUserKey string `envconfig:"ENCODINGCOM_USER_KEY"`
}
