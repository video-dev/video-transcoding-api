package config

import (
	"github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
	"github.com/NYTimes/video-transcoding-api/db/redis/storage"
	"github.com/fsouza/gizmo-stackdriver-logging"
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
	Hybrik                 *Hybrik
	Zencoder               *Zencoder
	Bitmovin               *Bitmovin
	Log                    *logging.Config
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
	APIKey           string `envconfig:"BITMOVIN_API_KEY"`
	Endpoint         string `envconfig:"BITMOVIN_ENDPOINT" default:"https://api.bitmovin.com/v1/"`
	Timeout          uint   `envconfig:"BITMOVIN_TIMEOUT" default:"5"`
	AccessKeyID      string `envconfig:"BITMOVIN_AWS_ACCESS_KEY_ID"`
	SecretAccessKey  string `envconfig:"BITMOVIN_AWS_SECRET_ACCESS_KEY"`
	Destination      string `envconfig:"BITMOVIN_DESTINATION"`
	AWSStorageRegion string `envconfig:"BITMOVIN_AWS_STORAGE_REGION" default:"US_EAST_1"`
	EncodingRegion   string `envconfig:"BITMOVIN_ENCODING_REGION" default:"AWS_US_EAST_1"`
	EncodingVersion  string `envconfig:"BITMOVIN_ENCODING_VERSION" default:"STABLE"`
}

// Hybrik represents the set of configurations for the Hybrik
// provider.
type Hybrik struct {
	URL            string `envconfig:"HYBRIK_URL"`
	ComplianceDate string `envconfig:"HYBRIK_COMPLIANCE_DATE" default:"20170601"`
	OAPIKey        string `envconfig:"HYBRIK_OAPI_KEY"`
	OAPISecret     string `envconfig:"HYBRIK_OAPI_SECRET"`
	AuthKey        string `envconfig:"HYBRIK_AUTH_KEY"`
	AuthSecret     string `envconfig:"HYBRIK_AUTH_SECRET"`
	Destination    string `envconfig:"HYBRIK_DESTINATION"`
	PresetPath     string `envconfig:"HYBRIK_PRESET_PATH" default:"transcoding-api-presets"`
}

// LoadConfig loads the configuration of the API using environment variables.
func LoadConfig() *Config {
	var cfg Config
	config.LoadEnvConfig(&cfg)
	return &cfg
}
