package config

import (
	"github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
	"github.com/NYTimes/video-transcoding-api/db/redis/storage"
)

// Config is a struct to contain all the needed configuration for the
// Transcoding API.
type Config struct {
	Server                 *server.Config
	SwaggerManifest        string `envconfig:"SWAGGER_MANIFEST_PATH"`
	DefaultSegmentDuration uint   `envconfig:"DEFAULT_SEGMENT_DURATION" default:"5"`
	Redis                  *storage.Config
	EncodingCom            *EncodingCom
	ElasticTranscoder      *ElasticTranscoder
	ElementalConductor     *ElementalConductor
	Zencoder               *Zencoder
	Bitmovin               *Bitmovin
}

// EncodingCom represents the set of configurations for the Encoding.com
// provider.
type EncodingCom struct {
	UserID         string `envconfig:"ENCODINGCOM_USER_ID"`
	UserKey        string `envconfig:"ENCODINGCOM_USER_KEY"`
	Destination    string `envconfig:"ENCODINGCOM_DESTINATION"`
	Region         string `envconfig:"ENCODINGCOM_REGION"`
	StatusEndpoint string `envconfig:"ENCODINGCOM_STATUS_ENDPOINT" default:"http://status.encoding.com"`
}

// Zencoder represents the set of configurations for the Zencoder
// provider.
type Zencoder struct {
	APIKey      string `envconfig:"ZENCODER_API_KEY"`
	Destination string `envconfig:"ZENCODER_DESTINATION"`
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

// Bitmovin represents the set of configurations for the Bitmovin
// provider.
type Bitmovin struct {
	APIKey   string `envconfig:"BITMOVIN_API_KEY"`
	Endpoint string `envconfig:"BITMOVIN_ENDPOINT" default:"https://api.bitmovin.com/v1/"`
	Timeout  uint   `envconfig:"BITMOVIN_TIMEOUT" default:"5"`
}

// LoadConfig loads the configuration of the API using environment variables.
func LoadConfig() *Config {
	cfg := Config{
		Redis:              new(storage.Config),
		EncodingCom:        new(EncodingCom),
		ElasticTranscoder:  new(ElasticTranscoder),
		ElementalConductor: new(ElementalConductor),
		Bitmovin:           new(Bitmovin),
		Server:             new(server.Config),
	}
	config.LoadEnvConfig(&cfg)
	loadFromEnv(cfg.Redis, cfg.EncodingCom, cfg.ElasticTranscoder, cfg.ElementalConductor, cfg.Bitmovin, cfg.Server)
	return &cfg
}

func loadFromEnv(cfgs ...interface{}) {
	for _, cfg := range cfgs {
		config.LoadEnvConfig(cfg)
	}
}
