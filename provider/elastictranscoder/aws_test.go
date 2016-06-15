package elastictranscoder

import (
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/provider"
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
	provider, err := elasticTranscoderProvider(&cfg)
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
	provider, err := elasticTranscoderProvider(&cfg)
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
		provider, err := elasticTranscoderProvider(&cfg)
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
	presetmaps := []db.PresetMap{
		{
			Name: "mp4_720p",
			ProviderMapping: map[string]string{
				Name:    "93239832-0001",
				"other": "irrelevant",
			},
			OutputOpts: db.OutputOptions{Extension: "mp4"},
		},
		{
			Name: "webm_720p",
			ProviderMapping: map[string]string{
				Name:    "93239832-0002",
				"other": "irrelevant",
			},
			OutputOpts: db.OutputOptions{Extension: "webm"},
		},
		{
			Name: "mov_1080p",
			ProviderMapping: map[string]string{
				Name:    "93239832-0003",
				"other": "irrelevant",
			},
			OutputOpts: db.OutputOptions{Extension: "mov"},
		},
	}

	transcodeProfile := provider.TranscodeProfile{
		SourceMedia:     source,
		Presets:         presetmaps,
		StreamingParams: provider.StreamingParams{},
	}
	jobStatus, err := prov.Transcode(&db.Job{ID: "job-123"}, transcodeProfile)
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
		Input:      &elastictranscoder.JobInput{Key: aws.String(source)},
		Outputs: []*elastictranscoder.CreateJobOutput{
			{PresetId: aws.String("93239832-0001"), Key: aws.String("job-123/dir/mp4_720p/file.mp4")},
			{PresetId: aws.String("93239832-0002"), Key: aws.String("job-123/dir/webm_720p/file.webm")},
			{PresetId: aws.String("93239832-0003"), Key: aws.String("job-123/dir/mov_1080p/file.mov")},
		},
	}
	if !reflect.DeepEqual(*jobInput, expectedJobInput) {
		t.Errorf("Elastic Transcoder: wrong input. Want %#v. Got %#v.", expectedJobInput, *jobInput)
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
	presets := []db.PresetMap{
		{
			Name: "hls_360p",
			ProviderMapping: map[string]string{
				Name:    "93239832-0001-hls",
				"other": "irrelevant",
			},
			OutputOpts: db.OutputOptions{Extension: "hls"},
		},

		{
			Name: "hls_480p",
			ProviderMapping: map[string]string{
				Name:    "93239832-0002-hls",
				"other": "irrelevant",
			},
			OutputOpts: db.OutputOptions{Extension: "hls"},
		},
		{
			Name: "hls_720p",
			ProviderMapping: map[string]string{
				Name:    "93239832-0003-hls",
				"other": "irrelevant",
			},
			OutputOpts: db.OutputOptions{Extension: "hls"},
		},
	}

	transcodeProfile := provider.TranscodeProfile{
		SourceMedia:     source,
		Presets:         presets,
		StreamingParams: provider.StreamingParams{Protocol: "asdf", SegmentDuration: 3},
	}
	jobStatus, err := prov.Transcode(&db.Job{ID: "job-123"}, transcodeProfile)
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
		Input:      &elastictranscoder.JobInput{Key: aws.String(source)},
		Outputs: []*elastictranscoder.CreateJobOutput{
			{PresetId: aws.String("93239832-0001-hls"), Key: aws.String("job-123/dir/file/hls_360p/video"), SegmentDuration: aws.String("3")},
			{PresetId: aws.String("93239832-0002-hls"), Key: aws.String("job-123/dir/file/hls_480p/video"), SegmentDuration: aws.String("3")},
			{PresetId: aws.String("93239832-0003-hls"), Key: aws.String("job-123/dir/file/hls_720p/video"), SegmentDuration: aws.String("3")},
		},
		Playlists: []*elastictranscoder.CreateJobPlaylist{
			{
				Format: aws.String("HLSv3"),
				Name:   aws.String("job-123/dir/file/master"),
				OutputKeys: []*string{
					aws.String("job-123/dir/file/hls_360p/video"),
					aws.String("job-123/dir/file/hls_480p/video"),
					aws.String("job-123/dir/file/hls_720p/video"),
				},
			},
		},
	}
	if !reflect.DeepEqual(*jobInput, expectedJobInput) {
		t.Errorf("Elastic Transcoder: wrong input. Want %#v. Got %#v.", expectedJobInput, *jobInput)
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
	presets := []db.PresetMap{
		{
			Name: "hls_360p",
			ProviderMapping: map[string]string{
				Name:    "93239832-0001-hls",
				"other": "irrelevant",
			},
			OutputOpts: db.OutputOptions{Extension: "hls"},
		},

		{
			Name: "hls_480p",
			ProviderMapping: map[string]string{
				Name:    "93239832-0002-hls",
				"other": "irrelevant",
			},
			OutputOpts: db.OutputOptions{Extension: "hls"},
		},
		{
			Name: "hls_720p",
			ProviderMapping: map[string]string{
				Name:    "93239832-0003-hls",
				"other": "irrelevant",
			},
			OutputOpts: db.OutputOptions{Extension: "hls"},
		},
		{
			Name: "mp4_720p",
			ProviderMapping: map[string]string{
				Name:    "93239832-0004",
				"other": "irrelevant",
			},
			OutputOpts: db.OutputOptions{Extension: "mp4"},
		},
	}

	transcodeProfile := provider.TranscodeProfile{
		SourceMedia:     source,
		Presets:         presets,
		StreamingParams: provider.StreamingParams{Protocol: "asdf", SegmentDuration: 3},
	}
	jobStatus, err := prov.Transcode(&db.Job{ID: "job-123"}, transcodeProfile)
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
		Input:      &elastictranscoder.JobInput{Key: aws.String(source)},
		Outputs: []*elastictranscoder.CreateJobOutput{
			{PresetId: aws.String("93239832-0001-hls"), Key: aws.String("job-123/dir/file/hls_360p/video"), SegmentDuration: aws.String("3")},
			{PresetId: aws.String("93239832-0002-hls"), Key: aws.String("job-123/dir/file/hls_480p/video"), SegmentDuration: aws.String("3")},
			{PresetId: aws.String("93239832-0003-hls"), Key: aws.String("job-123/dir/file/hls_720p/video"), SegmentDuration: aws.String("3")},
			{PresetId: aws.String("93239832-0004"), Key: aws.String("job-123/dir/mp4_720p/file.mp4")},
		},
		Playlists: []*elastictranscoder.CreateJobPlaylist{
			{
				Format: aws.String("HLSv3"),
				Name:   aws.String("job-123/dir/file/master"),
				OutputKeys: []*string{
					aws.String("job-123/dir/file/hls_360p/video"),
					aws.String("job-123/dir/file/hls_480p/video"),
					aws.String("job-123/dir/file/hls_720p/video"),
				},
			},
		},
	}
	if !reflect.DeepEqual(*jobInput, expectedJobInput) {
		t.Errorf("Elastic Transcoder: wrong input. Want %#v. Got %#v.", expectedJobInput, *jobInput)
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
	presets := []db.PresetMap{
		{
			Name: "mp4_720p",
			ProviderMapping: map[string]string{
				Name:    "93239832-0001",
				"other": "irrelevant",
			},
			OutputOpts: db.OutputOptions{Extension: "mp4"},
		},
		{
			Name: "hls_1080p",
			ProviderMapping: map[string]string{
				Name:    "93239832-0002",
				"other": "irrelevant",
			},
			OutputOpts: db.OutputOptions{Extension: "ts"},
		},
	}
	transcodeProfile := provider.TranscodeProfile{
		SourceMedia:     source,
		Presets:         presets,
		StreamingParams: provider.StreamingParams{},
	}
	jobStatus, err := prov.Transcode(&db.Job{ID: "job-1"}, transcodeProfile)
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
		Input:      &elastictranscoder.JobInput{Key: aws.String("some/dir/with/subdir/file.mov")},
		Outputs: []*elastictranscoder.CreateJobOutput{
			{PresetId: aws.String("93239832-0001"), Key: aws.String("job-1/some/dir/with/subdir/mp4_720p/file.mp4")},
			{PresetId: aws.String("93239832-0002"), Key: aws.String("job-1/some/dir/with/subdir/hls_1080p/file.ts")},
		},
	}
	if !reflect.DeepEqual(*jobInput, expectedJobInput) {
		t.Errorf("Elastic Transcoder: wrong input. Want %#v. Got %#v.", expectedJobInput, *jobInput)
	}
}

func TestJobStatusOutputDestination(t *testing.T) {
	var tests = []struct {
		job      elastictranscoder.Job
		expected string
	}{
		{
			elastictranscoder.Job{
				PipelineId:      aws.String("mypipeline"),
				OutputKeyPrefix: aws.String("output-prefix"),
				Outputs: []*elastictranscoder.JobOutput{
					{
						// Output for HLS preset
						Key: aws.String("/job-id/some-dir/file-name/preset-name/video"),
					},
					{
						// Output for MP4 preset
						Key: aws.String("/job-id/some-dir/preset-name/file-name.mp4"),
					},
				},
			},
			"s3://some bucket/output-prefix/job-id/some-dir",
		},
	}
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
	for _, test := range tests {
		got, err := prov.getOutputDestination(&test.job)
		if err != nil {
			t.Error(err.Error())
		}
		if got != test.expected {
			t.Errorf("Wrong output destination. Want %q. Got %q", test.expected, got)
		}
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
	presets := []db.PresetMap{
		{
			Name:            "mp4_720p",
			ProviderMapping: map[string]string{"other": "irrelevant"},
			OutputOpts:      db.OutputOptions{Extension: "mp4"},
		},
	}
	transcodeProfile := provider.TranscodeProfile{
		SourceMedia:     source,
		Presets:         presets,
		StreamingParams: provider.StreamingParams{},
	}
	jobStatus, err := prov.Transcode(&db.Job{ID: "job-123"}, transcodeProfile)
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
	transcodeProfile := provider.TranscodeProfile{
		SourceMedia:     source,
		Presets:         nil,
		StreamingParams: provider.StreamingParams{},
	}
	jobStatus, err := prov.Transcode(&db.Job{ID: "job-123"}, transcodeProfile)
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
	presets := []db.PresetMap{
		{
			Name: "mp4_720p",
			ProviderMapping: map[string]string{
				Name:    "93239832-0001",
				"other": "irrelevant",
			},
			OutputOpts: db.OutputOptions{Extension: "mp4"},
		},
		{
			Name: "webm_720p",
			ProviderMapping: map[string]string{
				Name:    "93239832-0002",
				"other": "irrelevant",
			},
			OutputOpts: db.OutputOptions{Extension: "webm"},
		},
	}
	source := "dir/file.mov"
	transcodeProfile := provider.TranscodeProfile{
		SourceMedia:     source,
		Presets:         presets,
		StreamingParams: provider.StreamingParams{},
	}
	jobStatus, err := prov.Transcode(&db.Job{ID: "job-123"}, transcodeProfile)
	if err != nil {
		t.Fatal(err)
	}
	id := jobStatus.ProviderJobID
	jobStatus, err = prov.JobStatus(id)
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := provider.JobStatus{
		ProviderJobID: id,
		Status:        provider.StatusFinished,
		ProviderStatus: map[string]interface{}{
			"outputs": map[string]interface{}{
				"job-123/dir/mp4_720p/file.mp4":   "it's finished!",
				"job-123/dir/webm_720p/file.webm": "it's finished!",
			},
		},
		OutputDestination: "s3://some bucket/job-123",
	}
	if !reflect.DeepEqual(*jobStatus, expectedJobStatus) {
		t.Errorf("Wrong JobStatus. Want %#v. Got %#v.", expectedJobStatus, *jobStatus)
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

	inputPreset := provider.Preset{
		Name:         "preset_name",
		Description:  "description here",
		Container:    "mp4",
		Profile:      "Main",
		ProfileLevel: "3.1",
		RateControl:  "VBR",
		Video: provider.VideoPreset{
			Height:        "720",
			Codec:         "h264",
			Bitrate:       "2500000",
			GopSize:       "90",
			GopMode:       "fixed",
			InterlaceMode: "progressive",
		},
		Audio: provider.AudioPreset{
			Codec:   "aac",
			Bitrate: "64000",
		},
	}

	presetID, _ := prov.CreatePreset(inputPreset)

	if !reflect.DeepEqual(presetID, "preset_name-abc123") {
		t.Errorf("CreatePreset: want %s. Got %s", presetID, "preset_name-abc123")
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
	jobStatus, err := provider.JobStatus("idk")
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
	jobStatus, err := provider.JobStatus("idk")
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
