package hybrik

import (
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/provider"
	hwrapper "github.com/hybrik/hybrik-sdk-go"
)

const (
	// Name describes the name of the transcoder
	Name          = "hybrik"
	queued        = "queued"
	active        = "active"
	completed     = "completed"
	failed        = "failed"
	activeRunning = "running"
	activeWaiting = "waiting"
	hls           = "hls"
)

var (
	// ErrBitrateNan is an error returned when the bitrate field of db.Preset is not a valid number
	ErrBitrateNan = fmt.Errorf("bitrate not a number")

	// ErrPresetOutputMatch represents an error in the hybrik encoding-wrapper provider.
	ErrPresetOutputMatch = fmt.Errorf("preset retrieved does not map to hybrik.Preset struct")

	// ErrVideoWidthNan is an error returned when the preset video width of db.Preset is not a valid number
	ErrVideoWidthNan = fmt.Errorf("preset video width not a number")
	// ErrVideoHeightNan is an error returned when the preset video height of db.Preset is not a valid number
	ErrVideoHeightNan = fmt.Errorf("preset video height not a number")

	// ErrUnsupportedContainer is returned when the container format is not present in the provider's capabilities list
	ErrUnsupportedContainer = fmt.Errorf("container format unsupported. Hybrik provider capabilities may need to be updated")
)

func init() {
	provider.Register(Name, hybrikTranscoderFactory)
}

type hybrikProvider struct {
	c      hwrapper.ClientInterface
	config *config.Hybrik
}

func (hp hybrikProvider) String() string {
	return "Hybrik"
}

func hybrikTranscoderFactory(cfg *config.Config) (provider.TranscodingProvider, error) {
	api, err := hwrapper.NewClient(hwrapper.Config{
		URL:            cfg.Hybrik.URL,
		ComplianceDate: cfg.Hybrik.ComplianceDate,
		OAPIKey:        cfg.Hybrik.OAPIKey,
		OAPISecret:     cfg.Hybrik.OAPISecret,
		AuthKey:        cfg.Hybrik.AuthKey,
		AuthSecret:     cfg.Hybrik.AuthSecret,
	})
	if err != nil {
		return &hybrikProvider{}, err
	}

	return &hybrikProvider{
		c:      api,
		config: cfg.Hybrik,
	}, nil
}

func (hp *hybrikProvider) Transcode(job *db.Job) (*provider.JobStatus, error) {
	cj, err := hp.presetsToTranscodeJob(job)
	if err != nil {
		return &provider.JobStatus{}, err
	}

	id, err := hp.c.QueueJob(cj)
	if err != nil {
		return &provider.JobStatus{}, err
	}

	return &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: id,
		Status:        provider.StatusQueued,
	}, nil
}

func (hp *hybrikProvider) mountTranscodeElement(elementID, id, outputFilename, destination string, duration uint, preset hwrapper.Preset) (hwrapper.Element, error) {
	var e hwrapper.Element
	var subLocation *hwrapper.TranscodeLocation

	// outputFilename can be "test.mp4", or "subfolder1/subfodler2/test.mp4"
	// Handling accordingly
	subPath := path.Dir(outputFilename)
	outputFilePattern := path.Base(outputFilename)
	if subPath != "." && subPath != "/" {
		subLocation = &hwrapper.TranscodeLocation{
			StorageProvider: "relative",
			Path:            subPath,
		}
	}

	// create the transcode element
	e = hwrapper.Element{
		UID:  "transcode_task" + elementID,
		Kind: "transcode",
		Task: &hwrapper.ElementTaskOptions{
			Name: "Transcode - " + preset.Name,
		},
		Preset: &hwrapper.TranscodePreset{
			Key: preset.Name,
		},
		Payload: hwrapper.LocationTargetPayload{
			Location: hwrapper.TranscodeLocation{
				StorageProvider: "s3",
				Path:            fmt.Sprintf("%s/j%s", destination, id),
			},
			Targets: []hwrapper.TranscodeLocationTarget{
				{
					Location:    subLocation,
					FilePattern: outputFilePattern,
					Container: hwrapper.TranscodeTargetContainer{
						SegmentDuration: duration,
					},
				},
			},
		},
	}

	return e, nil
}

type presetResult struct {
	presetID string
	preset   interface{}
}

func makeGetPresetRequest(hp *hybrikProvider, presetID string, ch chan *presetResult) {
	presetOutput, err := hp.GetPreset(presetID)
	result := new(presetResult)
	result.presetID = presetID
	if err != nil {
		result.preset = err
		ch <- result
	} else {
		result.preset = presetOutput
		ch <- result
	}
}

func (hp *hybrikProvider) presetsToTranscodeJob(job *db.Job) (string, error) {
	elements := []hwrapper.Element{}
	var hlsElementIds []int

	// create a source element
	sourceElement := hwrapper.Element{
		UID:  "source_file",
		Kind: "source",
		Payload: hwrapper.ElementPayload{
			Kind: "asset_url",
			Payload: hwrapper.AssetPayload{
				StorageProvider: "s3",
				URL:             job.SourceMedia,
			},
		},
	}

	elements = append(elements, sourceElement)

	presetCh := make(chan *presetResult)
	presets := make(map[string]interface{})

	for _, output := range job.Outputs {
		presetID, ok := output.Preset.ProviderMapping[Name]
		if !ok {
			return "", provider.ErrPresetMapNotFound
		}

		if _, ok := presets[presetID]; ok {
			continue
		}

		presets[presetID] = nil

		go makeGetPresetRequest(hp, presetID, presetCh)
	}

	for i := 0; i < len(presets); i++ {
		res := <-presetCh
		err, isErr := res.preset.(error)
		if isErr {
			return "", fmt.Errorf("Error getting preset info: %s", err.Error())
		}

		presets[res.presetID] = res.preset
	}

	// create transcode elements for each target
	// TODO: This can be optimized further with regards to combining tasks so that they run in the same machine. Requires some discussion
	elementID := 0
	for _, output := range job.Outputs {
		presetID, ok := output.Preset.ProviderMapping[Name]
		if !ok {
			return "", provider.ErrPresetMapNotFound
		}

		presetOutput, ok := presets[presetID]
		if !ok {
			return "", fmt.Errorf("Hybrik preset not found in preset results")
		}

		preset, ok := presetOutput.(hwrapper.Preset)
		if !ok {
			return "", ErrPresetOutputMatch
		}

		var segmentDur uint
		// track the hls outputs so we can later connect them to a manifest creator task
		if len(preset.Payload.Targets) > 0 && preset.Payload.Targets[0].Container.Kind == hls {
			hlsElementIds = append(hlsElementIds, elementID)
			segmentDur = job.StreamingParams.SegmentDuration
		}

		e, err := hp.mountTranscodeElement(strconv.Itoa(elementID), job.ID, output.FileName, hp.config.Destination, segmentDur, preset)
		if err != nil {
			return "", err
		}

		elements = append(elements, e)

		elementID++
	}

	// connect the source element to each of the transcode elements
	var transcodeSuccessConnections []hwrapper.ToSuccess
	for i := 0; i < elementID; i++ {
		transcodeSuccessConnections = append(transcodeSuccessConnections, hwrapper.ToSuccess{Element: "transcode_task" + strconv.Itoa(i)})
	}

	// create the full job structure
	cj := hwrapper.CreateJob{
		Name: fmt.Sprintf("Job %s [%s]", job.ID, path.Base(job.SourceMedia)),
		Payload: hwrapper.CreateJobPayload{
			Elements: elements,
			Connections: []hwrapper.Connection{
				{
					From: []hwrapper.ConnectionFrom{
						{
							Element: "source_file",
						},
					},
					To: hwrapper.ConnectionTo{
						Success: transcodeSuccessConnections,
					},
				},
			},
		},
	}

	// check if we need to add a master manifest task element
	if job.StreamingParams.Protocol == hls {
		manifestOutputDir := fmt.Sprintf("%s/j%s", hp.config.Destination, job.ID)
		manifestSubDir := path.Dir(job.StreamingParams.PlaylistFileName)
		manifestFilePattern := path.Base(job.StreamingParams.PlaylistFileName)

		if manifestSubDir != "." && manifestSubDir != "/" {
			manifestOutputDir = path.Join(manifestOutputDir, manifestSubDir)
		}

		manifestElement := hwrapper.Element{
			UID:  "manifest_creator",
			Kind: "manifest_creator",
			Payload: hwrapper.ManifestCreatorPayload{
				Location: hwrapper.TranscodeLocation{
					StorageProvider: "s3",
					Path:            manifestOutputDir,
				},
				FilePattern: manifestFilePattern,
				Kind:        hls,
			},
		}

		cj.Payload.Elements = append(cj.Payload.Elements, manifestElement)

		var manifestFromConnections []hwrapper.ConnectionFrom
		for _, hlsElementID := range hlsElementIds {
			manifestFromConnections = append(manifestFromConnections, hwrapper.ConnectionFrom{Element: "transcode_task" + strconv.Itoa(hlsElementID)})
		}

		cj.Payload.Connections = append(cj.Payload.Connections,
			hwrapper.Connection{
				From: manifestFromConnections,
				To: hwrapper.ConnectionTo{
					Success: []hwrapper.ToSuccess{
						{Element: "manifest_creator"},
					},
				},
			},
		)

	}

	resp, err := json.Marshal(cj)
	if err != nil {
		return "", err
	}

	return string(resp), nil
}

func (hp *hybrikProvider) JobStatus(job *db.Job) (*provider.JobStatus, error) {
	ji, err := hp.c.GetJobInfo(job.ProviderJobID)
	if err != nil {
		return &provider.JobStatus{}, err
	}

	var status provider.Status
	switch ji.Status {
	case active:
		fallthrough
	case activeRunning:
		fallthrough
	case activeWaiting:
		status = provider.StatusStarted
	case queued:
		status = provider.StatusQueued
	case completed:
		status = provider.StatusFinished
	case failed:
		status = provider.StatusFailed
	}

	return &provider.JobStatus{
		ProviderJobID: job.ProviderJobID,
		ProviderName:  hp.String(),
		Progress:      float64(ji.Progress),
		Status:        status,
	}, nil
}

func (hp *hybrikProvider) CancelJob(id string) error {
	return hp.c.StopJob(id)
}

func (hp *hybrikProvider) CreatePreset(preset db.Preset) (string, error) {
	var minGOPFrames, maxGOPFrames, gopSize int

	gopSize, err := strconv.Atoi(preset.Video.GopSize)
	if err != nil {
		return "", err
	}

	if preset.Video.GopMode == "fixed" {
		minGOPFrames = gopSize
		maxGOPFrames = gopSize
	} else {
		maxGOPFrames = gopSize
	}

	container := ""
	for _, c := range hp.Capabilities().OutputFormats {
		if preset.Container == c || (preset.Container == "m3u8" && c == hls) {
			container = c
		}
	}

	if container == "" {
		return "", ErrUnsupportedContainer
	}

	bitrate, err := strconv.Atoi(preset.Video.Bitrate)
	if err != nil {
		return "", ErrBitrateNan
	}

	audioBitrate, err := strconv.Atoi(preset.Audio.Bitrate)
	if err != nil {
		return "", ErrBitrateNan
	}

	var videoWidth *int
	var videoHeight *int

	if preset.Video.Width != "" {
		var presetWidth int
		presetWidth, err = strconv.Atoi(preset.Video.Width)
		if err != nil {
			return "", ErrVideoWidthNan
		}
		videoWidth = &presetWidth
	}

	if preset.Video.Height != "" {
		var presetHeight int
		presetHeight, err = strconv.Atoi(preset.Video.Height)
		if err != nil {
			return "", ErrVideoHeightNan
		}
		videoHeight = &presetHeight
	}

	videoProfile := strings.ToLower(preset.Video.Profile)
	videoLevel := preset.Video.ProfileLevel

	// TODO: Understand video-transcoding-api profile + level settings in relation to vp8
	// For now, we will omit and leave to encoder defaults
	if preset.Video.Codec == "vp8" {
		videoProfile = ""
		videoLevel = ""
	}

	p := hwrapper.Preset{
		Key:         preset.Name,
		Name:        preset.Name,
		Description: preset.Description,
		Kind:        "transcode",
		Path:        hp.config.PresetPath,
		Payload: hwrapper.PresetPayload{
			Targets: []hwrapper.PresetTarget{
				{
					FilePattern: "",
					Container:   hwrapper.TranscodeContainer{Kind: container},
					Video: hwrapper.VideoTarget{
						Width:         videoWidth,
						Height:        videoHeight,
						Codec:         preset.Video.Codec,
						BitrateKb:     bitrate / 1000,
						MinGOPFrames:  minGOPFrames,
						MaxGOPFrames:  maxGOPFrames,
						Profile:       videoProfile,
						Level:         videoLevel,
						InterlaceMode: preset.Video.InterlaceMode,
					},
					Audio: []hwrapper.AudioTarget{
						{
							Codec:     preset.Audio.Codec,
							BitrateKb: audioBitrate / 1000,
						},
					},
					ExistingFiles: "replace",
					UID:           "target",
				},
			},
		},
	}

	resultPreset, err := hp.c.CreatePreset(p)
	if err != nil {
		return "", err
	}

	return resultPreset.Name, nil
}

func (hp *hybrikProvider) DeletePreset(presetID string) error {
	return hp.c.DeletePreset(presetID)
}

func (hp *hybrikProvider) GetPreset(presetID string) (interface{}, error) {
	preset, err := hp.c.GetPreset(presetID)
	if err != nil {
		return nil, err
	}

	return preset, nil
}

// Healthcheck should return nil if the provider is currently available
// for transcoding videos, otherwise it should return an error
// explaining what's going on.
func (hp *hybrikProvider) Healthcheck() error {
	// For now, just call list jobs. If this errors, then we can consider the service unhealthy
	_, err := hp.c.CallAPI("GET", "/jobs/info", nil, nil)
	return err
}

// Capabilities describes the capabilities of the provider.
func (hp *hybrikProvider) Capabilities() provider.Capabilities {
	// we can support quite a bit more format wise, but unsure of schema so limiting to known supported video-transcoding-api formats for now...
	return provider.Capabilities{
		InputFormats:  []string{"prores", "h264"},
		OutputFormats: []string{"mp4", "hls", "webm", "mov"},
		Destinations:  []string{"s3"},
	}
}
