package config

import (
	"os"
	"testing"

	"github.com/NYTimes/gizmo/server"
	"github.com/NYTimes/video-transcoding-api/db/redis/storage"
	"github.com/fsouza/gizmo-stackdriver-logging"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestLoadConfigFromEnv(t *testing.T) {
	os.Clearenv()
	accessLog := "/var/log/transcoding-api-access.log"
	setEnvs(map[string]string{
		"SENTINEL_ADDRS":                           "10.10.10.10:26379,10.10.10.11:26379,10.10.10.12:26379",
		"SENTINEL_MASTER_NAME":                     "super-master",
		"REDIS_ADDR":                               "localhost:6379",
		"REDIS_PASSWORD":                           "super-secret",
		"REDIS_POOL_SIZE":                          "100",
		"REDIS_POOL_TIMEOUT_SECONDS":               "10",
		"ENCODINGCOM_USER_ID":                      "myuser",
		"ENCODINGCOM_USER_KEY":                     "secret-key",
		"ENCODINGCOM_DESTINATION":                  "https://safe-stuff",
		"ENCODINGCOM_STATUS_ENDPOINT":              "https://safe-status",
		"ENCODINGCOM_REGION":                       "sa-east-1",
		"AWS_ACCESS_KEY_ID":                        "AKIANOTREALLY",
		"AWS_SECRET_ACCESS_KEY":                    "secret-key",
		"AWS_REGION":                               "us-east-1",
		"ELASTICTRANSCODER_PIPELINE_ID":            "mypipeline",
		"ELEMENTALCONDUCTOR_HOST":                  "elemental-server",
		"ELEMENTALCONDUCTOR_USER_LOGIN":            "myuser",
		"ELEMENTALCONDUCTOR_API_KEY":               "secret-key",
		"ELEMENTALCONDUCTOR_AUTH_EXPIRES":          "30",
		"ELEMENTALCONDUCTOR_AWS_ACCESS_KEY_ID":     "AKIANOTREALLY",
		"ELEMENTALCONDUCTOR_AWS_SECRET_ACCESS_KEY": "secret-key",
		"ELEMENTALCONDUCTOR_DESTINATION":           "https://safe-stuff",
		"BITMOVIN_API_KEY":                         "secret-key",
		"BITMOVIN_ENDPOINT":                        "bitmovin",
		"BITMOVIN_TIMEOUT":                         "3",
		"BITMOVIN_AWS_ACCESS_KEY_ID":               "AKIANOTREALLY",
		"BITMOVIN_AWS_SECRET_ACCESS_KEY":           "secret-key",
		"BITMOVIN_DESTINATION":                     "https://safe-stuff",
		"BITMOVIN_AWS_STORAGE_REGION":              "US_WEST_1",
		"BITMOVIN_ENCODING_REGION":                 "GOOGLE_EUROPE_WEST_1",
		"BITMOVIN_ENCODING_VERSION":                "notstable",
		"SWAGGER_MANIFEST_PATH":                    "/opt/video-transcoding-api-swagger.json",
		"HTTP_ACCESS_LOG":                          accessLog,
		"HTTP_PORT":                                "8080",
		"DEFAULT_SEGMENT_DURATION":                 "3",
		"LOGGING_LEVEL":                            "debug",
	})
	cfg := LoadConfig()
	expectedCfg := Config{
		SwaggerManifest:        "/opt/video-transcoding-api-swagger.json",
		DefaultSegmentDuration: 3,
		Redis: &storage.Config{
			SentinelAddrs:      "10.10.10.10:26379,10.10.10.11:26379,10.10.10.12:26379",
			SentinelMasterName: "super-master",
			RedisAddr:          "localhost:6379",
			Password:           "super-secret",
			PoolSize:           100,
			PoolTimeout:        10,
		},
		EncodingCom: &EncodingCom{
			UserID:         "myuser",
			UserKey:        "secret-key",
			Destination:    "https://safe-stuff",
			StatusEndpoint: "https://safe-status",
			Region:         "sa-east-1",
		},
		Hybrik: &Hybrik{
			ComplianceDate: "20170601",
			PresetPath:     "transcoding-api-presets",
		},
		Zencoder: &Zencoder{},
		ElasticTranscoder: &ElasticTranscoder{
			AccessKeyID:     "AKIANOTREALLY",
			SecretAccessKey: "secret-key",
			Region:          "us-east-1",
			PipelineID:      "mypipeline",
		},
		ElementalConductor: &ElementalConductor{
			Host:            "elemental-server",
			UserLogin:       "myuser",
			APIKey:          "secret-key",
			AuthExpires:     30,
			AccessKeyID:     "AKIANOTREALLY",
			SecretAccessKey: "secret-key",
			Destination:     "https://safe-stuff",
		},
		Bitmovin: &Bitmovin{
			APIKey:           "secret-key",
			Endpoint:         "bitmovin",
			Timeout:          3,
			AccessKeyID:      "AKIANOTREALLY",
			SecretAccessKey:  "secret-key",
			AWSStorageRegion: "US_WEST_1",
			Destination:      "https://safe-stuff",
			EncodingRegion:   "GOOGLE_EUROPE_WEST_1",
			EncodingVersion:  "notstable",
		},
		Server: &server.Config{
			HTTPPort:      8080,
			HTTPAccessLog: &accessLog,
		},
		Log: &logging.Config{
			StackDriverErrorLogName: "error_log",

			Level: "debug",
		},
	}
	diff := cmp.Diff(*cfg, expectedCfg, cmpopts.IgnoreUnexported(server.Config{}))
	if diff != "" {
		t.Errorf("LoadConfig(): wrong config\nWant %#v\nGot %#v\nDiff: %v", expectedCfg, *cfg, diff)
	}
}

func TestLoadConfigFromEnvWithDefaults(t *testing.T) {
	os.Clearenv()
	accessLog := "/var/log/transcoding-api-access.log"
	setEnvs(map[string]string{
		"SENTINEL_ADDRS":                           "10.10.10.10:26379,10.10.10.11:26379,10.10.10.12:26379",
		"SENTINEL_MASTER_NAME":                     "super-master",
		"REDIS_PASSWORD":                           "super-secret",
		"REDIS_POOL_SIZE":                          "100",
		"REDIS_POOL_TIMEOUT_SECONDS":               "10",
		"REDIS_IDLE_TIMEOUT_SECONDS":               "30",
		"REDIS_IDLE_CHECK_FREQUENCY_SECONDS":       "20",
		"ENCODINGCOM_USER_ID":                      "myuser",
		"ENCODINGCOM_USER_KEY":                     "secret-key",
		"ENCODINGCOM_DESTINATION":                  "https://safe-stuff",
		"AWS_ACCESS_KEY_ID":                        "AKIANOTREALLY",
		"AWS_SECRET_ACCESS_KEY":                    "secret-key",
		"AWS_REGION":                               "us-east-1",
		"ELASTICTRANSCODER_PIPELINE_ID":            "mypipeline",
		"ELEMENTALCONDUCTOR_HOST":                  "elemental-server",
		"ELEMENTALCONDUCTOR_USER_LOGIN":            "myuser",
		"ELEMENTALCONDUCTOR_API_KEY":               "secret-key",
		"ELEMENTALCONDUCTOR_AUTH_EXPIRES":          "30",
		"ELEMENTALCONDUCTOR_AWS_ACCESS_KEY_ID":     "AKIANOTREALLY",
		"ELEMENTALCONDUCTOR_AWS_SECRET_ACCESS_KEY": "secret-key",
		"ELEMENTALCONDUCTOR_DESTINATION":           "https://safe-stuff",
		"BITMOVIN_API_KEY":                         "secret-key",
		"BITMOVIN_AWS_ACCESS_KEY_ID":               "AKIANOTREALLY",
		"BITMOVIN_AWS_SECRET_ACCESS_KEY":           "secret-key",
		"BITMOVIN_DESTINATION":                     "https://safe-stuff",
		"SWAGGER_MANIFEST_PATH":                    "/opt/video-transcoding-api-swagger.json",
		"HTTP_ACCESS_LOG":                          accessLog,
		"HTTP_PORT":                                "8080",
	})
	cfg := LoadConfig()
	expectedCfg := Config{
		SwaggerManifest:        "/opt/video-transcoding-api-swagger.json",
		DefaultSegmentDuration: 5,
		Redis: &storage.Config{
			SentinelAddrs:      "10.10.10.10:26379,10.10.10.11:26379,10.10.10.12:26379",
			SentinelMasterName: "super-master",
			RedisAddr:          "127.0.0.1:6379",
			Password:           "super-secret",
			PoolSize:           100,
			PoolTimeout:        10,
			IdleCheckFrequency: 20,
			IdleTimeout:        30,
		},
		EncodingCom: &EncodingCom{
			UserID:         "myuser",
			UserKey:        "secret-key",
			Destination:    "https://safe-stuff",
			StatusEndpoint: "http://status.encoding.com",
		},
		ElasticTranscoder: &ElasticTranscoder{
			AccessKeyID:     "AKIANOTREALLY",
			SecretAccessKey: "secret-key",
			Region:          "us-east-1",
			PipelineID:      "mypipeline",
		},
		ElementalConductor: &ElementalConductor{
			Host:            "elemental-server",
			UserLogin:       "myuser",
			APIKey:          "secret-key",
			AuthExpires:     30,
			AccessKeyID:     "AKIANOTREALLY",
			SecretAccessKey: "secret-key",
			Destination:     "https://safe-stuff",
		},
		Hybrik: &Hybrik{
			ComplianceDate: "20170601",
			PresetPath:     "transcoding-api-presets",
		},
		Zencoder: &Zencoder{},
		Bitmovin: &Bitmovin{
			APIKey:           "secret-key",
			Endpoint:         "https://api.bitmovin.com/v1/",
			Timeout:          5,
			AccessKeyID:      "AKIANOTREALLY",
			SecretAccessKey:  "secret-key",
			Destination:      "https://safe-stuff",
			AWSStorageRegion: "US_EAST_1",
			EncodingRegion:   "AWS_US_EAST_1",
			EncodingVersion:  "STABLE",
		},
		Server: &server.Config{
			HTTPPort:      8080,
			HTTPAccessLog: &accessLog,
		},
		Log: &logging.Config{
			Level: "info",

			StackDriverErrorLogName: "error_log",
		},
	}
	diff := cmp.Diff(*cfg, expectedCfg, cmpopts.IgnoreUnexported(server.Config{}))
	if diff != "" {
		t.Errorf("LoadConfig(): wrong config\nWant %#v\nGot %#v\nDiff: %v", expectedCfg, *cfg, diff)
	}
}

func setEnvs(envs map[string]string) {
	for k, v := range envs {
		os.Setenv(k, v)
	}
}
