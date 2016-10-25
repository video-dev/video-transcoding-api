package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/NYTimes/gizmo/server"
	"github.com/NYTimes/video-transcoding-api/db/redis/storage"
	"github.com/marzagao/envconfigfromfile"
)

func TestLoadConfigFromEnv(t *testing.T) {
	os.Clearenv()
	accessLog := "/var/log/transcoding-api-access.log"
	gcpCredsTestFilePath := "testdata/fake_gcp_creds.json"
	gcpCredsTestFileContents, _ := ioutil.ReadFile(gcpCredsTestFilePath)
	setEnvs(map[string]string{
		"SENTINEL_ADDRS":                           "10.10.10.10:26379,10.10.10.11:26379,10.10.10.12:26379",
		"SENTINEL_MASTER_NAME":                     "supermaster",
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
		"SWAGGER_MANIFEST_PATH":                    "/opt/video-transcoding-api-swagger.json",
		"HTTP_ACCESS_LOG":                          accessLog,
		"HTTP_PORT":                                "8080",
		"DEFAULT_SEGMENT_DURATION":                 "3",
		"GCP_CREDENTIALS_FILE":                     gcpCredsTestFilePath,
	})
	cfg := LoadConfig()
	expectedCfg := Config{
		SwaggerManifest:        "/opt/video-transcoding-api-swagger.json",
		DefaultSegmentDuration: 3,
		Redis: &storage.Config{
			SentinelAddrs:      "10.10.10.10:26379,10.10.10.11:26379,10.10.10.12:26379",
			SentinelMasterName: "supermaster",
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
		Server: &server.Config{
			HTTPPort:      8080,
			HTTPAccessLog: &accessLog,
		},
		GCPCredentials: &envconfigfromfile.EnvConfigFromFile{
			FilePath: gcpCredsTestFilePath,
			Value:    string(gcpCredsTestFileContents),
		},
	}
	if cfg.SwaggerManifest != expectedCfg.SwaggerManifest {
		t.Errorf("LoadConfig(): wrong swagger manifest. Want %q. Got %q", expectedCfg.SwaggerManifest, cfg.SwaggerManifest)
	}
	if cfg.DefaultSegmentDuration != expectedCfg.DefaultSegmentDuration {
		t.Errorf("LoadConfig(): wrong default segment duration. Want %q. Got %q", expectedCfg.DefaultSegmentDuration, cfg.DefaultSegmentDuration)
	}
	if !reflect.DeepEqual(*cfg.Redis, *expectedCfg.Redis) {
		t.Errorf("LoadConfig(): wrong Redis config returned. Want %#v. Got %#v.", *expectedCfg.Redis, *cfg.Redis)
	}
	if !reflect.DeepEqual(*cfg.EncodingCom, *expectedCfg.EncodingCom) {
		t.Errorf("LoadConfig(): wrong EncodingCom config returned. Want %#v. Got %#v.", *expectedCfg.EncodingCom, *cfg.EncodingCom)
	}
	if !reflect.DeepEqual(*cfg.ElasticTranscoder, *expectedCfg.ElasticTranscoder) {
		t.Errorf("LoadConfig(): wrong ElasticTranscoder config returned. Want %#v. Got %#v.", *expectedCfg.ElasticTranscoder, *cfg.ElasticTranscoder)
	}
	if !reflect.DeepEqual(*cfg.ElementalConductor, *expectedCfg.ElementalConductor) {
		t.Errorf("LoadConfig(): wrong Elemental Conductor config returned. Want %#v. Got %#v.", *expectedCfg.ElementalConductor, *cfg.ElementalConductor)
	}
	if !reflect.DeepEqual(*cfg.GCPCredentials, *expectedCfg.GCPCredentials) {
		t.Errorf("LoadConfig(): Wrong GCPCredentials returned. Want %#v. Got %#v.", *expectedCfg.GCPCredentials, *cfg.GCPCredentials)
	}
}

func TestLoadConfigFromEnvWithDefauts(t *testing.T) {
	os.Clearenv()
	accessLog := "/var/log/transcoding-api-access.log"
	setEnvs(map[string]string{
		"SENTINEL_ADDRS":                           "10.10.10.10:26379,10.10.10.11:26379,10.10.10.12:26379",
		"SENTINEL_MASTER_NAME":                     "supermaster",
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
			SentinelMasterName: "supermaster",
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
		Server: &server.Config{
			HTTPPort:      8080,
			HTTPAccessLog: &accessLog,
		},
	}
	if cfg.SwaggerManifest != expectedCfg.SwaggerManifest {
		t.Errorf("LoadConfig(): wrong swagger manifest. Want %q. Got %q", expectedCfg.SwaggerManifest, cfg.SwaggerManifest)
	}
	if cfg.DefaultSegmentDuration != expectedCfg.DefaultSegmentDuration {
		t.Errorf("LoadConfig(): wrong default segment duration. Want %q. Got %q", expectedCfg.DefaultSegmentDuration, cfg.DefaultSegmentDuration)
	}
	if !reflect.DeepEqual(*cfg.Redis, *expectedCfg.Redis) {
		t.Errorf("LoadConfig(): wrong Redis config returned. Want %#v. Got %#v.", *expectedCfg.Redis, *cfg.Redis)
	}
	if !reflect.DeepEqual(*cfg.EncodingCom, *expectedCfg.EncodingCom) {
		t.Errorf("LoadConfig(): wrong EncodingCom config returned. Want %#v. Got %#v.", *expectedCfg.EncodingCom, *cfg.EncodingCom)
	}
	if !reflect.DeepEqual(*cfg.ElasticTranscoder, *expectedCfg.ElasticTranscoder) {
		t.Errorf("LoadConfig(): wrong ElasticTranscoder config returned. Want %#v. Got %#v.", *expectedCfg.ElasticTranscoder, *cfg.ElasticTranscoder)
	}
	if !reflect.DeepEqual(*cfg.ElementalConductor, *expectedCfg.ElementalConductor) {
		t.Errorf("LoadConfig(): wrong Elemental Conductor config returned. Want %#v. Got %#v.", *expectedCfg.ElementalConductor, *cfg.ElementalConductor)
	}
	if !reflect.DeepEqual(*cfg.Server, *expectedCfg.Server) {
		t.Errorf("LoadConfig(): wrong Server config returned. Want %#v. Got %#v.", *expectedCfg.Server, *cfg.Server)
	}
}

func setEnvs(envs map[string]string) {
	for k, v := range envs {
		os.Setenv(k, v)
	}
}
