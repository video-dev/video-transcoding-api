package bitmovin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/provider"
	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/bitmovintypes"
	"github.com/bitmovin/bitmovin-go/models"
)

func TestFactoryIsRegistered(t *testing.T) {
	_, err := provider.GetProviderFactory(Name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBitmovinFactory(t *testing.T) {
	cfg := config.Config{
		Bitmovin: &config.Bitmovin{
			APIKey:           "apikey",
			Endpoint:         "this.is.my.endpoint",
			Timeout:          uint(5),
			AccessKeyID:      "accesskey",
			SecretAccessKey:  "secretaccesskey",
			Destination:      "s3://some-output-bucket/",
			EncodingRegion:   "AWS_US_EAST_1",
			AWSStorageRegion: "US_EAST_1",
		},
	}
	provider, err := bitmovinFactory(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	bitmovinProvider, ok := provider.(*bitmovinProvider)
	if !ok {
		t.Fatalf("Wrong provider returned. Want elementalConductorProvider instance. Got %#v.", provider)
	}
	expected := &bitmovin.Bitmovin{
		APIKey:     stringToPtr("apikey"),
		APIBaseURL: stringToPtr("this.is.my.endpoint"),
	}
	if *bitmovinProvider.client.APIKey != *expected.APIKey {
		t.Errorf("Factory: wrong APIKey returned. Want %#v. Got %#v.", expected.APIKey, *bitmovinProvider.client.APIKey)
	}
	if *bitmovinProvider.client.APIBaseURL != *expected.APIBaseURL {
		t.Errorf("Factory: wrong APIKey returned. Want %#v. Got %#v.", expected.APIBaseURL, *bitmovinProvider.client.APIBaseURL)
	}
}

func TestCreateH264Preset(t *testing.T) {
	testPresetName := "this_is_an_audio_config_uuid"
	preset := getH264Preset()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/configurations/audio/aac":
			resp := models.AACCodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.AACCodecConfigurationData{
					Result: models.AACCodecConfiguration{
						ID: stringToPtr("this_is_an_audio_config_uuid"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264":
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						ID: stringToPtr(testPresetName),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	presetName, err := prov.CreatePreset(preset)
	if err != nil {
		t.Fatal(err)
	}
	if presetName != testPresetName {
		t.Error("expected ", testPresetName, "got ", presetName)
	}
}

func TestCreateVP8Preset(t *testing.T) {
	testPresetName := "this_is_an_audio_config_uuid"
	preset := getVP8Preset()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/configurations/audio/vorbis":
			resp := models.VorbisCodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.VorbisCodecConfigurationData{
					Result: models.VorbisCodecConfiguration{
						ID: stringToPtr("this_is_an_audio_config_uuid"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/vp8":
			resp := models.VP8CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.VP8CodecConfigurationData{
					Result: models.VP8CodecConfiguration{
						ID: stringToPtr(testPresetName),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	presetName, err := prov.CreatePreset(preset)
	if err != nil {
		t.Fatal(err)
	}
	if presetName != testPresetName {
		t.Error("expected ", testPresetName, "got ", presetName)
	}
}

func TestCreatePresetFailsOnAPIError(t *testing.T) {
	preset := getH264Preset()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/configurations/audio/aac":
			resp := models.AACCodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusError,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	_, err := prov.CreatePreset(preset)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
}

func TestCreatePresetFailsOnGenericError(t *testing.T) {
	preset := getH264Preset()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/configurations/audio/aac":
			fmt.Fprintln(w, "Not proper json")
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	_, err := prov.CreatePreset(preset)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
}

func TestDeletePresetH264(t *testing.T) {
	testPresetID := "i_want_to_delete_this"
	audioPresetID := "embedded_audio_id"
	customData := make(map[string]interface{})
	customData["audio"] = audioPresetID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/configurations/video/h264/" + testPresetID + "/customData":
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/audio/aac/" + audioPresetID:
			resp := models.AACCodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/" + testPresetID:
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	err := prov.DeletePreset(testPresetID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeletePresetVP8(t *testing.T) {
	testPresetID := "i_want_to_delete_this"
	audioPresetID := "embedded_audio_id"
	customData := make(map[string]interface{})
	customData["audio"] = audioPresetID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/configurations/video/vp8/" + testPresetID + "/customData":
			resp := models.VP8CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.VP8CodecConfigurationData{
					Result: models.VP8CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/audio/vorbis/" + audioPresetID:
			resp := models.VorbisCodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/" + testPresetID:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - no API found with those values"))
		case "/encoding/configurations/video/vp8/" + testPresetID:
			resp := models.VP8CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	err := prov.DeletePreset(testPresetID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeletePresetFailsOnAPIError(t *testing.T) {
	testPresetID := "i_want_to_delete_this"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/configurations/video/h264/" + testPresetID:
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusError,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	err := prov.DeletePreset(testPresetID)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
}

func TestDeletePresetFailsOnGenericErrors(t *testing.T) {
	testPresetID := "i_want_to_delete_this"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/configurations/video/h264/" + testPresetID:
			fmt.Fprintln(w, "Not proper json")
		case "/encoding/configurations/video/vp8/" + testPresetID:
			fmt.Fprintln(w, "Not proper json")
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	err := prov.DeletePreset(testPresetID)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
}

func TestGetPresetH264(t *testing.T) {
	testPresetID := "this_is_a_video_preset_id"
	audioPresetID := "this_is_a_audio_preset_id"
	customData := make(map[string]interface{})
	customData["audio"] = audioPresetID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/configurations/video/h264/" + testPresetID + "/customData":
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/audio/aac/" + audioPresetID:
			resp := models.AACCodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/" + testPresetID:
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	i, err := prov.GetPreset(testPresetID)
	if err != nil {
		t.Fatal(err)
	}
	expected := bitmovinH264Preset{
		Video: models.H264CodecConfiguration{CustomData: customData},
		Audio: models.AACCodecConfiguration{},
	}
	if !reflect.DeepEqual(i, expected) {
		t.Errorf("GetPreset: want %#v. Got %#v", expected, i)
	}
}

func TestGetPresetVP8(t *testing.T) {
	testPresetID := "this_is_a_video_preset_id"
	audioPresetID := "this_is_a_audio_preset_id"
	customData := make(map[string]interface{})
	customData["audio"] = audioPresetID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/configurations/video/h264/" + testPresetID:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - no API found with those values"))
		case "/encoding/configurations/video/vp8/" + testPresetID:
			resp := models.VP8CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/vp8/" + testPresetID + "/customData":
			resp := models.VP8CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.VP8CodecConfigurationData{
					Result: models.VP8CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/audio/vorbis/" + audioPresetID:
			resp := models.VorbisCodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	i, err := prov.GetPreset(testPresetID)
	if err != nil {
		t.Fatal(err)
	}
	expected := bitmovinVP8Preset{
		Video: models.VP8CodecConfiguration{CustomData: customData},
		Audio: models.VorbisCodecConfiguration{},
	}
	if !reflect.DeepEqual(i, expected) {
		t.Errorf("GetPreset: want %#v. Got %#v", expected, i)
	}
}

func TestGetPresetFailsOnAPIError(t *testing.T) {
	testPresetID := "this_is_a_video_preset_id"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/configurations/video/h264/" + testPresetID:
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusError,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	i, err := prov.GetPreset(testPresetID)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	if i != nil {
		t.Errorf("GetPreset: got unexpected non-nil result: %#v", i)
	}
}

func TestGetPresetFailsOnGenericErrors(t *testing.T) {
	testPresetID := "this_is_a_video_preset_id"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/configurations/video/h264/" + testPresetID:
			fmt.Fprintln(w, "Not proper json")
		case "/encoding/configurations/video/vp8/" + testPresetID:
			fmt.Fprintln(w, "Not proper json")
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	i, err := prov.GetPreset(testPresetID)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	if i != nil {
		t.Errorf("GetPreset: got unexpected non-nil result: %#v", i)
	}
}

func TestTranscodeWithS3Input(t *testing.T) {
	s3InputID := "this_is_the_s3_input_id"
	s3OutputID := "this_is_the_s3_output_id"
	encodingID := "this_is_the_master_encoding_id"
	manifestID := "this_is_the_master_manifest_id"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/inputs/s3":
			resp := models.S3InputResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.S3InputData{
					Result: models.S3InputItem{
						ID: stringToPtr(s3InputID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/outputs/s3":
			resp := models.S3OutputResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.S3OutputData{
					Result: models.S3OutputItem{
						ID: stringToPtr(s3OutputID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID1/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID1"
			customData["container"] = "mp4"
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID2/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID2"
			customData["container"] = "m3u8"
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID3/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID3"
			customData["container"] = "m3u8"
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/vp8/videoID4/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID4"
			customData["container"] = "webm"
			resp := models.VP8CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.VP8CodecConfigurationData{
					Result: models.VP8CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID5/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID5"
			customData["container"] = "mov"
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls":
			resp := models.HLSManifestResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.HLSManifestData{
					Result: models.HLSManifest{
						ID: stringToPtr(manifestID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings":
			resp := models.EncodingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.EncodingData{
					Result: models.Encoding{
						ID: stringToPtr(encodingID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID1",
			"/encoding/configurations/video/h264/videoID2",
			"/encoding/configurations/video/h264/videoID3",
			"/encoding/configurations/video/h264/videoID5":
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID4":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - no API found with those values"))
		case "/encoding/configurations/video/vp8/videoID4":
			resp := models.VP8CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/streams":
			resp := models.StreamResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StreamData{
					Result: models.Stream{
						ID: stringToPtr("this_is_a_stream_id"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/muxings/mp4":
			resp := models.MP4MuxingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/muxings/ts":
			resp := models.TSMuxingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.TSMuxingData{
					Result: models.TSMuxing{
						ID: stringToPtr("this_is_a_ts_muxing_id"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/muxings/progressive-webm":
			resp := models.MP4MuxingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/muxings/progressive-mov":
			resp := models.ProgressiveMOVMuxingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls/" + manifestID + "/media":
			resp := models.MediaInfoResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls/" + manifestID + "/streams":
			resp := models.StreamInfoResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/start":
			resp := models.StartStopResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit " + r.URL.Path))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.Transcode(getJob("s3://bucket/folder/filename.mp4"))
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: encodingID,
		Status:        provider.StatusQueued,
	}
	if !reflect.DeepEqual(jobStatus, expectedJobStatus) {
		t.Errorf("Job Status: want %#v. Got %#v", expectedJobStatus, jobStatus)
	}
}

func TestTranscodeWithHTTPInput(t *testing.T) {
	httpInputID := "this_is_the_s3_input_id"
	s3OutputID := "this_is_the_s3_output_id"
	encodingID := "this_is_the_master_encoding_id"
	manifestID := "this_is_the_master_manifest_id"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/inputs/http":
			resp := models.HTTPInputResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.HTTPInputData{
					Result: models.HTTPInputItem{
						ID: stringToPtr(httpInputID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/outputs/s3":
			resp := models.S3OutputResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.S3OutputData{
					Result: models.S3OutputItem{
						ID: stringToPtr(s3OutputID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID1/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID1"
			customData["container"] = "mp4"
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID2/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID2"
			customData["container"] = "m3u8"
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID3/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID3"
			customData["container"] = "m3u8"
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/vp8/videoID4/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID4"
			customData["container"] = "webm"
			resp := models.VP8CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.VP8CodecConfigurationData{
					Result: models.VP8CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID5/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID5"
			customData["container"] = "mov"
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls":
			resp := models.HLSManifestResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.HLSManifestData{
					Result: models.HLSManifest{
						ID: stringToPtr(manifestID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings":
			resp := models.EncodingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.EncodingData{
					Result: models.Encoding{
						ID: stringToPtr(encodingID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID1",
			"/encoding/configurations/video/h264/videoID2",
			"/encoding/configurations/video/h264/videoID3",
			"/encoding/configurations/video/h264/videoID5":
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID4":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - no API found with those values"))
		case "/encoding/configurations/video/vp8/videoID4":
			resp := models.VP8CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/streams":
			resp := models.StreamResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StreamData{
					Result: models.Stream{
						ID: stringToPtr("this_is_a_stream_id"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/muxings/mp4":
			resp := models.MP4MuxingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/muxings/ts":
			resp := models.TSMuxingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.TSMuxingData{
					Result: models.TSMuxing{
						ID: stringToPtr("this_is_a_ts_muxing_id"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/muxings/progressive-webm":
			resp := models.MP4MuxingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/muxings/progressive-mov":
			resp := models.ProgressiveMOVMuxingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls/" + manifestID + "/media":
			resp := models.MediaInfoResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls/" + manifestID + "/streams":
			resp := models.StreamInfoResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/start":
			resp := models.StartStopResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit " + r.URL.Path))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.Transcode(getJob("http://bucket.com/folder/filename.mp4"))
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: encodingID,
		Status:        provider.StatusQueued,
	}
	if !reflect.DeepEqual(jobStatus, expectedJobStatus) {
		t.Errorf("Job Status: want %#v. Got %#v", expectedJobStatus, jobStatus)
	}
}

func TestTranscodeWithHTTPSInput(t *testing.T) {
	httpsInputID := "this_is_the_s3_input_id"
	s3OutputID := "this_is_the_s3_output_id"
	encodingID := "this_is_the_master_encoding_id"
	manifestID := "this_is_the_master_manifest_id"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/inputs/https":
			resp := models.HTTPSInputResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.HTTPSInputData{
					Result: models.HTTPSInputItem{
						ID: stringToPtr(httpsInputID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/outputs/s3":
			resp := models.S3OutputResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.S3OutputData{
					Result: models.S3OutputItem{
						ID: stringToPtr(s3OutputID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID1/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID1"
			customData["container"] = "mp4"
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID2/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID2"
			customData["container"] = "m3u8"
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID3/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID3"
			customData["container"] = "m3u8"
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/vp8/videoID4/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID4"
			customData["container"] = "webm"
			resp := models.VP8CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.VP8CodecConfigurationData{
					Result: models.VP8CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID5/customData":
			customData := make(map[string]interface{})
			customData["audio"] = "audioID5"
			customData["container"] = "mov"
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.H264CodecConfigurationData{
					Result: models.H264CodecConfiguration{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls":
			resp := models.HLSManifestResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.HLSManifestData{
					Result: models.HLSManifest{
						ID: stringToPtr(manifestID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings":
			resp := models.EncodingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.EncodingData{
					Result: models.Encoding{
						ID: stringToPtr(encodingID),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID1",
			"/encoding/configurations/video/h264/videoID2",
			"/encoding/configurations/video/h264/videoID3",
			"/encoding/configurations/video/h264/videoID5":
			resp := models.H264CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/configurations/video/h264/videoID4":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - no API found with those values"))
		case "/encoding/configurations/video/vp8/videoID4":
			resp := models.VP8CodecConfigurationResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/streams":
			resp := models.StreamResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StreamData{
					Result: models.Stream{
						ID: stringToPtr("this_is_a_stream_id"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/muxings/mp4":
			resp := models.MP4MuxingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/muxings/ts":
			resp := models.TSMuxingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.TSMuxingData{
					Result: models.TSMuxing{
						ID: stringToPtr("this_is_a_ts_muxing_id"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/muxings/progressive-webm":
			resp := models.MP4MuxingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/muxings/progressive-mov":
			resp := models.ProgressiveMOVMuxingResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls/" + manifestID + "/media":
			resp := models.MediaInfoResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls/" + manifestID + "/streams":
			resp := models.StreamInfoResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + encodingID + "/start":
			resp := models.StartStopResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit " + r.URL.Path))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.Transcode(getJob("https://bucket.com/folder/filename.mp4"))
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: encodingID,
		Status:        provider.StatusQueued,
	}
	if !reflect.DeepEqual(jobStatus, expectedJobStatus) {
		t.Errorf("Job Status: want %#v. Got %#v", expectedJobStatus, jobStatus)
	}
}

func TestTranscodeFailsOnAPIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/outputs/s3":
			resp := models.S3OutputResponse{
				Status: bitmovintypes.ResponseStatusError,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit " + r.URL.Path))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.Transcode(getJob("s3://bucket/folder/filename.mp4"))
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	if jobStatus != nil {
		t.Errorf("Transcode: got unexpected non-nil result: %#v", jobStatus)
	}
}

func TestTranscodeFailsOnGenericError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/outputs/s3":
			fmt.Fprintln(w, "Not proper json")
		default:
			t.Fatal(errors.New("unexpected path hit " + r.URL.Path))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.Transcode(getJob("s3://bucket/folder/filename.mp4"))
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	if jobStatus != nil {
		t.Errorf("Transcode: got unexpected non-nil result: %#v", jobStatus)
	}
}

func TestJobStatusReturnsFinishedIfEncodeAndManifestAreFinished(t *testing.T) {
	const (
		testJobID    = "this_is_a_job_id"
		manifestID   = "this_is_the_underlying_manifest_id"
		mp4MuxingID  = "test_mp4_muxing_id"
		webmMuxingID = "test_webm_muxing_id"
		movMuxingID  = "test_mov_muxing_id"
	)

	customData := make(map[string]interface{})
	customData["manifest"] = manifestID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status:   stringToPtr("FINISHED"),
						Progress: floatToPtr(100),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/customData":
			resp := models.CustomDataResponse{
				Data: models.Data{
					Result: models.Result{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls/" + manifestID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status: stringToPtr("FINISHED"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/streams":
			resp := models.StreamListResponse{
				Data: models.StreamListData{
					Result: models.StreamListResult{
						Items: []models.Stream{{ID: stringToPtr("new_stream")}},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/streams/new_stream/input":
			resp := models.StreamInputResponse{
				Data: models.StreamInputData{
					Result: models.StreamInputResult{
						Duration: floatToPtr(42.1),
						Bitrate:  intToPtr(25e5),
						VideoStreams: []models.StreamInputVideo{
							{
								ID:       stringToPtr("hd-stream"),
								Codec:    stringToPtr("h265"),
								Duration: floatToPtr(42.1),
								Width:    intToPtr(1280),
								Height:   intToPtr(720),
							},
						},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/muxings/mp4":
			resp := models.MP4MuxingListResponse{
				Data: models.MP4MuxingListData{
					Result: models.MP4MuxingListResult{
						TotalCount: intToPtr(1),
						Items: []models.MP4Muxing{{
							ID:       stringToPtr(mp4MuxingID),
							Filename: stringToPtr("test_file.mp4"),
						}},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/muxings/mp4/" + mp4MuxingID + "/information":
			resp := models.MP4MuxingInformationResponse{
				Data: models.MP4MuxingInformationData{
					Result: models.MP4MuxingInformationResult{
						ContainerFormat: stringToPtr("mpeg-4"),
						FileSize:        intToPtr(3),
						VideoTracks: []models.VideoTrack{{
							Codec:       stringToPtr("h264"),
							FrameWidth:  intToPtr(1280),
							FrameHeight: intToPtr(720),
						}},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/muxings/progressive-webm":
			resp := models.ProgressiveWebMMuxingListResponse{
				Data: models.ProgressiveWebMMuxingListData{
					Result: models.ProgressiveWebMMuxingListResult{
						TotalCount: intToPtr(1),
						Items: []models.ProgressiveWebMMuxing{
							{
								ID:       stringToPtr(webmMuxingID),
								Filename: stringToPtr("test_file.webm"),
							},
						},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/muxings/progressive-webm/" + webmMuxingID + "/information":
			resp := models.ProgressiveWebMMuxingInformationResponse{
				Data: models.ProgressiveWebMMuxingInformationData{
					Result: models.ProgressiveWebMMuxingInformationResult{
						ContainerFormat: stringToPtr("webm"),
						FileSize:        intToPtr(9),
						VideoTracks: []models.VideoTrack{{
							Codec:       stringToPtr("vp8"),
							FrameWidth:  intToPtr(1280),
							FrameHeight: intToPtr(720),
						}},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/muxings/progressive-mov":
			resp := models.ProgressiveMOVMuxingListResponse{
				Data: models.ProgressiveMOVMuxingListData{
					Result: models.ProgressiveMOVMuxingListResult{
						TotalCount: intToPtr(1),
						Items: []models.ProgressiveMOVMuxing{
							{
								ID:       stringToPtr(movMuxingID),
								Filename: stringToPtr("test_file.mov"),
							},
						},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/muxings/progressive-mov/" + movMuxingID + "/information":
			resp := models.ProgressiveMOVMuxingInformationResponse{
				Data: models.ProgressiveMOVMuxingInformationData{
					Result: models.ProgressiveMOVMuxingInformationResult{
						ContainerFormat: stringToPtr("mov"),
						FileSize:        intToPtr(9),
						VideoTracks: []models.VideoTrack{{
							Codec:       stringToPtr("h264"),
							FrameWidth:  intToPtr(1920),
							FrameHeight: intToPtr(1080),
						}},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatalf("unexpected path hit: %v", r.URL.Path)
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: testJobID})
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: testJobID,
		Status:        provider.StatusFinished,
		Progress:      100,
		ProviderStatus: map[string]interface{}{
			"message":        "",
			"originalStatus": "FINISHED",
			"manifestStatus": "FINISHED",
		},
		SourceInfo: provider.SourceInfo{
			Duration:   42100 * time.Millisecond,
			VideoCodec: "h265",
			Width:      1280,
			Height:     720,
		},
		Output: provider.JobOutput{
			Destination: "s3://some-output-bucket/job-123/",
			Files: []provider.OutputFile{
				{
					Path:       "s3://some-output-bucket/job-123/test_file.mp4",
					Container:  "mpeg-4",
					FileSize:   3,
					VideoCodec: "h264",
					Width:      1280,
					Height:     720,
				},
				{
					Path:       "s3://some-output-bucket/job-123/test_file.webm",
					Container:  "webm",
					FileSize:   9,
					VideoCodec: "vp8",
					Width:      1280,
					Height:     720,
				},
				{
					Path:       "s3://some-output-bucket/job-123/test_file.mov",
					Container:  "mov",
					FileSize:   9,
					VideoCodec: "h264",
					Width:      1920,
					Height:     1080,
				},
			},
		},
	}
	if !reflect.DeepEqual(jobStatus, expectedJobStatus) {
		t.Errorf("Job Status\nWant %#v\nGot  %#v", expectedJobStatus, jobStatus)
	}
}

func TestJobStatusReturnsFinishedIfEncodeIsFinishedAndNoManifestGenerationIsNeeded(t *testing.T) {
	const (
		testJobID    = "this_is_a_job_id"
		mp4MuxingID  = "test_mp4_muxing_id"
		webmMuxingID = "test_webm_muxing_id"
		movMuxingID  = "test_mov_muxing_id"
	)

	customData := make(map[string]interface{})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status:   stringToPtr("FINISHED"),
						Progress: floatToPtr(100),
					},
					Message: stringToPtr("it's done!"),
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/customData":
			resp := models.CustomDataResponse{
				Data: models.Data{
					Result: models.Result{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/streams":
			resp := models.StreamListResponse{
				Data: models.StreamListData{
					Result: models.StreamListResult{
						Items: []models.Stream{{ID: stringToPtr("test_stream")}},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/streams/test_stream/input":
			resp := models.StreamInputResponse{
				Data: models.StreamInputData{
					Result: models.StreamInputResult{
						Duration: floatToPtr(32),
						Bitrate:  intToPtr(25e6),
						VideoStreams: []models.StreamInputVideo{
							{
								ID:       stringToPtr("video-stream-id"),
								Codec:    stringToPtr("h264"),
								Duration: floatToPtr(32.1),
								Width:    intToPtr(1920),
								Height:   intToPtr(1080),
							},
						},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/muxings/mp4":
			resp := models.MP4MuxingListResponse{
				Data: models.MP4MuxingListData{
					Result: models.MP4MuxingListResult{
						TotalCount: intToPtr(1),
						Items: []models.MP4Muxing{{
							ID:       stringToPtr(mp4MuxingID),
							Filename: stringToPtr("test_file.mp4"),
						}},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/muxings/mp4/" + mp4MuxingID + "/information":
			resp := models.MP4MuxingInformationResponse{
				Data: models.MP4MuxingInformationData{
					Result: models.MP4MuxingInformationResult{
						ContainerFormat: stringToPtr("mpeg-4"),
						FileSize:        intToPtr(3),
						VideoTracks: []models.VideoTrack{{
							Codec:       stringToPtr("h264"),
							FrameWidth:  intToPtr(1280),
							FrameHeight: intToPtr(720),
						}},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/muxings/progressive-webm":
			resp := models.ProgressiveWebMMuxingListResponse{
				Data: models.ProgressiveWebMMuxingListData{
					Result: models.ProgressiveWebMMuxingListResult{
						TotalCount: intToPtr(1),
						Items: []models.ProgressiveWebMMuxing{
							{
								ID:       stringToPtr(webmMuxingID),
								Filename: stringToPtr("test_file.webm"),
							},
						},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/muxings/progressive-webm/" + webmMuxingID + "/information":
			resp := models.ProgressiveWebMMuxingInformationResponse{
				Data: models.ProgressiveWebMMuxingInformationData{
					Result: models.ProgressiveWebMMuxingInformationResult{
						ContainerFormat: stringToPtr("webm"),
						FileSize:        intToPtr(9),
						VideoTracks: []models.VideoTrack{{
							Codec:       stringToPtr("vp8"),
							FrameWidth:  intToPtr(1280),
							FrameHeight: intToPtr(720),
						}},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/muxings/progressive-mov":
			resp := models.ProgressiveMOVMuxingListResponse{
				Data: models.ProgressiveMOVMuxingListData{
					Result: models.ProgressiveMOVMuxingListResult{
						TotalCount: intToPtr(1),
						Items: []models.ProgressiveMOVMuxing{
							{
								ID:       stringToPtr(movMuxingID),
								Filename: stringToPtr("test_file.mov"),
							},
						},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/muxings/progressive-mov/" + movMuxingID + "/information":
			resp := models.ProgressiveMOVMuxingInformationResponse{
				Data: models.ProgressiveMOVMuxingInformationData{
					Result: models.ProgressiveMOVMuxingInformationResult{
						ContainerFormat: stringToPtr("mov"),
						FileSize:        intToPtr(9),
						VideoTracks: []models.VideoTrack{{
							Codec:       stringToPtr("h264"),
							FrameWidth:  intToPtr(1920),
							FrameHeight: intToPtr(1080),
						}},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatalf("unexpected path hit: %v", r.URL.Path)
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: testJobID})
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: testJobID,
		Status:        provider.StatusFinished,
		Progress:      100,
		ProviderStatus: map[string]interface{}{
			"message":        "it's done!",
			"originalStatus": "FINISHED",
		},
		SourceInfo: provider.SourceInfo{
			Duration:   32 * time.Second,
			VideoCodec: "h264",
			Width:      1920,
			Height:     1080,
		},
		Output: provider.JobOutput{
			Destination: "s3://some-output-bucket/job-123/",
			Files: []provider.OutputFile{
				{
					Path:       "s3://some-output-bucket/job-123/test_file.mp4",
					Container:  "mpeg-4",
					FileSize:   3,
					VideoCodec: "h264",
					Width:      1280,
					Height:     720,
				},
				{
					Path:       "s3://some-output-bucket/job-123/test_file.webm",
					Container:  "webm",
					FileSize:   9,
					VideoCodec: "vp8",
					Width:      1280,
					Height:     720,
				},
				{
					Path:       "s3://some-output-bucket/job-123/test_file.mov",
					Container:  "mov",
					FileSize:   9,
					VideoCodec: "h264",
					Width:      1920,
					Height:     1080,
				},
			},
		},
	}
	if !reflect.DeepEqual(jobStatus, expectedJobStatus) {
		t.Errorf("Job Status\nWant %#v\nGot  %#v", expectedJobStatus, jobStatus)
	}
}

func TestJobStatusReturnsQueuedIfEncodeIsCreated(t *testing.T) {
	testJobID := "this_is_a_job_id"
	manifestID := "this_is_the_underlying_manifest_id"
	customData := make(map[string]interface{})
	customData["manifest"] = manifestID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status: stringToPtr("CREATED"),
					},
					Message: stringToPtr("pending, pending"),
				},
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: testJobID})
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: testJobID,
		Status:        provider.StatusQueued,
		ProviderStatus: map[string]interface{}{
			"message":        "pending, pending",
			"originalStatus": "CREATED",
		},
		Output: provider.JobOutput{
			Destination: "s3://some-output-bucket/job-123/",
		},
	}
	if !reflect.DeepEqual(jobStatus, expectedJobStatus) {
		t.Errorf("Job Status\nWant %#v\nGot  %#v", expectedJobStatus, jobStatus)
	}
}

func TestJobStatusReturnsStartedIfEncodeIsRunning(t *testing.T) {
	testJobID := "this_is_a_job_id"
	manifestID := "this_is_the_underlying_manifest_id"
	customData := make(map[string]interface{})
	customData["manifest"] = manifestID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status:   stringToPtr("RUNNING"),
						Progress: floatToPtr(33),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: testJobID})
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: testJobID,
		Progress:      33,
		Status:        provider.StatusStarted,
		ProviderStatus: map[string]interface{}{
			"message":        "",
			"originalStatus": "RUNNING",
		},
		Output: provider.JobOutput{
			Destination: "s3://some-output-bucket/job-123/",
		},
	}
	if !reflect.DeepEqual(jobStatus, expectedJobStatus) {
		t.Errorf("Job Status\nWant %#v\nGot  %#v", expectedJobStatus, jobStatus)
	}
}

func TestJobStatusReturnsFailedIfEncodeFailed(t *testing.T) {
	testJobID := "this_is_a_job_id"
	manifestID := "this_is_the_underlying_manifest_id"
	customData := make(map[string]interface{})
	customData["manifest"] = manifestID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status: stringToPtr("ERROR"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: testJobID})
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: testJobID,
		Status:        provider.StatusFailed,
		ProviderStatus: map[string]interface{}{
			"message":        "",
			"originalStatus": "ERROR",
		},
		Output: provider.JobOutput{
			Destination: "s3://some-output-bucket/job-123/",
		},
	}
	if !reflect.DeepEqual(jobStatus, expectedJobStatus) {
		t.Errorf("Job Status\nWant %#v\nGot  %#v", expectedJobStatus, jobStatus)
	}
}

func TestJobStatusReturnsUnknownOnAPIError(t *testing.T) {
	testJobID := "this_is_a_job_id"
	manifestID := "this_is_the_underlying_manifest_id"
	customData := make(map[string]interface{})
	customData["manifest"] = manifestID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusError,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: testJobID})
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: testJobID,
		Status:        provider.StatusUnknown,
		ProviderStatus: map[string]interface{}{
			"message":        "",
			"originalStatus": "",
		},
		Output: provider.JobOutput{
			Destination: "s3://some-output-bucket/job-123/",
		},
	}
	if !reflect.DeepEqual(jobStatus, expectedJobStatus) {
		t.Errorf("Job Status\nWant %#v\nGot  %#v", expectedJobStatus, jobStatus)
	}
}

func TestJobStatusReturnsErrorOnFailureToListStreams(t *testing.T) {
	testJobID := "this_is_a_job_id"
	manifestID := "this_is_the_underlying_manifest_id"
	customData := make(map[string]interface{})
	customData["manifest"] = manifestID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status:   stringToPtr("FINISHED"),
						Progress: floatToPtr(100),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/customData":
			resp := models.CustomDataResponse{
				Data: models.Data{
					Result: models.Result{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls/" + manifestID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status: stringToPtr("FINISHED"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/streams":
			w.Write([]byte("internal server error"))
		default:
			t.Fatalf("unexpected path hit: %v", r.URL.Path)
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: testJobID})
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	if jobStatus != nil {
		t.Errorf("Got unexpected non-nil JobStatus: %#v", jobStatus)
	}
}

func TestJobStatusReturnsErrorOnFailureToGetInputData(t *testing.T) {
	testJobID := "this_is_a_job_id"
	manifestID := "this_is_the_underlying_manifest_id"
	customData := make(map[string]interface{})
	customData["manifest"] = manifestID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status:   stringToPtr("FINISHED"),
						Progress: floatToPtr(100),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/customData":
			resp := models.CustomDataResponse{
				Data: models.Data{
					Result: models.Result{
						CustomData: customData,
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/manifests/hls/" + manifestID + "/status":
			resp := models.StatusResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
				Data: models.StatusData{
					Result: models.StatusResult{
						Status: stringToPtr("FINISHED"),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/streams":
			resp := models.StreamListResponse{
				Data: models.StreamListData{
					Result: models.StreamListResult{
						Items: []models.Stream{{ID: stringToPtr("new_stream")}},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/encoding/encodings/" + testJobID + "/streams/new_stream/input":
			w.Write([]byte("internal server error"))
		default:
			t.Fatalf("unexpected path hit: %v", r.URL.Path)
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: testJobID})
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	if jobStatus != nil {
		t.Errorf("Got unexpected non-nil JobStatus: %#v", jobStatus)
	}
}

func TestJobStatusReturnsErrorOnGenericError(t *testing.T) {
	testJobID := "this_is_a_job_id"
	manifestID := "this_is_the_underlying_manifest_id"
	customData := make(map[string]interface{})
	customData["manifest"] = manifestID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/status":
			fmt.Fprintln(w, "Not proper json")
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: testJobID})
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
	if jobStatus != nil {
		t.Errorf("Got unexpected non-nil JobStatus: %#v", jobStatus)
	}
}

func TestCancelJob(t *testing.T) {
	testJobID := "this_is_a_job_id"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/stop":
			resp := models.StartStopResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	err := prov.CancelJob(testJobID)
	if err != nil {
		t.Fatal(err)
	}

}

func TestCancelJobFailsOnAPIError(t *testing.T) {
	testJobID := "this_is_a_job_id"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/stop":
			resp := models.StartStopResponse{
				Status: bitmovintypes.ResponseStatusError,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	err := prov.CancelJob(testJobID)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
}

func TestCancelJobFailsOnGenericError(t *testing.T) {
	testJobID := "this_is_a_job_id"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings/" + testJobID + "/stop":
			fmt.Fprintln(w, "Not proper json")
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	err := prov.CancelJob(testJobID)
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
}

func TestHealthCheck(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings":
			resp := models.EncodingListResponse{
				Status: bitmovintypes.ResponseStatusSuccess,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	err := prov.Healthcheck()
	if err != nil {
		t.Fatal(err)
	}
}

func TestHealthCheckFailsOnAPIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings":
			resp := models.EncodingListResponse{
				Status: bitmovintypes.ResponseStatusError,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	err := prov.Healthcheck()
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
}

func TestHealthCheckFailsOnGenericError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/encoding/encodings":
			fmt.Fprintln(w, "Not proper json")
		default:
			t.Fatal(errors.New("unexpected path hit"))
		}
	}))
	defer ts.Close()
	prov := getBitmovinProvider(ts.URL)
	err := prov.Healthcheck()
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
}

func TestCapabilities(t *testing.T) {
	var prov bitmovinProvider
	expected := provider.Capabilities{
		InputFormats:  []string{"prores", "h264"},
		OutputFormats: []string{"mp4", "mov", "hls", "webm"},
		Destinations:  []string{"s3"},
	}
	cap := prov.Capabilities()
	if !reflect.DeepEqual(cap, expected) {
		t.Errorf("Capabilities: want %#v. Got %#v", expected, cap)
	}
}

func getBitmovinProvider(url string) bitmovinProvider {
	client := bitmovin.NewBitmovin("apikey", url+"/", int64(5))
	return bitmovinProvider{
		client: client,
		config: &config.Bitmovin{
			APIKey:          "apikey",
			Endpoint:        url + "/",
			Timeout:         uint(5),
			AccessKeyID:     "accesskey",
			SecretAccessKey: "secretaccesskey",
			Destination:     "s3://some-output-bucket/",
		},
	}
}

func getH264Preset() db.Preset {
	return db.Preset{
		Audio: db.AudioPreset{
			Bitrate: "128000",
			Codec:   "aac",
		},
		Container:   "mp4",
		Description: "my nice preset",
		Name:        "mp4_1080p",
		RateControl: "VBR",
		Video: db.VideoPreset{
			Profile:      "main",
			ProfileLevel: "3.1",
			Bitrate:      "3500000",
			Codec:        "h264",
			GopMode:      "fixed",
			GopSize:      "90",
			Height:       "1080",
		},
	}
}

func getVP8Preset() db.Preset {
	return db.Preset{
		Audio: db.AudioPreset{
			Bitrate: "128000",
			Codec:   "vorbis",
		},
		Container:   "mp4",
		Description: "my nice preset",
		Name:        "mp4_1080p",
		RateControl: "VBR",
		Video: db.VideoPreset{
			Profile:      "main",
			ProfileLevel: "3.1",
			Bitrate:      "3500000",
			Codec:        "vp8",
			GopMode:      "fixed",
			GopSize:      "90",
			Height:       "1080",
		},
	}
}

func getJob(sourceMedia string) *db.Job {
	presets := []db.PresetMap{
		{
			Name: "mp4_1080p",
			ProviderMapping: map[string]string{
				Name: "videoID1",
			},
			OutputOpts: db.OutputOptions{Extension: "mp4"},
		},
		{
			Name: "hls_360p",
			ProviderMapping: map[string]string{
				Name: "videoID2",
			},
			OutputOpts: db.OutputOptions{Extension: "m3u8"},
		},
		{
			Name: "hls_480p",
			ProviderMapping: map[string]string{
				Name: "videoID3",
			},
			OutputOpts: db.OutputOptions{Extension: "m3u8"},
		},
		{
			Name: "webm_480p",
			ProviderMapping: map[string]string{
				Name: "videoID4",
			},
			OutputOpts: db.OutputOptions{Extension: "webm"},
		},
		{
			Name: "1080p_synd",
			ProviderMapping: map[string]string{
				Name: "videoID5",
			},
			OutputOpts: db.OutputOptions{Extension: "mov"},
		},
	}
	outputs := make([]db.TranscodeOutput, len(presets))
	for i, preset := range presets {
		fileName := "output-" + preset.Name + "." + preset.OutputOpts.Extension
		if preset.OutputOpts.Extension == "m3u8" {
			fileName = "hls/output-" + preset.Name + ".m3u8"
		}
		outputs[i] = db.TranscodeOutput{
			Preset:   preset,
			FileName: fileName,
		}
	}
	job := &db.Job{
		ProviderName: Name,
		SourceMedia:  sourceMedia,
		StreamingParams: db.StreamingParams{
			SegmentDuration:  uint(4),
			PlaylistFileName: "hls/master_playlist.m3u8",
		},
		Outputs: outputs,
	}
	return job
}
