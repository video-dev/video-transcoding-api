// Package elementalcloud provides a implementation of the provider that uses the
// Elemental Cloud API for transcoding media files.
//
// It doesn't expose any public type. In order to use the provider, one must
// import this package and then grab the factory from the provider package:
//
//     import (
//         "github.com/nytm/video-transcoding-api/provider"
//         "github.com/nytm/video-transcoding-api/provider/elementalcloud"
//     )
//
//     func UseProvider() {
//         factory, err := provider.GetProviderFactory(elementalcloud.Name)
//         // handle err and use factory to get an instance of the provider.
//     }
package elementalcloud

import (
	"github.com/NYTimes/encoding-wrapper/elementalcloud"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/provider"
)

// Name is the name used for registering the Elemental Cloud provider in the
// registry of providers.
const Name = "elementalcloud"

var errElementalCloudInvalidConfig = provider.InvalidConfigError("missing Elemental user login or api key. Please define the environment variables ELEMENTALCLOUD_USER_LOGIN and ELEMENTALCLOUD_API_KEY or set these values in the configuration file")

func init() {
	provider.Register(Name, elementalCloudFactory)
}

type elementalCloudProvider struct {
	config *config.Config
	client *elementalcloud.Client
}

func (p *elementalCloudProvider) TranscodeWithPresets(source string, presets []string) (*provider.JobStatus, error) {
	return nil, nil
}

func (p *elementalCloudProvider) JobStatus(id string) (*provider.JobStatus, error) {
	_, err := p.client.GetJob(id)
	if err != nil {
		return nil, err
	}
	return &provider.JobStatus{ProviderName: Name}, nil
}

func elementalCloudFactory(cfg *config.Config) (provider.TranscodingProvider, error) {
	if cfg.ElementalCloud.Host == "" || cfg.ElementalCloud.UserLogin == "" ||
		cfg.ElementalCloud.APIKey == "" || cfg.ElementalCloud.AuthExpires == 0 {
		return nil, errElementalCloudInvalidConfig
	}
	client := elementalcloud.NewClient(
		cfg.ElementalCloud.Host,
		cfg.ElementalCloud.UserLogin,
		cfg.ElementalCloud.APIKey,
		cfg.ElementalCloud.AuthExpires,
	)
	return &elementalCloudProvider{client: client, config: cfg}, nil
}
