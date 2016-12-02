package elastictranscoder

import (
	"errors"
	"os"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/provider"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
	"github.com/kr/pretty"
)

func TestFactoryIsRegistered(t *testing.T) {
	_, err := provider.GetProviderFactory(Name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestElasticTranscoderProvider(t *testing.T) {
	cfg := config.Config{
		ElasticTranscoder: &config.ElasticTranscoder{
			AccessKeyID:     "AKIANOTREALLY",
			SecretAccessKey: "really-secret",
			PipelineID:      "mypipeline",
			Region:          "sa-east-1",
		},
	}
	provider, err := elasticTranscoderFactory(&cfg)
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
	provider, err := elasticTranscoderFactory(&cfg)
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
		provider, err := elasticTranscoderFactory(&cfg)
		if provider != nil {
			t.Errorf("Got unexpected non-nil provider: %#v", provider)
		}
		if err != errAWSInvalidConfig {
			t.Errorf("Wrong error returned. Want errAWSInvalidConfig. Got %#v", err)
		}
	}
}

func TestAWSTranscode(t *testing.T) {
	fakeTranscoder := newFakeElasticTranscoder()
	prov := &awsProvider{
		c: fakeTranscoder,
		config: &config.ElasticTranscoder{
			AccessKeyID:     "AKIA",
			SecretAccessKey: "secret",
			Region:          "sa-east-1",
			PipelineID:      "mypipeline",
		},
	}
	source := "dir/file.mov"
	outputs := []db.TranscodeOutput{
		{
			FileName: "output-720p.mp4",
			Preset: db.PresetMap{
				Name: "mp4_720p",
				ProviderMapping: map[string]string{
					Name:    "93239832-0001",
					"other": "irrelevant",
				},
				OutputOpts: db.OutputOptions{Extension: "mp4"},
			},
		},
		{
			FileName: "output-720p.webm",
			Preset: db.PresetMap{
				Name: "webm_720p",
				ProviderMapping: map[string]string{
					Name:    "93239832-0002",
					"other": "irrelevant",
				},
				OutputOpts: db.OutputOptions{Extension: "webm"},
			},
		},
		{
			FileName: "output-1080p.mov",
			Preset: db.PresetMap{
				Name: "mov_1080p",
				ProviderMapping: map[string]string{
					Name:    "93239832-0003",
					"other": "irrelevant",
				},
				OutputOpts: db.OutputOptions{Extension: "mov"},
			},
		},
	}

	jobStatus, err := prov.Transcode(&db.Job{
		ID:              "job-123",
		SourceMedia:     source,
		Outputs:         outputs,
		StreamingParams: db.StreamingParams{},
	})
	if err != nil {
		t.Fatal(err)
	}
	if m, _ := regexp.MatchString(`^job-[a-f0-9]{8}$`, jobStatus.ProviderJobID); !m {
		t.Errorf("Elastic Transcoder: invalid id returned - %q", jobStatus.ProviderJobID)
	}
	if jobStatus.Status != provider.StatusQueued {
		t.Errorf("Elastic Transcoder: wrong status returned. Want queued. Got %v", jobStatus.Status)
	}
	if jobStatus.ProviderName != Name {
		t.Errorf("Elastic Transcoder: wrong provider name returned. Want %q. Got %q", Name, jobStatus.ProviderName)
	}

	if len(fakeTranscoder.jobs) != 1 {
		t.Fatal("Did not send any job request to the server.")
	}
	jobInput := fakeTranscoder.jobs[jobStatus.ProviderJobID]

	expectedJobInput := elastictranscoder.CreateJobInput{
		PipelineId: aws.String("mypipeline"),
		Input: &elastictranscoder.JobInput{
			Key: aws.String(source),
			DetectedProperties: &elastictranscoder.DetectedProperties{
				DurationMillis: aws.Int64(120e3),
				FileSize:       aws.Int64(60356779),
				Height:         aws.Int64(1080),
				Width:          aws.Int64(1920),
			},
		},
		Outputs: []*elastictranscoder.CreateJobOutput{
			{PresetId: aws.String("93239832-0001"), Key: aws.String("job-123/output-720p.mp4")},
			{PresetId: aws.String("93239832-0002"), Key: aws.String("job-123/output-720p.webm")},
			{PresetId: aws.String("93239832-0003"), Key: aws.String("job-123/output-1080p.mov")},
		},
	}
	if !reflect.DeepEqual(*jobInput, expectedJobInput) {
		t.Errorf("Elastic Transcoder: wrong input\nWant %#v\nGot  %#v", expectedJobInput, *jobInput)
	}
}

func TestAWSTranscodeAdaptiveStreaming(t *testing.T) {
	fakeTranscoder := newFakeElasticTranscoder()
	prov := &awsProvider{
		c: fakeTranscoder,
		config: &config.ElasticTranscoder{
			AccessKeyID:     "AKIA",
			SecretAccessKey: "secret",
			Region:          "sa-east-1",
			PipelineID:      "mypipeline",
		},
	}
	source := "dir/file.mov"
	outputs := []db.TranscodeOutput{
		{
			FileName: "output_360p_hls/video.m3u8",
			Preset: db.PresetMap{
				Name: "hls_360p",
				ProviderMapping: map[string]string{
					Name:    "93239832-0001-hls",
					"other": "irrelevant",
				},
				OutputOpts: db.OutputOptions{Extension: "hls"},
			},
		},
		{
			FileName: "output_480p_hls/video.m3u8",
			Preset: db.PresetMap{
				Name: "hls_480p",
				ProviderMapping: map[string]string{
					Name:    "93239832-0002-hls",
					"other": "irrelevant",
				},
				OutputOpts: db.OutputOptions{Extension: "hls"},
			},
		},
		{
			FileName: "output_720p_hls/video.m3u8",
			Preset: db.PresetMap{
				Name: "hls_720p",
				ProviderMapping: map[string]string{
					Name:    "93239832-0003-hls",
					"other": "irrelevant",
				},
				OutputOpts: db.OutputOptions{Extension: "hls"},
			},
		},
	}

	jobStatus, err := prov.Transcode(&db.Job{
		ID:          "job-123",
		SourceMedia: source,
		Outputs:     outputs,
		StreamingParams: db.StreamingParams{
			PlaylistFileName: "video.m3u8",
			Protocol:         "asdf",
			SegmentDuration:  3,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if m, _ := regexp.MatchString(`^job-[a-f0-9]{8}$`, jobStatus.ProviderJobID); !m {
		t.Errorf("Elastic Transcoder: invalid id returned - %q", jobStatus.ProviderJobID)
	}
	if jobStatus.Status != provider.StatusQueued {
		t.Errorf("Elastic Transcoder: wrong status returned. Want queued. Got %v", jobStatus.Status)
	}
	if jobStatus.ProviderName != Name {
		t.Errorf("Elastic Transcoder: wrong provider name returned. Want %q. Got %q", Name, jobStatus.ProviderName)
	}

	if len(fakeTranscoder.jobs) != 1 {
		t.Fatal("Did not send any job request to the server.")
	}
	jobInput := fakeTranscoder.jobs[jobStatus.ProviderJobID]

	expectedJobInput := elastictranscoder.CreateJobInput{
		PipelineId: aws.String("mypipeline"),
		Input: &elastictranscoder.JobInput{
			Key: aws.String(source),
			DetectedProperties: &elastictranscoder.DetectedProperties{
				DurationMillis: aws.Int64(120e3),
				FileSize:       aws.Int64(60356779),
				Height:         aws.Int64(1080),
				Width:          aws.Int64(1920),
			},
		},
		Outputs: []*elastictranscoder.CreateJobOutput{
			{PresetId: aws.String("93239832-0001-hls"), Key: aws.String("job-123/output_360p_hls/video"), SegmentDuration: aws.String("3")},
			{PresetId: aws.String("93239832-0002-hls"), Key: aws.String("job-123/output_480p_hls/video"), SegmentDuration: aws.String("3")},
			{PresetId: aws.String("93239832-0003-hls"), Key: aws.String("job-123/output_720p_hls/video"), SegmentDuration: aws.String("3")},
		},
		Playlists: []*elastictranscoder.CreateJobPlaylist{
			{
				Format: aws.String("HLSv3"),
				Name:   aws.String("job-123/video"),
				OutputKeys: []*string{
					aws.String("job-123/output_360p_hls/video"),
					aws.String("job-123/output_480p_hls/video"),
					aws.String("job-123/output_720p_hls/video"),
				},
			},
		},
	}
	if !reflect.DeepEqual(*jobInput, expectedJobInput) {
		t.Errorf("Elastic Transcoder: wrong input\nWant %#v\nGot  %#v", expectedJobInput, *jobInput)
	}
}

func TestAWSTranscodeAdaptiveAndNonAdaptiveStreaming(t *testing.T) {
	fakeTranscoder := newFakeElasticTranscoder()
	prov := &awsProvider{
		c: fakeTranscoder,
		config: &config.ElasticTranscoder{
			AccessKeyID:     "AKIA",
			SecretAccessKey: "secret",
			Region:          "sa-east-1",
			PipelineID:      "mypipeline",
		},
	}
	source := "dir/file.mov"
	outputs := []db.TranscodeOutput{
		{
			FileName: "hls/output_hls_360p/video.m3u8",
			Preset: db.PresetMap{
				Name: "hls_360p",
				ProviderMapping: map[string]string{
					Name:    "93239832-0001-hls",
					"other": "irrelevant",
				},
				OutputOpts: db.OutputOptions{Extension: "hls"},
			},
		},
		{
			FileName: "hls/output_hls_480p/video.m3u8",
			Preset: db.PresetMap{
				Name: "hls_480p",
				ProviderMapping: map[string]string{
					Name:    "93239832-0002-hls",
					"other": "irrelevant",
				},
				OutputOpts: db.OutputOptions{Extension: "hls"},
			},
		},
		{
			FileName: "hls/output_hls_720p/video.m3u8",
			Preset: db.PresetMap{
				Name: "hls_720p",
				ProviderMapping: map[string]string{
					Name:    "93239832-0003-hls",
					"other": "irrelevant",
				},
				OutputOpts: db.OutputOptions{Extension: "hls"},
			},
		},
		{
			FileName: "output_720p.mp4",
			Preset: db.PresetMap{
				Name: "mp4_720p",
				ProviderMapping: map[string]string{
					Name:    "93239832-0004",
					"other": "irrelevant",
				},
				OutputOpts: db.OutputOptions{Extension: "mp4"},
			},
		},
	}

	jobStatus, err := prov.Transcode(&db.Job{
		ID:          "job-123",
		SourceMedia: source,
		Outputs:     outputs,
		StreamingParams: db.StreamingParams{
			PlaylistFileName: "hls/video.m3u8",
			Protocol:         "asdf",
			SegmentDuration:  3,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if m, _ := regexp.MatchString(`^job-[a-f0-9]{8}$`, jobStatus.ProviderJobID); !m {
		t.Errorf("Elastic Transcoder: invalid id returned - %q", jobStatus.ProviderJobID)
	}
	if jobStatus.Status != provider.StatusQueued {
		t.Errorf("Elastic Transcoder: wrong status returned. Want queued. Got %v", jobStatus.Status)
	}
	if jobStatus.ProviderName != Name {
		t.Errorf("Elastic Transcoder: wrong provider name returned. Want %q. Got %q", Name, jobStatus.ProviderName)
	}

	if len(fakeTranscoder.jobs) != 1 {
		t.Fatal("Did not send any job request to the server.")
	}
	jobInput := fakeTranscoder.jobs[jobStatus.ProviderJobID]

	expectedJobInput := elastictranscoder.CreateJobInput{
		PipelineId: aws.String("mypipeline"),
		Input: &elastictranscoder.JobInput{
			Key: aws.String(source),
			DetectedProperties: &elastictranscoder.DetectedProperties{
				DurationMillis: aws.Int64(120e3),
				FileSize:       aws.Int64(60356779),
				Height:         aws.Int64(1080),
				Width:          aws.Int64(1920),
			},
		},
		Outputs: []*elastictranscoder.CreateJobOutput{
			{PresetId: aws.String("93239832-0001-hls"), Key: aws.String("job-123/hls/output_hls_360p/video"), SegmentDuration: aws.String("3")},
			{PresetId: aws.String("93239832-0002-hls"), Key: aws.String("job-123/hls/output_hls_480p/video"), SegmentDuration: aws.String("3")},
			{PresetId: aws.String("93239832-0003-hls"), Key: aws.String("job-123/hls/output_hls_720p/video"), SegmentDuration: aws.String("3")},
			{PresetId: aws.String("93239832-0004"), Key: aws.String("job-123/output_720p.mp4")},
		},
		Playlists: []*elastictranscoder.CreateJobPlaylist{
			{
				Format: aws.String("HLSv3"),
				Name:   aws.String("job-123/hls/video"),
				OutputKeys: []*string{
					aws.String("job-123/hls/output_hls_360p/video"),
					aws.String("job-123/hls/output_hls_480p/video"),
					aws.String("job-123/hls/output_hls_720p/video"),
				},
			},
		},
	}
	if !reflect.DeepEqual(*jobInput, expectedJobInput) {
		t.Errorf("Elastic Transcoder: wrong input\nWant %#v\nGot  %#v", expectedJobInput, *jobInput)
	}
}

func TestAWSTranscodeNormalizedSource(t *testing.T) {
	fakeTranscoder := newFakeElasticTranscoder()
	prov := &awsProvider{
		c: fakeTranscoder,
		config: &config.ElasticTranscoder{
			AccessKeyID:     "AKIA",
			SecretAccessKey: "secret",
			Region:          "sa-east-1",
			PipelineID:      "mypipeline",
		},
	}
	source := "s3://bucketname/some/dir/with/subdir/file.mov"
	outputs := []db.TranscodeOutput{
		{
			FileName: "output_720p.mp4",
			Preset: db.PresetMap{
				Name: "mp4_720p",
				ProviderMapping: map[string]string{
					Name:    "93239832-0001",
					"other": "irrelevant",
				},
				OutputOpts: db.OutputOptions{Extension: "mp4"},
			},
		},
		{
			FileName: "output_1080p.webm",
			Preset: db.PresetMap{
				Name: "webm_1080p",
				ProviderMapping: map[string]string{
					Name:    "93239832-0002",
					"other": "irrelevant",
				},
				OutputOpts: db.OutputOptions{Extension: "webm"},
			},
		},
	}
	jobStatus, err := prov.Transcode(&db.Job{
		ID:              "job-1",
		SourceMedia:     source,
		Outputs:         outputs,
		StreamingParams: db.StreamingParams{},
	})
	if err != nil {
		t.Fatal(err)
	}
	if m, _ := regexp.MatchString(`^job-[a-f0-9]{8}$`, jobStatus.ProviderJobID); !m {
		t.Errorf("Elastic Transcoder: invalid id returned - %q", jobStatus.ProviderJobID)
	}
	if jobStatus.Status != provider.StatusQueued {
		t.Errorf("Elastic Transcoder: wrong status returned. Want queued. Got %v", jobStatus.Status)
	}
	if jobStatus.ProviderName != Name {
		t.Errorf("Elastic Transcoder: wrong provider name returned. Want %q. Got %q", Name, jobStatus.ProviderName)
	}

	if len(fakeTranscoder.jobs) != 1 {
		t.Fatal("Did not send any job request to the server.")
	}
	jobInput := fakeTranscoder.jobs[jobStatus.ProviderJobID]

	expectedJobInput := elastictranscoder.CreateJobInput{
		PipelineId: aws.String("mypipeline"),
		Input: &elastictranscoder.JobInput{
			Key: aws.String("some/dir/with/subdir/file.mov"),
			DetectedProperties: &elastictranscoder.DetectedProperties{
				DurationMillis: aws.Int64(120e3),
				FileSize:       aws.Int64(60356779),
				Height:         aws.Int64(1080),
				Width:          aws.Int64(1920),
			},
		},
		Outputs: []*elastictranscoder.CreateJobOutput{
			{PresetId: aws.String("93239832-0001"), Key: aws.String("job-1/output_720p.mp4")},
			{PresetId: aws.String("93239832-0002"), Key: aws.String("job-1/output_1080p.webm")},
		},
	}
	if !reflect.DeepEqual(*jobInput, expectedJobInput) {
		t.Errorf("Elastic Transcoder: wrong input\nWant %#v\nGot  %#v", expectedJobInput, *jobInput)
	}
}

func TestAWSTranscodePresetNotFound(t *testing.T) {
	fakeTranscoder := newFakeElasticTranscoder()
	prov := &awsProvider{
		c: fakeTranscoder,
		config: &config.ElasticTranscoder{
			AccessKeyID:     "AKIA",
			SecretAccessKey: "secret",
			Region:          "sa-east-1",
			PipelineID:      "mypipeline",
		},
	}
	source := "s3://bucketname/some/dir/with/subdir/file.mov"
	outputs := []db.TranscodeOutput{
		{
			FileName: "output_720p.mp4",
			Preset: db.PresetMap{
				Name:            "mp4_720p",
				ProviderMapping: map[string]string{"other": "irrelevant"},
				OutputOpts:      db.OutputOptions{Extension: "mp4"},
			},
		},
	}
	jobStatus, err := prov.Transcode(&db.Job{
		ID:              "job-123",
		SourceMedia:     source,
		Outputs:         outputs,
		StreamingParams: db.StreamingParams{},
	})
	if err != provider.ErrPresetMapNotFound {
		t.Errorf("Wrong error returned. Want %#v. Got %#v", provider.ErrPresetMapNotFound, err)
	}
	if jobStatus != nil {
		t.Errorf("Got unexpected non-nil JobStatus: %#v", jobStatus)
	}
}

func TestAWSTranscodeAWSFailureInAmazon(t *testing.T) {
	prepErr := errors.New("something went wrong")
	fakeTranscoder := newFakeElasticTranscoder()
	fakeTranscoder.prepareFailure("CreateJob", prepErr)
	prov := &awsProvider{
		c: fakeTranscoder,
		config: &config.ElasticTranscoder{
			AccessKeyID:     "AKIA",
			SecretAccessKey: "secret",
			Region:          "sa-east-1",
			PipelineID:      "mypipeline",
		},
	}
	source := "dir/file.mp4"
	jobStatus, err := prov.Transcode(&db.Job{
		ID:              "job-123",
		SourceMedia:     source,
		StreamingParams: db.StreamingParams{},
	})
	if jobStatus != nil {
		t.Errorf("Got unexpected non-nil status: %#v", jobStatus)
	}
	if err != prepErr {
		t.Errorf("Got wrong error. Want %q. Got %q", prepErr.Error(), err.Error())
	}
}

func TestAWSJobStatus(t *testing.T) {
	fakeTranscoder := newFakeElasticTranscoder()
	prov := &awsProvider{
		c: fakeTranscoder,
		config: &config.ElasticTranscoder{
			AccessKeyID:     "AKIA",
			SecretAccessKey: "secret",
			Region:          "sa-east-1",
			PipelineID:      "mypipeline",
		},
	}
	outputs := []db.TranscodeOutput{
		{
			FileName: "output_720p.mp4",
			Preset: db.PresetMap{
				Name: "mp4_720p",
				ProviderMapping: map[string]string{
					Name:    "93239832-0001",
					"other": "irrelevant",
				},
				OutputOpts: db.OutputOptions{Extension: "mp4"},
			},
		},
		{
			FileName: "output_720p.webm",
			Preset: db.PresetMap{
				Name: "webm_720p",
				ProviderMapping: map[string]string{
					Name:    "93239832-0002",
					"other": "irrelevant",
				},
				OutputOpts: db.OutputOptions{Extension: "webm"},
			},
		},
		{
			FileName: "hls/output_720p.m3u8",
			Preset: db.PresetMap{
				Name: "hls_720p",
				ProviderMapping: map[string]string{
					Name:    "hls-93239832-0003",
					"other": "irrelevant",
				},
				OutputOpts: db.OutputOptions{Extension: "m3u8"},
			},
		},
	}
	source := "dir/file.mov"
	jobStatus, err := prov.Transcode(&db.Job{
		ID:          "job-123",
		SourceMedia: source,
		Outputs:     outputs,
		StreamingParams: db.StreamingParams{
			PlaylistFileName: "hls/index.m3u8",
			SegmentDuration:  3,
			Protocol:         "hls",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	jobStatus, err = prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: jobStatus.ProviderJobID})
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := provider.JobStatus{
		ProviderJobID: jobStatus.ProviderJobID,
		Status:        provider.StatusFinished,
		StatusMessage: "it's finished!",
		Progress:      100,
		ProviderStatus: map[string]interface{}{
			"outputs": map[string]interface{}{
				"job-123/output_720p.mp4":  "it's finished!",
				"job-123/output_720p.webm": "it's finished!",
				"job-123/hls/output_720p":  "it's finished!",
			},
		},
		SourceInfo: provider.SourceInfo{
			Duration: 120 * time.Second,
			Width:    1920,
			Height:   1080,
		},
		Output: provider.JobOutput{
			Destination: "s3://some bucket/job-123",
			Files: []provider.OutputFile{
				{
					Path:       "s3://some bucket/job-123/output_720p.mp4",
					Container:  "mp4",
					VideoCodec: "H.264",
					Width:      0,
					Height:     720,
				},
				{
					Path:       "s3://some bucket/job-123/output_720p.webm",
					Container:  "webm",
					VideoCodec: "VP8",
					Width:      0,
					Height:     720,
				},
				{
					Path:      "s3://some bucket/job-123/hls/index.m3u8",
					Container: "m3u8",
				},
			},
		},
	}
	if !reflect.DeepEqual(*jobStatus, expectedJobStatus) {
		t.Errorf("Wrong JobStatus\nWant %#v\nGot  %#v", expectedJobStatus, *jobStatus)
	}
}

func TestAWSJobStatusNoDetectedProperties(t *testing.T) {
	fakeTranscoder := newFakeElasticTranscoder()
	prov := &awsProvider{
		c: fakeTranscoder,
		config: &config.ElasticTranscoder{
			AccessKeyID:     "AKIA",
			SecretAccessKey: "secret",
			Region:          "sa-east-1",
			PipelineID:      "mypipeline",
		},
	}
	outputs := []db.TranscodeOutput{
		{
			FileName: "output_720p.mp4",
			Preset: db.PresetMap{
				Name: "mp4_720p",
				ProviderMapping: map[string]string{
					Name:    "93239832-0001",
					"other": "irrelevant",
				},
				OutputOpts: db.OutputOptions{Extension: "mp4"},
			},
		},
		{
			FileName: "output_720p.webm",
			Preset: db.PresetMap{
				Name: "webm_720p",
				ProviderMapping: map[string]string{
					Name:    "93239832-0002",
					"other": "irrelevant",
				},
				OutputOpts: db.OutputOptions{Extension: "webm"},
			},
		},
	}
	jobStatus, err := prov.Transcode(&db.Job{
		ID:              "job-123",
		SourceMedia:     "dir/file.mov",
		Outputs:         outputs,
		StreamingParams: db.StreamingParams{},
	})
	if err != nil {
		t.Fatal(err)
	}
	fakeTranscoder.jobs[jobStatus.ProviderJobID].Input.DetectedProperties = nil
	jobStatus, err = prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: jobStatus.ProviderJobID})
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := provider.JobStatus{
		ProviderJobID: jobStatus.ProviderJobID,
		Status:        provider.StatusFinished,
		StatusMessage: "it's finished!",
		Progress:      100,
		ProviderStatus: map[string]interface{}{
			"outputs": map[string]interface{}{
				"job-123/output_720p.mp4":  "it's finished!",
				"job-123/output_720p.webm": "it's finished!",
			},
		},
		Output: provider.JobOutput{
			Destination: "s3://some bucket/job-123",
			Files: []provider.OutputFile{
				{
					Path:       "s3://some bucket/job-123/output_720p.mp4",
					Container:  "mp4",
					VideoCodec: "H.264",
					Width:      0,
					Height:     720,
				},
				{
					Path:       "s3://some bucket/job-123/output_720p.webm",
					Container:  "webm",
					VideoCodec: "VP8",
					Width:      0,
					Height:     720,
				},
			},
		},
	}
	if !reflect.DeepEqual(*jobStatus, expectedJobStatus) {
		t.Errorf("Wrong JobStatus\nWant %#v\nGot  %#v", expectedJobStatus, *jobStatus)
	}
}

func TestAWSCreatePreset(t *testing.T) {
	fakeTranscoder := newFakeElasticTranscoder()
	prov := &awsProvider{
		c: fakeTranscoder,
		config: &config.ElasticTranscoder{
			AccessKeyID:     "AKIA",
			SecretAccessKey: "secret",
			Region:          "sa-east-1",
			PipelineID:      "mypipeline",
		},
	}

	inputPreset := db.Preset{
		Name:        "preset_name",
		Description: "description here",
		Container:   "mp4",
		RateControl: "VBR",
		Video: db.VideoPreset{
			Profile:       "Main",
			ProfileLevel:  "3.1",
			Height:        "720",
			Codec:         "h264",
			Bitrate:       "2500000",
			GopSize:       "90",
			GopMode:       "fixed",
			InterlaceMode: "progressive",
		},
		Audio: db.AudioPreset{
			Codec:   "aac",
			Bitrate: "64000",
		},
	}

	presetID, _ := prov.CreatePreset(inputPreset)

	if !reflect.DeepEqual(presetID, "preset_name-abc123") {
		t.Errorf("CreatePreset: want %s. Got %s", presetID, "preset_name-abc123")
	}
}

func TestCreateVideoPreset(t *testing.T) {
	fakeTranscoder := newFakeElasticTranscoder()
	prov := &awsProvider{
		c: fakeTranscoder,
		config: &config.ElasticTranscoder{
			AccessKeyID:     "AKIA",
			SecretAccessKey: "secret",
			Region:          "sa-east-1",
			PipelineID:      "mypipeline",
		},
	}
	var tests = []struct {
		givenTestCase       string
		givenPreset         db.Preset
		expectedVideoParams *elastictranscoder.VideoParameters
	}{
		{
			"H.264 preset",
			db.Preset{
				Container: "m3u8",
				Video: db.VideoPreset{
					Profile:      "Main",
					ProfileLevel: "3.1",
					Codec:        "h264",
				},
			},
			&elastictranscoder.VideoParameters{
				BitRate: aws.String("0"),
				Codec:   aws.String("H.264"),
				CodecOptions: map[string]*string{
					"MaxReferenceFrames": aws.String("2"),
					"Profile":            aws.String("main"),
					"Level":              aws.String("3.1"),
				},
				DisplayAspectRatio: aws.String("auto"),
				FrameRate:          aws.String("auto"),
				KeyframesMaxDist:   aws.String(""),
				MaxHeight:          aws.String("auto"),
				MaxWidth:           aws.String("auto"),
				PaddingPolicy:      aws.String("Pad"),
				SizingPolicy:       aws.String("Fill"),
			},
		},
		{
			"WEBM vp8 preset",
			db.Preset{
				Container: "webm",
				Video: db.VideoPreset{
					Codec:   "vp8",
					GopSize: "90",
				},
			},
			&elastictranscoder.VideoParameters{
				BitRate: aws.String("0"),
				Codec:   aws.String("vp8"),
				CodecOptions: map[string]*string{
					"Profile": aws.String("0"),
				},
				DisplayAspectRatio: aws.String("auto"),
				FrameRate:          aws.String("auto"),
				KeyframesMaxDist:   aws.String("90"),
				MaxHeight:          aws.String("auto"),
				MaxWidth:           aws.String("auto"),
				PaddingPolicy:      aws.String("Pad"),
				SizingPolicy:       aws.String("Fill"),
			},
		},
		{
			"WEBM vp9 preset",
			db.Preset{
				Container: "webm",
				Video: db.VideoPreset{
					Codec:   "vp9",
					GopSize: "90",
				},
			},
			&elastictranscoder.VideoParameters{
				BitRate: aws.String("0"),
				Codec:   aws.String("vp9"),
				CodecOptions: map[string]*string{
					"Profile": aws.String("0"),
				},
				DisplayAspectRatio: aws.String("auto"),
				FrameRate:          aws.String("auto"),
				KeyframesMaxDist:   aws.String("90"),
				MaxHeight:          aws.String("auto"),
				MaxWidth:           aws.String("auto"),
				PaddingPolicy:      aws.String("Pad"),
				SizingPolicy:       aws.String("Fill"),
			},
		},
		{
			"MP4 preset",
			db.Preset{
				Container: "mp4",
				Video: db.VideoPreset{
					Profile:      "Main",
					ProfileLevel: "3.1",
					Codec:        "h264",
					GopSize:      "90",
				},
			},
			&elastictranscoder.VideoParameters{
				BitRate: aws.String("0"),
				Codec:   aws.String("H.264"),
				CodecOptions: map[string]*string{
					"MaxReferenceFrames": aws.String("2"),
					"Profile":            aws.String("main"),
					"Level":              aws.String("3.1"),
				},
				DisplayAspectRatio: aws.String("auto"),
				FrameRate:          aws.String("auto"),
				KeyframesMaxDist:   aws.String("90"),
				MaxHeight:          aws.String("auto"),
				MaxWidth:           aws.String("auto"),
				PaddingPolicy:      aws.String("Pad"),
				SizingPolicy:       aws.String("Fill"),
			},
		},
	}
	for _, test := range tests {
		videoParams := prov.createVideoPreset(test.givenPreset)
		if !reflect.DeepEqual(test.expectedVideoParams, videoParams) {
			t.Errorf("%s: CreateVideoPreset: want %s. Got %s", test.givenTestCase, test.expectedVideoParams, videoParams)
			pretty.Fdiff(os.Stderr, videoParams, test.expectedVideoParams)
		}
	}
}

func TestCreateAudioPreset(t *testing.T) {
	fakeTranscoder := newFakeElasticTranscoder()
	prov := &awsProvider{
		c: fakeTranscoder,
		config: &config.ElasticTranscoder{
			AccessKeyID:     "AKIA",
			SecretAccessKey: "secret",
			Region:          "sa-east-1",
			PipelineID:      "mypipeline",
		},
	}
	var tests = []struct {
		givenTestCase       string
		givenPreset         db.Preset
		expectedAudioParams *elastictranscoder.AudioParameters
	}{
		{
			"AAC preset",
			db.Preset{
				Audio: db.AudioPreset{
					Codec: "aac",
				},
			},
			&elastictranscoder.AudioParameters{
				BitRate:    aws.String("0"),
				Channels:   aws.String("auto"),
				Codec:      aws.String("AAC"),
				SampleRate: aws.String("auto"),
			},
		},
		{
			"libvorbis preset",
			db.Preset{
				Audio: db.AudioPreset{
					Codec: "libvorbis",
				},
			},
			&elastictranscoder.AudioParameters{
				BitRate:    aws.String("0"),
				Channels:   aws.String("auto"),
				Codec:      aws.String("vorbis"),
				SampleRate: aws.String("auto"),
			},
		},
	}
	for _, test := range tests {
		audioParams := prov.createAudioPreset(test.givenPreset)
		if !reflect.DeepEqual(test.expectedAudioParams, audioParams) {
			t.Errorf("%s: CreateAudioPreset: want %s. Got %s", test.givenTestCase, test.expectedAudioParams, audioParams)
			pretty.Fdiff(os.Stderr, audioParams, test.expectedAudioParams)
		}
	}
}

func TestAWSJobStatusNotFound(t *testing.T) {
	fakeTranscoder := newFakeElasticTranscoder()
	provider := &awsProvider{
		c: fakeTranscoder,
		config: &config.ElasticTranscoder{
			AccessKeyID:     "AKIA",
			SecretAccessKey: "secret",
			Region:          "sa-east-1",
			PipelineID:      "mypipeline",
		},
	}
	jobStatus, err := provider.JobStatus(&db.Job{ProviderJobID: "idk"})
	if err == nil {
		t.Fatal("Got unexpected <nil> error")
	}
	expectedErrMsg := "job not found"
	if err.Error() != expectedErrMsg {
		t.Errorf("Got wrong error message. Want %q. Got %q", expectedErrMsg, err.Error())
	}
	if jobStatus != nil {
		t.Errorf("Got unexpected non-nil JobStatus: %#v", jobStatus)
	}
}

func TestAWSJobStatusInternalError(t *testing.T) {
	prepErr := errors.New("failed to get job status")
	fakeTranscoder := newFakeElasticTranscoder()
	fakeTranscoder.prepareFailure("ReadJob", prepErr)
	provider := &awsProvider{
		c: fakeTranscoder,
		config: &config.ElasticTranscoder{
			AccessKeyID:     "AKIA",
			SecretAccessKey: "secret",
			Region:          "sa-east-1",
			PipelineID:      "mypipeline",
		},
	}
	jobStatus, err := provider.JobStatus(&db.Job{ProviderJobID: "idk"})
	if jobStatus != nil {
		t.Errorf("Got unexpected non-nil JobStatus: %#v", jobStatus)
	}
	if err != prepErr {
		t.Errorf("Got wrong error. Want %q. Got %q", prepErr.Error(), err.Error())
	}
}

func TestAWSStatusMap(t *testing.T) {
	var tests = []struct {
		input  string
		output provider.Status
	}{
		{"Submitted", provider.StatusQueued},
		{"Progressing", provider.StatusStarted},
		{"Canceled", provider.StatusCanceled},
		{"Error", provider.StatusFailed},
		{"Complete", provider.StatusFinished},
		{"unknown", provider.StatusFailed},
	}
	var prov awsProvider
	for _, test := range tests {
		result := prov.statusMap(test.input)
		if result != test.output {
			t.Errorf("statusMap(%q): wrong result. Want %q. Got %q", test.input, test.output, result)
		}
	}
}

func TestCancelJob(t *testing.T) {
	fakeTranscoder := newFakeElasticTranscoder()
	prov := &awsProvider{
		c: fakeTranscoder,
		config: &config.ElasticTranscoder{
			AccessKeyID:     "AKIA",
			SecretAccessKey: "secret",
			Region:          "sa-east-1",
			PipelineID:      "mypipeline",
		},
	}
	err := prov.CancelJob("idk")
	if err != nil {
		t.Fatal(err)
	}
	if id := aws.StringValue(fakeTranscoder.canceledJobs[0].Id); id != "idk" {
		t.Errorf("wrong job canceled. Want %q. Got %q", "idk", id)
	}
}

func TestCancelJobInternalError(t *testing.T) {
	prepErr := errors.New("failed to cancel job")
	fakeTranscoder := newFakeElasticTranscoder()
	fakeTranscoder.prepareFailure("CancelJob", prepErr)
	provider := &awsProvider{
		c: fakeTranscoder,
		config: &config.ElasticTranscoder{
			AccessKeyID:     "AKIA",
			SecretAccessKey: "secret",
			Region:          "sa-east-1",
			PipelineID:      "mypipeline",
		},
	}
	err := provider.CancelJob("idk")
	if err != prepErr {
		t.Errorf("wrong error returned.\nWant %#v\nGot  %#v", prepErr, err)
	}
}

func TestHealthcheck(t *testing.T) {
	fakeTranscoder := newFakeElasticTranscoder()
	provider := &awsProvider{
		c: fakeTranscoder,
		config: &config.ElasticTranscoder{
			AccessKeyID:     "AKIA",
			SecretAccessKey: "secret",
			Region:          "sa-east-1",
			PipelineID:      "mypipeline",
		},
	}
	err := provider.Healthcheck()
	if err != nil {
		t.Fatal(err)
	}
}

func TestHealthcheckFailure(t *testing.T) {
	prepErr := errors.New("something went wrong")
	fakeTranscoder := newFakeElasticTranscoder()
	fakeTranscoder.prepareFailure("ReadPipeline", prepErr)
	provider := &awsProvider{
		c: fakeTranscoder,
		config: &config.ElasticTranscoder{
			AccessKeyID:     "AKIA",
			SecretAccessKey: "secret",
			Region:          "sa-east-1",
			PipelineID:      "mypipeline",
		},
	}
	err := provider.Healthcheck()
	if err != prepErr {
		t.Errorf("Wrong error returned. Want %#v.Got %#v", prepErr, err)
	}
}

func TestCapabilities(t *testing.T) {
	var prov awsProvider
	expected := provider.Capabilities{
		InputFormats:  []string{"h264"},
		OutputFormats: []string{"mp4", "hls", "webm"},
		Destinations:  []string{"s3"},
	}
	cap := prov.Capabilities()
	if !reflect.DeepEqual(cap, expected) {
		t.Errorf("Capabilities: want %#v. Got %#v", expected, cap)
	}
}
