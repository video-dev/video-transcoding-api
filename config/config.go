package config

import "github.com/NYTimes/gizmo/config"

// Config is a struct to contain all the needed configuration for the
// Transcoding API.
type Config struct {
	*config.Server
	EncodingComUserID  string `envconfig:"ENCODINGCOM_USER_ID"`
	EncodingComUserKey string `envconfig:"ENCODINGCOM_USER_KEY"`

	// S3Bucket is the name of the bucket on S3 for storing input files and
	// also the output files.
	S3Bucket string `envconfig:"TRANSCODING_S3_BUCKET"`
}
