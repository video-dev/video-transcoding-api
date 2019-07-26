package mediaconvert

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert"
	"github.com/google/go-cmp/cmp"
	"github.com/video-dev/video-transcoding-api/v2/config"
	"github.com/video-dev/video-transcoding-api/v2/db"
	"github.com/video-dev/video-transcoding-api/v2/provider"
)

var (
	defaultPreset = db.Preset{
		Name:        "preset_name",
		Description: "test_desc",
		Container:   "mp4",
		RateControl: "VBR",
		TwoPass:     true,
		Video: db.VideoPreset{
			Profile:       "high",
			ProfileLevel:  "4.1",
			Width:         "300",
			Height:        "400",
			Codec:         "h264",
			Bitrate:       "400000",
			GopSize:       "120",
			InterlaceMode: "progressive",
		},
		Audio: db.AudioPreset{
			Codec:   "aac",
			Bitrate: "20000",
		},
	}

	defaultJob = db.Job{
		ID:           "jobID",
		ProviderName: Name,
		SourceMedia:  "s3://some/path.mp4",
		Outputs: []db.TranscodeOutput{
			{
				Preset: db.PresetMap{
					Name: "preset1",
					ProviderMapping: map[string]string{
						"mediaconvert": "preset1",
					},
				},
				FileName: "file1.mp4",
			},
			{
				Preset: db.PresetMap{
					Name: "preset2",
					ProviderMapping: map[string]string{
						"mediaconvert": "preset2",
					},
				},
				FileName: "file2.mp4",
			},
		},
		StreamingParams: db.StreamingParams{
			SegmentDuration: 6,
		},
	}
)

func Test_mcProvider_CreatePreset(t *testing.T) {
	tests := []struct {
		name          string
		preset        db.Preset
		wantPresetReq mediaconvert.CreatePresetInput
		wantErr       bool
	}{
		{
			name:   "a valid h264/aac mp4 preset results in the expected mediaconvert preset sent to AWS API",
			preset: defaultPreset,
			wantPresetReq: mediaconvert.CreatePresetInput{
				Name:        aws.String("preset_name"),
				Description: aws.String("test_desc"),
				Settings: &mediaconvert.PresetSettings{
					ContainerSettings: &mediaconvert.ContainerSettings{
						Container: mediaconvert.ContainerTypeMp4,
					},
					VideoDescription: &mediaconvert.VideoDescription{
						Height:            aws.Int64(400),
						Width:             aws.Int64(300),
						RespondToAfd:      mediaconvert.RespondToAfdNone,
						ScalingBehavior:   mediaconvert.ScalingBehaviorDefault,
						TimecodeInsertion: mediaconvert.VideoTimecodeInsertionDisabled,
						AntiAlias:         mediaconvert.AntiAliasEnabled,
						CodecSettings: &mediaconvert.VideoCodecSettings{
							Codec: mediaconvert.VideoCodecH264,
							H264Settings: &mediaconvert.H264Settings{
								Bitrate:            aws.Int64(400000),
								CodecLevel:         mediaconvert.H264CodecLevelLevel41,
								CodecProfile:       mediaconvert.H264CodecProfileHigh,
								InterlaceMode:      mediaconvert.H264InterlaceModeProgressive,
								QualityTuningLevel: mediaconvert.H264QualityTuningLevelMultiPassHq,
								RateControlMode:    mediaconvert.H264RateControlModeVbr,
								GopSize:            aws.Float64(120),
							},
						},
					},
					AudioDescriptions: []mediaconvert.AudioDescription{
						{
							CodecSettings: &mediaconvert.AudioCodecSettings{
								Codec: mediaconvert.AudioCodecAac,
								AacSettings: &mediaconvert.AacSettings{
									Bitrate:         aws.Int64(20000),
									CodecProfile:    mediaconvert.AacCodecProfileLc,
									CodingMode:      mediaconvert.AacCodingModeCodingMode20,
									RateControlMode: mediaconvert.AacRateControlModeCbr,
									SampleRate:      aws.Int64(defaultAudioSampleRate),
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			client := &testMediaConvertClient{t: t}
			p := &mcProvider{client: client}
			_, err := p.CreatePreset(tt.preset)
			if (err != nil) != tt.wantErr {
				t.Errorf("mcProvider.CreatePreset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if g, e := *client.createPresetCalledWith, tt.wantPresetReq; !reflect.DeepEqual(g, e) {
				t.Fatalf("CreatePreset(): wrong preset request\nWant %+v\nGot %+v\nDiff %s", e,
					g, cmp.Diff(e, g))
			}
		})
	}
}

func Test_mcProvider_CreatePreset_fields(t *testing.T) {
	tests := []struct {
		name           string
		presetModifier func(preset db.Preset) db.Preset
		assertion      func(*mediaconvert.CreatePresetInput, *testing.T)
		wantErrMsg     string
	}{
		{
			name: "hls presets are set correctly",
			presetModifier: func(p db.Preset) db.Preset {
				p.Container = "m3u8"
				return p
			},
			assertion: func(input *mediaconvert.CreatePresetInput, t *testing.T) {
				if g, e := input.Settings.ContainerSettings.Container, mediaconvert.ContainerTypeM3u8; g != e {
					t.Fatalf("got %q, expected %q", g, e)
				}
			},
		},
		{
			name: "unrecognized containers return an error",
			presetModifier: func(p db.Preset) db.Preset {
				p.Container = "unrecognized"
				return p
			},
			wantErrMsg: `mapping preset container to MediaConvert container: container "unrecognized" not supported with mediaconvert`,
		},
		{
			name: "unrecognized h264 codec returns an error",
			presetModifier: func(p db.Preset) db.Preset {
				p.Video.Codec = "vp9001"
				return p
			},
			wantErrMsg: `generating video preset: video codec "vp9001" is not yet supported with mediaconvert`,
		},
		{
			name: "unrecognized h264 codec profile returns an error",
			presetModifier: func(p db.Preset) db.Preset {
				p.Video.Profile = "8000"
				return p
			},
			wantErrMsg: `generating video preset: h264 profile "8000" is not supported with mediaconvert`,
		},
		{
			name: "unrecognized h264 codec level returns an error",
			presetModifier: func(p db.Preset) db.Preset {
				p.Video.ProfileLevel = "9001"
				return p
			},
			wantErrMsg: `generating video preset: h264 level "9001" is not supported with mediaconvert`,
		},
		{
			name: "bad video width returns an error",
			presetModifier: func(p db.Preset) db.Preset {
				p.Video.Width = "s"
				return p
			},
			wantErrMsg: `generating video preset: parsing video width "s" to int64: strconv.ParseInt: parsing "s": invalid syntax`,
		},
		{
			name: "bad video height returns an error",
			presetModifier: func(p db.Preset) db.Preset {
				p.Video.Height = "h"
				return p
			},
			wantErrMsg: `generating video preset: parsing video height "h" to int64: strconv.ParseInt: parsing "h": invalid syntax`,
		},
		{
			name: "bad video bitrate returns an error",
			presetModifier: func(p db.Preset) db.Preset {
				p.Video.Bitrate = "bitrate"
				return p
			},
			wantErrMsg: `generating video preset: parsing video bitrate "bitrate" to int64: strconv.ParseInt: parsing "bitrate": invalid syntax`,
		},
		{
			name: "bad video gop size returns an error",
			presetModifier: func(p db.Preset) db.Preset {
				p.Video.GopSize = "gop"
				return p
			},
			wantErrMsg: `generating video preset: parsing gop size "gop" to int64: strconv.ParseFloat: parsing "gop": invalid syntax`,
		},
		{
			name: "unrecognized rate control mode returns an error",
			presetModifier: func(p db.Preset) db.Preset {
				p.RateControl = "not supported"
				return p
			},
			wantErrMsg: `generating video preset: rate control mode "not supported" is not supported with mediaconvert`,
		},
		{
			name: "unrecognized interlace modes return an error",
			presetModifier: func(p db.Preset) db.Preset {
				p.Video.InterlaceMode = "unsupported mode"
				return p
			},
			wantErrMsg: `generating video preset: h264 interlace mode "unsupported mode" is not supported with mediaconvert`,
		},
		{
			name: "unrecognized audio bitrate returns an error",
			presetModifier: func(p db.Preset) db.Preset {
				p.Audio.Bitrate = "aud_bitrate"
				return p
			},
			wantErrMsg: `generating audio preset: parsing audio bitrate "aud_bitrate" to int64: strconv.ParseInt: parsing "aud_bitrate": invalid syntax`,
		},
		{
			name: "unrecognized audio codec returns an error",
			presetModifier: func(p db.Preset) db.Preset {
				p.Audio.Codec = "aab"
				return p
			},
			wantErrMsg: `generating audio preset: audio codec "aab" is not yet supported with mediaconvert`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			client := &testMediaConvertClient{t: t}
			p := &mcProvider{client: client}
			_, err := p.CreatePreset(tt.presetModifier(defaultPreset))
			if err != nil && tt.wantErrMsg != err.Error() {
				t.Errorf("mcProvider.CreatePreset() error = %v, wantErr %q", err, tt.wantErrMsg)
				return
			}

			if tt.assertion != nil {
				tt.assertion(client.createPresetCalledWith, t)
			}
		})
	}
}

func Test_mcProvider_GetPreset(t *testing.T) {
	presetID := "some_preset"
	client := &testMediaConvertClient{t: t}
	p := &mcProvider{client: client}
	_, err := p.GetPreset(presetID)
	if err != nil {
		t.Fatalf("expected GetPreset() not to return an error, got: %v", err)
	}

	if g, e := *client.getPresetCalledWith, presetID; g != e {
		t.Fatalf("got %q, expected %q", g, e)
	}
}

func Test_mcProvider_DeletePreset(t *testing.T) {
	presetID := "some_preset_id"
	client := &testMediaConvertClient{t: t}
	p := &mcProvider{client: client}
	err := p.DeletePreset(presetID)
	if err != nil {
		t.Fatalf("expected DeletePreset() not to return an error, got: %v", err)
	}

	if g, e := client.deletePresetCalledWith, presetID; g != e {
		t.Fatalf("got %q, expected %q", g, e)
	}
}

func Test_mcProvider_Transcode(t *testing.T) {
	tests := []struct {
		name                string
		job                 *db.Job
		presetContainerType mediaconvert.ContainerType
		destination         string
		wantJobReq          mediaconvert.CreateJobInput
		wantErr             bool
	}{
		{
			name:                "a valid mp4 transcode job is mapped correctly to a mediaconvert job input",
			job:                 &defaultJob,
			presetContainerType: mediaconvert.ContainerTypeMp4,
			destination:         "s3://some/destination",
			wantJobReq: mediaconvert.CreateJobInput{
				Role:  aws.String(""),
				Queue: aws.String(""),
				Settings: &mediaconvert.JobSettings{
					Inputs: []mediaconvert.Input{
						{
							AudioSelectors: map[string]mediaconvert.AudioSelector{
								"Audio Selector 1": {
									DefaultSelection: mediaconvert.AudioDefaultSelectionDefault,
								},
							},
							FileInput: aws.String("s3://some/path.mp4"),
							VideoSelector: &mediaconvert.VideoSelector{
								ColorSpace: mediaconvert.ColorSpaceFollow,
							},
						},
					},
					OutputGroups: []mediaconvert.OutputGroup{
						{
							OutputGroupSettings: &mediaconvert.OutputGroupSettings{
								Type: mediaconvert.OutputGroupTypeFileGroupSettings,
								FileGroupSettings: &mediaconvert.FileGroupSettings{
									Destination: aws.String("s3://some/destination/jobID/"),
								},
							},
							Outputs: []mediaconvert.Output{
								{
									NameModifier: aws.String("file1"),
									Preset:       aws.String("preset1"),
									Extension:    aws.String("mp4"),
								},
								{
									NameModifier: aws.String("file2"),
									Preset:       aws.String("preset2"),
									Extension:    aws.String("mp4"),
								},
							},
						},
					},
				},
			},
		},
		{
			name:                "a valid hls transcode job is mapped correctly to a mediaconvert job input",
			job:                 &defaultJob,
			presetContainerType: mediaconvert.ContainerTypeM3u8,
			destination:         "s3://some/destination",
			wantJobReq: mediaconvert.CreateJobInput{
				Role:  aws.String(""),
				Queue: aws.String(""),
				Settings: &mediaconvert.JobSettings{
					Inputs: []mediaconvert.Input{
						{
							AudioSelectors: map[string]mediaconvert.AudioSelector{
								"Audio Selector 1": {
									DefaultSelection: mediaconvert.AudioDefaultSelectionDefault,
								},
							},
							FileInput: aws.String("s3://some/path.mp4"),
							VideoSelector: &mediaconvert.VideoSelector{
								ColorSpace: mediaconvert.ColorSpaceFollow,
							},
						},
					},
					OutputGroups: []mediaconvert.OutputGroup{
						{
							OutputGroupSettings: &mediaconvert.OutputGroupSettings{
								Type: mediaconvert.OutputGroupTypeHlsGroupSettings,
								HlsGroupSettings: &mediaconvert.HlsGroupSettings{
									Destination:            aws.String("s3://some/destination/jobID/"),
									SegmentLength:          aws.Int64(6),
									MinSegmentLength:       aws.Int64(0),
									DirectoryStructure:     mediaconvert.HlsDirectoryStructureSingleDirectory,
									ManifestDurationFormat: mediaconvert.HlsManifestDurationFormatFloatingPoint,
									OutputSelection:        mediaconvert.HlsOutputSelectionManifestsAndSegments,
									SegmentControl:         mediaconvert.HlsSegmentControlSegmentedFiles,
								},
							},
							Outputs: []mediaconvert.Output{
								{
									NameModifier: aws.String("file1"),
									Preset:       aws.String("preset1"),
									Extension:    aws.String("mp4"),
								},
								{
									NameModifier: aws.String("file2"),
									Preset:       aws.String("preset2"),
									Extension:    aws.String("mp4"),
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			client := &testMediaConvertClient{t: t, getPresetContainerType: tt.presetContainerType}
			p := &mcProvider{client: client, cfg: &config.MediaConvert{
				Destination: tt.destination,
			}}
			_, err := p.Transcode(tt.job)
			if (err != nil) != tt.wantErr {
				t.Errorf("mcProvider.Transcode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if g, e := client.createJobCalledWith, tt.wantJobReq; !reflect.DeepEqual(g, e) {
				t.Fatalf("Transcode(): wrong job request\nWant %+v\nGot %+v\nDiff %s", e,
					g, cmp.Diff(e, g))
			}
		})
	}
}

func Test_mcProvider_CancelJob(t *testing.T) {
	jobID := "some_job_id"
	client := &testMediaConvertClient{t: t}
	p := &mcProvider{client: client}
	err := p.CancelJob(jobID)
	if err != nil {
		t.Fatalf("expected CancelJob() not to return an error, got: %v", err)
	}

	if g, e := client.cancelJobCalledWith, jobID; g != e {
		t.Fatalf("got %q, expected %q", g, e)
	}
}

func Test_mcProvider_Healthcheck(t *testing.T) {
	client := &testMediaConvertClient{t: t}
	p := &mcProvider{client: client}

	err := p.Healthcheck()
	if err != nil {
		t.Fatalf("expected Healthcheck() not to return an error, got: %v", err)
	}

	if !client.listJobsCalled {
		t.Fatal("expected Healthcheck() to call ListJobs")
	}
}

func Test_mcProvider_JobStatus(t *testing.T) {
	tests := []struct {
		name        string
		destination string
		mcJob       mediaconvert.Job
		wantStatus  provider.JobStatus
		wantErr     bool
	}{
		{
			name:        "a job that has been queued returns the correct status",
			destination: "s3://some/destination",
			mcJob: mediaconvert.Job{
				Status: mediaconvert.JobStatusSubmitted,
			},
			wantStatus: provider.JobStatus{
				Status:       provider.StatusQueued,
				ProviderName: Name,
				Output: provider.JobOutput{
					Destination: "s3://some/destination/jobID/",
				},
			},
		},
		{
			name:        "a job that is currently transcoding returns the correct status",
			destination: "s3://some/destination",
			mcJob: mediaconvert.Job{
				Status:             mediaconvert.JobStatusProgressing,
				JobPercentComplete: aws.Int64(42),
			},
			wantStatus: provider.JobStatus{
				Status:       provider.StatusStarted,
				ProviderName: Name,
				Progress:     42,
				Output: provider.JobOutput{
					Destination: "s3://some/destination/jobID/",
				},
			},
		},
		{
			name:        "a job that has finished transcoding returns the correct status",
			destination: "s3://some/destination",
			mcJob: mediaconvert.Job{
				Status: mediaconvert.JobStatusComplete,
				OutputGroupDetails: []mediaconvert.OutputGroupDetail{{
					OutputDetails: []mediaconvert.OutputDetail{
						{
							VideoDetails: &mediaconvert.VideoDetail{
								HeightInPx: aws.Int64(2160),
								WidthInPx:  aws.Int64(3840),
							},
						},
						{
							VideoDetails: &mediaconvert.VideoDetail{
								HeightInPx: aws.Int64(1080),
								WidthInPx:  aws.Int64(1920),
							},
						},
					},
				}},
			},
			wantStatus: provider.JobStatus{
				Status:       provider.StatusFinished,
				ProviderName: Name,
				Progress:     100,
				Output: provider.JobOutput{
					Destination: "s3://some/destination/jobID/",
					Files: []provider.OutputFile{
						{Height: 2160, Width: 3840},
						{Height: 1080, Width: 1920},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			client := &testMediaConvertClient{
				t:                   t,
				jobReturnedByGetJob: tt.mcJob,
			}

			p := &mcProvider{client: client, cfg: &config.MediaConvert{
				Destination: tt.destination,
			}}

			status, err := p.JobStatus(&defaultJob)
			if (err != nil) != tt.wantErr {
				t.Errorf("mcProvider.JobStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if g, e := status, &tt.wantStatus; !reflect.DeepEqual(g, e) {
				t.Fatalf("mcProvider.JobStatus(): wrong job request\nWant %+v\nGot %+v\nDiff %s", e,
					g, cmp.Diff(e, g))
			}
		})
	}
}
