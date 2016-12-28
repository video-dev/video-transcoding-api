package bitmovin

import (
	"reflect"
	"testing"

	"github.com/NYTimes/video-transcoding-api/provider"
)

func TestCapabilities(t *testing.T) {
	var prov bitmovinProvider
	expected := provider.Capabilities{
		InputFormats:  []string{"prores", "h264"},
		OutputFormats: []string{"mp4", "hls"},
		Destinations:  []string{"s3"},
	}
	cap := prov.Capabilities()
	if !reflect.DeepEqual(cap, expected) {
		t.Errorf("Capabilities: want %#v. Got %#v", expected, cap)
	}
}
