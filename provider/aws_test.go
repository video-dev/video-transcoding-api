package provider

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
	"github.com/nytm/video-transcoding-api/config"
)

func TestElasticTranscoderProvider(t *testing.T) {
	cfg := config.Config{
		ElasticTranscoder: &config.ElasticTranscoder{
			AccessKeyID:     "AKIANOTREALLY",
			SecretAccessKey: "really-secret",
			PipelineID:      "mypipeline",
			Region:          "sa-east-1",
		},
	}
	provider, err := ElasticTranscoderProvider(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	elasticProvider := provider.(*awsProvider)
	if !reflect.DeepEqual(*elasticProvider.config, *cfg.ElasticTranscoder) {
		t.Errorf("ElasticTranscoderProvider: did not store the proper config. Want %#v. Got %#v.", cfg.ElasticTranscoder, elasticProvider.config)
	}
	expectedCreds := credentials.Value{AccessKeyID: "AKIANOTREALLY", SecretAccessKey: "really-secret"}
	creds, err := elasticProvider.c.(*elastictranscoder.ElasticTranscoder).Config.Credentials.Get()
	if err != nil {
		t.Fatal(err)
	}

	// provider is not relevant
	creds.ProviderName = expectedCreds.ProviderName
	if !reflect.DeepEqual(creds, expectedCreds) {
		t.Errorf("ElasticTranscoderProvider: wrogn credentials. Want %#v. Got %#v.", expectedCreds, creds)
	}

	region := *elasticProvider.c.(*elastictranscoder.ElasticTranscoder).Config.Region
	if region != cfg.ElasticTranscoder.Region {
		t.Errorf("ElasticTranscoderProvider: wrong region. Want %q. Got %q.", cfg.ElasticTranscoder.Region, region)
	}
}

func TestElasticTranscoderProviderDefaultRegion(t *testing.T) {
	cfg := config.Config{
		ElasticTranscoder: &config.ElasticTranscoder{
			AccessKeyID:     "AKIANOTREALLY",
			SecretAccessKey: "really-secret",
			PipelineID:      "mypipeline",
		},
	}
	provider, err := ElasticTranscoderProvider(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	elasticProvider := provider.(*awsProvider)
	if !reflect.DeepEqual(*elasticProvider.config, *cfg.ElasticTranscoder) {
		t.Errorf("ElasticTranscoderProvider: did not store the proper config. Want %#v. Got %#v.", cfg.ElasticTranscoder, elasticProvider.config)
	}
	expectedCreds := credentials.Value{AccessKeyID: "AKIANOTREALLY", SecretAccessKey: "really-secret"}
	creds, err := elasticProvider.c.(*elastictranscoder.ElasticTranscoder).Config.Credentials.Get()
	if err != nil {
		t.Fatal(err)
	}

	// provider is not relevant
	creds.ProviderName = expectedCreds.ProviderName
	if !reflect.DeepEqual(creds, expectedCreds) {
		t.Errorf("ElasticTranscoderProvider: wrogn credentials. Want %#v. Got %#v.", expectedCreds, creds)
	}

	region := *elasticProvider.c.(*elastictranscoder.ElasticTranscoder).Config.Region
	if region != "us-east-1" {
		t.Errorf("ElasticTranscoderProvider: wrong region. Want %q. Got %q.", "us-east-1", region)
	}
}

func TestElasticTranscoderProviderValidation(t *testing.T) {
	var tests = []struct {
		accessKeyID     string
		secretAccessKey string
		pipelineID      string
	}{
		{"", "", ""},
		{"AKIANOTREALLY", "", ""},
		{"", "very-secret", ""},
		{"", "", "superpipeline"},
		{"AKIANOTREALLY", "very-secret", ""},
	}
	for _, test := range tests {
		cfg := config.Config{
			ElasticTranscoder: &config.ElasticTranscoder{
				AccessKeyID:     test.accessKeyID,
				SecretAccessKey: test.secretAccessKey,
				PipelineID:      test.pipelineID,
			},
		}
		provider, err := ElasticTranscoderProvider(&cfg)
		if provider != nil {
			t.Errorf("Got unexpected non-nil provider: %#v", provider)
		}
		if err != errAWSInvalidConfig {
			t.Errorf("Wrong error returned. Want errAWSInvalidConfig. Got %#v", err)
		}
	}
}
