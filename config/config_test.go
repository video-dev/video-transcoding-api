package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/NYTimes/gizmo/config"
)

func TestLoadConfigFromFile(t *testing.T) {
	cleanEnvs()
	fileName := "testdata/config.json"
	cfg := LoadConfig(fileName)
	expectedCfg := Config{
		Server: &config.Server{
			HTTPPort:      8090,
			HTTPAccessLog: "/var/log/myapp/access.log",
		},
		Redis: &Redis{
			SentinelAddrs:      "127.0.0.1:26379,127.0.0.2:26379,127.0.0.3:26379",
			SentinelMasterName: "mymaster",
			RedisAddr:          "127.0.0.1:6379",
			Password:           "secret",
			PoolSize:           90,
			PoolTimeout:        5,
		},
		EncodingCom: &EncodingCom{
			UserID:      "myuser",
			UserKey:     "superkey",
			Destination: "http://nice-destination",
		},
	}
	if !reflect.DeepEqual(*cfg.Server, *expectedCfg.Server) {
		t.Errorf("LoadConfig(%q): wrong Server config returned. Want %#v. Got %#v.", fileName, *expectedCfg.Server, *cfg.Server)
	}
	if !reflect.DeepEqual(*cfg.Redis, *expectedCfg.Redis) {
		t.Errorf("LoadConfig(%q): wrong Redis config returned. Want %#v. Got %#v.", fileName, *expectedCfg.Redis, *cfg.Redis)
	}
	if !reflect.DeepEqual(*cfg.EncodingCom, *expectedCfg.EncodingCom) {
		t.Errorf("LoadConfig(%q): wrong EncodingCom config returned. Want %#v. Got %#v.", fileName, *expectedCfg.EncodingCom, *cfg.EncodingCom)
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	setEnvs(map[string]string{
		"SENTINEL_ADDRS":             "10.10.10.10:26379,10.10.10.11:26379,10.10.10.12:26379",
		"SENTINEL_MASTER_NAME":       "supermaster",
		"REDIS_ADDR":                 "localhost:6379",
		"REDIS_PASSWORD":             "super-secret",
		"REDIS_POOL_SIZE":            "100",
		"REDIS_POOL_TIMEOUT_SECONDS": "10",
		"ENCODINGCOM_USER_ID":        "myuser",
		"ENCODINGCOM_USER_KEY":       "secret-key",
		"ENCODINGCOM_DESTINATION":    "https://safe-stuff",
	})
	fileName := "testdata/empty.json"
	cfg := LoadConfig(fileName)
	expectedCfg := Config{
		Redis: &Redis{
			SentinelAddrs:      "10.10.10.10:26379,10.10.10.11:26379,10.10.10.12:26379",
			SentinelMasterName: "supermaster",
			RedisAddr:          "localhost:6379",
			Password:           "super-secret",
			PoolSize:           100,
			PoolTimeout:        10,
		},
		EncodingCom: &EncodingCom{
			UserID:      "myuser",
			UserKey:     "secret-key",
			Destination: "https://safe-stuff",
		},
	}
	if !reflect.DeepEqual(*cfg.Redis, *expectedCfg.Redis) {
		t.Errorf("LoadConfig(%q): wrong Redis config returned. Want %#v. Got %#v.", fileName, *expectedCfg.Redis, *cfg.Redis)
	}
	if !reflect.DeepEqual(*cfg.EncodingCom, *expectedCfg.EncodingCom) {
		t.Errorf("LoadConfig(%q): wrong EncodingCom config returned. Want %#v. Got %#v.", fileName, *expectedCfg.EncodingCom, *cfg.EncodingCom)
	}
}

func TestLoadConfigOverride(t *testing.T) {
	cleanEnvs()
	setEnvs(map[string]string{
		"REDIS_PASSWORD":          "super-secret",
		"ENCODINGCOM_USER_ID":     "myuser",
		"ENCODINGCOM_USER_KEY":    "secret-key",
		"ENCODINGCOM_DESTINATION": "https://safe-stuff",
	})
	fileName := "testdata/config.json"
	cfg := LoadConfig(fileName)
	expectedCfg := Config{
		Server: &config.Server{
			HTTPPort:      8090,
			HTTPAccessLog: "/var/log/myapp/access.log",
		},
		Redis: &Redis{
			SentinelAddrs:      "127.0.0.1:26379,127.0.0.2:26379,127.0.0.3:26379",
			SentinelMasterName: "mymaster",
			RedisAddr:          "127.0.0.1:6379",
			Password:           "super-secret",
			PoolSize:           90,
			PoolTimeout:        5,
		},
		EncodingCom: &EncodingCom{
			UserID:      "myuser",
			UserKey:     "secret-key",
			Destination: "https://safe-stuff",
		},
	}
	if !reflect.DeepEqual(*cfg.Server, *expectedCfg.Server) {
		t.Errorf("LoadConfig(%q): wrong Server config returned. Want %#v. Got %#v.", fileName, *expectedCfg.Server, *cfg.Server)
	}
	if !reflect.DeepEqual(*cfg.Redis, *expectedCfg.Redis) {
		t.Errorf("LoadConfig(%q): wrong Redis config returned. Want %#v. Got %#v.", fileName, *expectedCfg.Redis, *cfg.Redis)
	}
	if !reflect.DeepEqual(*cfg.EncodingCom, *expectedCfg.EncodingCom) {
		t.Errorf("LoadConfig(%q): wrong EncodingCom config returned. Want %#v. Got %#v.", fileName, *expectedCfg.EncodingCom, *cfg.EncodingCom)
	}
}

func cleanEnvs() {
	envs := []string{
		"SENTINEL_ADDRS", "SENTINEL_MASTER_NAME", "REDIS_ADDR",
		"REDIS_PASSWORD", "ENCODINGCOM_USER_ID", "ENCODINGCOM_USER_KEY",
		"ENCODINGCOM_DESTINATION", "REDIS_POOL_SIZE", "REDIS_POOL_TIMEOUT_SECONDS",
	}
	for _, env := range envs {
		os.Unsetenv(env)
	}
}

func setEnvs(envs map[string]string) {
	for k, v := range envs {
		os.Setenv(k, v)
	}
}
