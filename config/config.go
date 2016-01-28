package config

import "github.com/NYTimes/gizmo/config"

// Config is a struct to contain all the needed configuration for the
// Transcoding API.
type Config struct {
	*config.Server
	*config.S3
	Redis
	EncodingCom
}

// Redis represents the Redis configuration. RedisAddr and SentinelAddrs
// configs are exclusive, and the API will prefer to use SentinelAddrs when
// both are defined.
type Redis struct {
	// Comma-separated list of sentinel servers.
	//
	// Example: 10.10.10.10:6379,10.10.10.1:6379,10.10.10.2:6379.
	SentinelAddrs      string `envconfig:"SENTINEL_ADDRS"`
	SentinelMasterName string `envconfig:"SENTINEL_MASTER_NAME"`

	RedisAddr string `envconfig:"REDIS_ADDR"`
	Password  string `envconfig:"REDIS_PASSWORD"`
}

// EncodingCom represents the set of configurations for the Encoding.com
// provider.
type EncodingCom struct {
	UserID  string `envconfig:"ENCODINGCOM_USER_ID"`
	UserKey string `envconfig:"ENCODINGCOM_USER_KEY"`
}
