package mediaconvert

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/mediaconvert"
	"github.com/video-dev/video-transcoding-api/v2/config"
)

var cfgWithoutCredsAndRegion = config.Config{
	MediaConvert: &config.MediaConvert{
		Endpoint: "http://some/endpoint",
		Queue:    "arn:some:queue",
		Role:     "arn:some:role",
	},
}

var cfgWithCredsAndRegion = config.Config{
	MediaConvert: &config.MediaConvert{
		AccessKeyID:     "cfg_access_key_id",
		SecretAccessKey: "cfg_secret_access_key",
		Endpoint:        "http://some/endpoint",
		Queue:           "arn:some:queue",
		Role:            "arn:some:role",
		Region:          "us-cfg-region-1",
	},
}

func Test_mediaconvertFactory(t *testing.T) {
	tests := []struct {
		name       string
		envVars    map[string]string
		cfg        config.Config
		wantErrMsg string
	}{
		{
			name: "when a config specifies aws credentials and region, those credentials are used",
			envVars: map[string]string{
				"AWS_ACCESS_KEY_ID":     "env_access_key_id",
				"AWS_SECRET_ACCESS_KEY": "env_secret_access_key",
				"AWS_DEFAULT_REGION":    "us-north-1",
			},
			cfg: cfgWithCredsAndRegion,
		},
		{
			name: "when a config does not specify aws credentials or region, credentials and region are loaded " +
				"from the environment",
			envVars: map[string]string{
				"AWS_ACCESS_KEY_ID":     "env_access_key_id",
				"AWS_SECRET_ACCESS_KEY": "env_secret_access_key",
				"AWS_DEFAULT_REGION":    "us-north-1",
			},
			cfg: cfgWithoutCredsAndRegion,
		},
		{
			name:       "an incomplete cfg results in an error returned",
			cfg:        config.Config{MediaConvert: &config.MediaConvert{}},
			wantErrMsg: "incomplete MediaConvert config",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				resetFunc, err := setenvReset(k, v)
				if err != nil {
					t.Errorf("running os env reset: %v", err)
				}
				defer resetFunc()
			}

			provider, err := mediaconvertFactory(&tt.cfg)
			if err != nil {
				if tt.wantErrMsg != err.Error() {
					t.Errorf("mcProvider.CreatePreset() error = %v, wantErr %q", err, tt.wantErrMsg)
				}
				return
			}

			p, ok := provider.(*mcProvider)
			if !ok {
				t.Error("factory didn't return a mediaconvert provider")
				return
			}

			_, ok = p.client.(*mediaconvert.Client)
			if !ok {
				t.Error("factory returned a mediaconvert provider with a non-aws client implementation")
				return
			}
		})
	}
}

func setenvReset(name, val string) (resetEnv func(), rerr error) {
	cached := os.Getenv(name)
	err := os.Setenv(name, val)
	if err != nil {
		return nil, err
	}
	return func() {
		os.Setenv(name, cached)
	}, nil
}
