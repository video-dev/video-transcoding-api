package config

import "github.com/NYTimes/gizmo/config"

// Config is a struct to contain all the needed
// configuration for our JSONService
type Config struct {
	*config.Server
	EncodingComUserID  string `envconfig:"ENCODINGCOM_USER_ID"`
	EncodingComUserKey string `envconfig:"ENCODINGCOM_USER_KEY"`
}
