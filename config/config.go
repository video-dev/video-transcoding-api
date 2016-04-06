package config

import "github.com/NYTimes/gizmo/config"

type defaultLoader interface {
	loadDefaults()
}

// Config is a struct to contain all the needed configuration for the
// Transcoding API.
type Config struct {
	*config.Server
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

// LoadConfig loads the configuration of the API using the provided file and
// environment variables. It will override settings defined in the file with
// the value of environment variables.
//
// Provide an empty file name for loading configuration exclusively from the
// environemtn.
func LoadConfig(fileName string) *Config {
	var cfg Config
	if fileName != "" {
		config.LoadJSONFile(fileName, &cfg)
	}
	config.LoadEnvConfig(&cfg)
	if cfg.Redis == nil {
		cfg.Redis = new(Redis)
	}
	if cfg.EncodingCom == nil {
		cfg.EncodingCom = new(EncodingCom)
	}
	if cfg.ElasticTranscoder == nil {
		cfg.ElasticTranscoder = new(ElasticTranscoder)
	}
	if cfg.ElementalConductor == nil {
		cfg.ElementalConductor = new(ElementalConductor)
	}
	loadFromEnvAndDefaults(cfg.Redis, cfg.EncodingCom, cfg.ElasticTranscoder, cfg.ElementalConductor)
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
