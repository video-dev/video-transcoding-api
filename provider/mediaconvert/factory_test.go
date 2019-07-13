package mediaconvert

import (
	"os"
	"reflect"
	"testing"

	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert"
	"github.com/google/go-cmp/cmp"
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
		wantCreds  aws.Credentials
		wantRegion string
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
			wantCreds: aws.Credentials{
				AccessKeyID:     "cfg_access_key_id",
				SecretAccessKey: "cfg_secret_access_key",
				Source:          aws.StaticCredentialsProviderName,
			},
			wantRegion: "us-cfg-region-1",
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
			wantCreds: aws.Credentials{
				AccessKeyID:     "env_access_key_id",
				SecretAccessKey: "env_secret_access_key",
				Source:          external.CredentialsSourceName,
			},
			wantRegion: "us-north-1",
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

			client, ok := p.client.(*mediaconvert.Client)
			if !ok {
				t.Error("factory returned a mediaconvert provider with a non-aws client implementation")
				return
			}

			creds, err := client.Credentials.Retrieve()
			if err != nil {
				t.Errorf("error retrieving aws credentials: %v", err)
			}

			if g, e := creds, tt.wantCreds; !reflect.DeepEqual(g, e) {
				t.Errorf("unexpected credentials\nWant %+v\nGot %+v\nDiff %s",
					e, g, cmp.Diff(e, g))
			}

			if g, e := client.Region, tt.wantRegion; g != e {
				t.Errorf("expected region %q, got %q", e, g)
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
