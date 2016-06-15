package config

import (
	"github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
)

type defaultLoader interface {
	loadDefaults()
}

// Config is a struct to contain all the needed configuration for the
// Transcoding API.
type Config struct {
	Server             *server.Config
	SwaggerManifest    string `envconfig:"SWAGGER_MANIFEST_PATH"`
	Redis              *Redis
	EncodingCom        *EncodingCom
	ElasticTranscoder  *ElasticTranscoder
	ElementalConductor *ElementalConductor
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

	RedisAddr          string `envconfig:"REDIS_ADDR"`
	Password           string `envconfig:"REDIS_PASSWORD"`
	PoolSize           int    `envconfig:"REDIS_POOL_SIZE"`
	PoolTimeout        int    `envconfig:"REDIS_POOL_TIMEOUT_SECONDS"`
	IdleTimeout        int    `envconfig:"REDIS_IDLE_TIMEOUT_SECONDS"`
	IdleCheckFrequency int    `envconfig:"REDIS_IDLE_CHECK_FREQUENCY_SECONDS"`
}

func (c *Redis) loadDefaults() {
	if c.RedisAddr == "" {
		c.RedisAddr = "127.0.0.1:6379"
	}
}

// EncodingCom represents the set of configurations for the Encoding.com
// provider.
type EncodingCom struct {
	UserID         string `envconfig:"ENCODINGCOM_USER_ID"`
	UserKey        string `envconfig:"ENCODINGCOM_USER_KEY"`
	Destination    string `envconfig:"ENCODINGCOM_DESTINATION"`
	Region         string `envconfig:"ENCODINGCOM_REGION"`
	StatusEndpoint string `envconfig:"ENCODINGCOM_STATUS_ENDPOINT"`
}

func (c *EncodingCom) loadDefaults() {
	if c.StatusEndpoint == "" {
		c.StatusEndpoint = "http://status.encoding.com"
	}
}

// ElasticTranscoder represents the set of configurations for the Elastic
// Transcoder provider.
type ElasticTranscoder struct {
	AccessKeyID     string `envconfig:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `envconfig:"AWS_SECRET_ACCESS_KEY"`
	Region          string `envconfig:"AWS_REGION"`
	PipelineID      string `envconfig:"ELASTICTRANSCODER_PIPELINE_ID"`
}

// ElementalConductor represents the set of configurations for the Elemental
// Conductor provider.
type ElementalConductor struct {
	Host            string `envconfig:"ELEMENTALCONDUCTOR_HOST"`
	UserLogin       string `envconfig:"ELEMENTALCONDUCTOR_USER_LOGIN"`
	APIKey          string `envconfig:"ELEMENTALCONDUCTOR_API_KEY"`
	AuthExpires     int    `envconfig:"ELEMENTALCONDUCTOR_AUTH_EXPIRES"`
	AccessKeyID     string `envconfig:"ELEMENTALCONDUCTOR_AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `envconfig:"ELEMENTALCONDUCTOR_AWS_SECRET_ACCESS_KEY"`
	Destination     string `envconfig:"ELEMENTALCONDUCTOR_DESTINATION"`
}

// LoadConfig loads the configuration of the API using environment variables.
func LoadConfig() *Config {
	cfg := Config{
		Redis:              new(Redis),
		EncodingCom:        new(EncodingCom),
		ElasticTranscoder:  new(ElasticTranscoder),
		ElementalConductor: new(ElementalConductor),
		Server:             new(server.Config),
	}
	config.LoadEnvConfig(&cfg)
	loadFromEnvAndDefaults(cfg.Redis, cfg.EncodingCom, cfg.ElasticTranscoder, cfg.ElementalConductor, cfg.Server)
	return &cfg
}

func loadFromEnvAndDefaults(cfgs ...interface{}) {
	for _, cfg := range cfgs {
		config.LoadEnvConfig(cfg)
		if dLoader, ok := cfg.(defaultLoader); ok {
			dLoader.loadDefaults()
		}
	}
}
