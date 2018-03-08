// Package encodingcom provides a implementation of the provider that uses the
// Encoding.com API for transcoding media files.
//
// It doesn't expose any public type. In order to use the provider, one must
// import this package and then grab the factory from the provider package:
//
//     import (
//         "github.com/NYTimes/video-transcoding-api/provider"
//         "github.com/NYTimes/video-transcoding-api/provider/encodingcom"
//     )
//
//     func UseProvider() {
//         factory, err := provider.GetProviderFactory(encodingcom.Name)
//         // handle err and use factory to get an instance of the provider.
//     }
package encodingcom

import (
	"errors"
	"fmt"
	"math"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/NYTimes/encoding-wrapper/encodingcom"
	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/provider"
)

// Name is the name used for registering the Encoding.com provider in the
// registry of providers.
const Name = "encodingcom"

var (
	kregexp      = regexp.MustCompile(`000$`)
	s3regexp     = regexp.MustCompile(`^s3://([^/_.]+)/([^?]+)(\?.+)?$`)
	httpS3Regexp = regexp.MustCompile(`https?://([^/_.]+)\.s3\.amazonaws\.com/(.+)$`)
)

var errEncodingComInvalidConfig = provider.InvalidConfigError("missing Encoding.com user id or key. Please define the environment variables ENCODINGCOM_USER_ID and ENCODINGCOM_USER_KEY or set these values in the configuration file")

const hlsOutput = "advanced_hls"

func init() {
	provider.Register(Name, encodingComFactory)
}

type encodingComProvider struct {
	config *config.Config
	client *encodingcom.Client
}

func (e *encodingComProvider) Transcode(job *db.Job) (*provider.JobStatus, error) {
	formats, err := e.presetsToFormats(job)
	if err != nil {
		return nil, fmt.Errorf("Error converting presets to formats on Transcode operation: %s", err.Error())
	}
	resp, err := e.client.AddMedia([]string{e.sourceMedia(job.SourceMedia)}, formats, e.config.EncodingCom.Region)
	if err != nil {
		return nil, fmt.Errorf("Error making AddMedia request for Transcode operation: %s", err.Error())
	}
	return &provider.JobStatus{
		ProviderJobID: resp.MediaID,
		StatusMessage: resp.Message,
		ProviderName:  Name,
	}, nil
}

func (e *encodingComProvider) CreatePreset(preset db.Preset) (string, error) {
	resp, err := e.client.SavePreset(preset.Name, e.presetToFormat(preset))
	if err != nil {
		return "", err
	}
	return resp.SavedPreset, nil
}

func (e *encodingComProvider) sourceMedia(original string) string {
	parts := s3regexp.FindStringSubmatch(original)
	if len(parts) > 0 {
		qs := parts[3]
		if qs == "" {
			qs = "?nocopy"
		}
		return fmt.Sprintf("https://%s.s3.amazonaws.com/%s%s", parts[1], parts[2], qs)
	}
	return original
}

func (e *encodingComProvider) presetToFormat(preset db.Preset) encodingcom.Format {
	falseYesNoBoolean := encodingcom.YesNoBoolean(false)
	format := encodingcom.Format{
		Output:      []string{preset.Container},
		Destination: []string{"ftp://username:password@yourftphost.com/video/encoded/test.flv"},
		TwoPass:     encodingcom.YesNoBoolean(preset.TwoPass),
	}
	if preset.Container == "m3u8" {
		format.Output = []string{hlsOutput}
		format.PackFiles = &falseYesNoBoolean
		format.Stream = e.buildStream(preset)
	} else {
		format.Bitrate = kregexp.ReplaceAllString(preset.Video.Bitrate, "k")
		format.AudioBitrate = kregexp.ReplaceAllString(preset.Audio.Bitrate, "k")
		format.Profile = preset.Video.Profile
		format.Gop = "cgop"
		format.Keyframe = []string{preset.Video.GopSize}
		format.AudioVolume = 100
		format.AudioCodec = e.getNormalizedCodec(preset.Audio.Codec)
		format.VideoCodec = e.getNormalizedCodec(preset.Video.Codec)
		format.Size = e.getSize(preset.Video.Width, preset.Video.Height)
	}
	return format
}

func (e *encodingComProvider) buildStream(preset db.Preset) []encodingcom.Stream {
	stream := encodingcom.Stream{
		Profile:      preset.Video.Profile,
		Keyframe:     preset.Video.GopSize,
		Bitrate:      kregexp.ReplaceAllString(preset.Video.Bitrate, "k"),
		AudioBitrate: kregexp.ReplaceAllString(preset.Audio.Bitrate, "k"),
		AudioVolume:  100,
	}
	stream.AudioCodec = e.getNormalizedCodec(preset.Audio.Codec)
	stream.VideoCodec = e.getNormalizedCodec(preset.Video.Codec)
	stream.Size = e.getSize(preset.Video.Width, preset.Video.Height)

	return []encodingcom.Stream{stream}
}

func (e *encodingComProvider) getSize(width string, height string) string {
	if width == "" {
		width = "0"
	}
	if height == "" {
		height = "0"
	}
	return width + "x" + height
}

func (e *encodingComProvider) getNormalizedCodec(codec string) string {
	audioCodecs := map[string]string{"aac": "dolby_aac", "vorbis": "libvorbis"}
	videoCodecs := map[string]string{"h264": "libx264", "vp8": "libvpx", "vp9": "libvpx-vp9"}
	if c, ok := audioCodecs[codec]; ok {
		return c
	} else if c, ok = videoCodecs[codec]; ok {
		return c
	}
	return codec
}

func (e *encodingComProvider) GetPreset(presetID string) (interface{}, error) {
	preset, err := e.client.GetPreset(presetID)
	if err != nil {
		return nil, err
	}
	return preset, nil
}

func (e *encodingComProvider) DeletePreset(presetID string) error {
	_, err := e.client.DeletePreset(presetID)
	return err
}

func (e *encodingComProvider) getDestinations(jobID, fileName string) []string {
	destination := e.buildDestination(e.config.EncodingCom.Destination, jobID, fileName)
	return []string{destination}
}

func (e *encodingComProvider) buildDestination(baseDestination, jobID, fileName string) string {
	outputPath := strings.TrimRight(baseDestination, "/")
	return outputPath + "/" + path.Join(jobID, fileName)
}

func (e *encodingComProvider) presetsToFormats(job *db.Job) ([]encodingcom.Format, error) {
	streams := []encodingcom.Stream{}
	formats := make([]encodingcom.Format, 0, len(job.Outputs))
	for _, output := range job.Outputs {
		presetName := output.Preset.Name
		presetID, ok := output.Preset.ProviderMapping[Name]
		if !ok {
			return nil, provider.ErrPresetMapNotFound
		}
		presetOutput, err := e.GetPreset(presetID)
		if err != nil {
			return nil, fmt.Errorf("Error getting preset info: %s", err.Error())
		}
		presetStruct := presetOutput.(*encodingcom.Preset)
		if presetStruct.Output == hlsOutput {
			for _, stream := range presetStruct.Format.Stream() {
				stream.SubPath = presetName
				streams = append(streams, stream)
			}
		} else {
			format := encodingcom.Format{
				OutputPreset: presetID,
				Destination:  e.getDestinations(job.ID, output.FileName),
			}
			formats = append(formats, format)
		}
	}
	if len(streams) > 0 {
		falseValue := encodingcom.YesNoBoolean(false)
		format := encodingcom.Format{
			Output:          []string{hlsOutput},
			Destination:     e.getDestinations(job.ID, job.StreamingParams.PlaylistFileName),
			SegmentDuration: job.StreamingParams.SegmentDuration,
			Stream:          streams,
			PackFiles:       &falseValue,
		}
		formats = append(formats, format)
	}
	return formats, nil
}

func (e *encodingComProvider) JobStatus(job *db.Job) (*provider.JobStatus, error) {
	resp, err := e.client.GetStatus([]string{job.ProviderJobID}, false)
	if err != nil {
		return nil, err
	}
	if len(resp) < 1 {
		return nil, errors.New("invalid value returned by the Encoding.com API: []")
	}
	var (
		sourceInfo provider.SourceInfo
		mediaInfo  *encodingcom.MediaInfo
	)
	status := e.statusMap(resp[0].MediaStatus)
	if status == provider.StatusFinished {
		sourceInfo, mediaInfo, err = e.sourceInfo(job.ProviderJobID)
		if err != nil {
			return nil, err
		}
	}
	return &provider.JobStatus{
		ProviderJobID: job.ProviderJobID,
		ProviderName:  "encoding.com",
		Status:        status,
		Progress:      resp[0].Progress,
		ProviderStatus: map[string]interface{}{
			"sourcefile":   resp[0].SourceFile,
			"timeleft":     resp[0].TimeLeft,
			"created":      resp[0].CreateDate,
			"started":      resp[0].StartDate,
			"finished":     resp[0].FinishDate,
			"formatStatus": e.getFormatStatus(resp),
		},
		Output: provider.JobOutput{
			Destination: e.getOutputDestination(job),
			Files:       e.getOutputDestinationStatus(resp, mediaInfo),
		},
		SourceInfo: sourceInfo,
	}, nil
}

func (e *encodingComProvider) sourceInfo(id string) (provider.SourceInfo, *encodingcom.MediaInfo, error) {
	var sourceInfo provider.SourceInfo
	info, err := e.client.GetMediaInfo(id)
	if err != nil {
		return sourceInfo, nil, err
	}
	sourceInfo.Width, sourceInfo.Height, err = e.parseSize(info.Size)
	sourceInfo.Duration = info.Duration
	sourceInfo.VideoCodec = info.VideoCodec
	return sourceInfo, info, err
}

func (e *encodingComProvider) parseSize(size string) (width int64, height int64, err error) {
	parts := strings.SplitN(size, "x", 2)
	if len(parts) < 2 {
		return width, height, fmt.Errorf("invalid size returned by the Encoding.com API: %q", size)
	}
	width, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return width, height, fmt.Errorf("invalid size returned by the Encoding.com API (%q): %s", size, err)
	}
	height, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return width, height, fmt.Errorf("invalid size returned by the Encoding.com API (%q): %s", size, err)
	}
	return width, height, nil
}

func (e *encodingComProvider) adjustSize(reportedSize string, sourceInfo *encodingcom.MediaInfo) (width int64, height int64, err error) {
	width, height, err = e.parseSize(reportedSize)
	if err != nil || sourceInfo == nil {
		return width, height, err
	}
	if width != 0 && height != 0 {
		return width, height, nil
	}
	sourceWidth, sourceHeight, err := e.parseSize(sourceInfo.Size)
	if err != nil {
		return 0, 0, err
	}
	if (sourceInfo.Rotation/90)%2 == 1 {
		sourceWidth, sourceHeight = sourceHeight, sourceWidth
	}
	if width == 0 {
		ratio := float64(sourceWidth) / float64(sourceHeight)
		width = int64(math.Floor(float64(height)*ratio + .5))
	}
	if height == 0 {
		ratio := float64(sourceHeight) / float64(sourceWidth)
		height = int64(math.Floor(float64(width)*ratio + .5))
	}
	return width, height, nil
}

func (e *encodingComProvider) getFormatStatus(status []encodingcom.StatusResponse) []string {
	formatStatusList := []string{}
	formats := status[0].Formats
	for _, formatStatus := range formats {
		formatStatusList = append(formatStatusList, formatStatus.Status)
	}
	return formatStatusList
}

func (e *encodingComProvider) getOutputDestinationStatus(status []encodingcom.StatusResponse, sourceInfo *encodingcom.MediaInfo) []provider.OutputFile {
	var outputFiles []provider.OutputFile
	formats := status[0].Formats
	for _, formatStatus := range formats {
		for idx, ds := range formatStatus.Destinations {
			destinationName := ds.Name
			if formatStatus.Output == hlsOutput {
				streams := formatStatus.Stream
				if idx < len(streams) {
					destinationNameParts := strings.Split(destinationName, "/")
					partsLength := len(destinationNameParts)
					fixedDestination := append(
						destinationNameParts[0:partsLength-2],
						streams[idx].SubPath,
						destinationNameParts[partsLength-1],
					)
					destinationName = strings.Join(fixedDestination, "/")
				}
			}
			container := formatStatus.Output
			if container == hlsOutput {
				container = "m3u8"
			}
			fileSize, _ := strconv.ParseInt(formatStatus.FileSize, 10, 64)
			file := provider.OutputFile{
				Path:       e.destinationMedia(destinationName),
				Container:  container,
				VideoCodec: formatStatus.VideoCodec,
				FileSize:   fileSize,
			}
			file.Width, file.Height, _ = e.adjustSize(formatStatus.Size, sourceInfo)
			outputFiles = append(outputFiles, file)
		}
	}
	return outputFiles
}

func (e *encodingComProvider) getOutputDestination(job *db.Job) string {
	parts := httpS3Regexp.FindStringSubmatch(strings.Trim(e.config.EncodingCom.Destination, "/"))
	if len(parts) > 0 {
		return fmt.Sprintf("s3://%s/%s/%s/", parts[1], parts[2], job.ID)
	}
	return strings.TrimRight(e.config.EncodingCom.Destination, "/") + "/" + job.ID
}

func (e *encodingComProvider) destinationMedia(input string) string {
	parts := httpS3Regexp.FindStringSubmatch(input)
	if len(parts) > 0 {
		return fmt.Sprintf("s3://%s/%s", parts[1], parts[2])
	}
	return input
}

func (e *encodingComProvider) statusMap(encodingComStatus string) provider.Status {
	switch strings.ToLower(encodingComStatus) {
	case "new", "waiting for encoder":
		return provider.StatusQueued
	case "downloading", "ready to process", "processing", "saving":
		return provider.StatusStarted
	case "finished":
		return provider.StatusFinished
	case "error":
		return provider.StatusFailed
	default:
		return provider.StatusUnknown
	}
}

func (e *encodingComProvider) CancelJob(id string) error {
	_, err := e.client.CancelMedia(id)
	return err
}

func (e *encodingComProvider) Healthcheck() error {
	status, err := encodingcom.APIStatus(e.config.EncodingCom.StatusEndpoint)
	if err != nil {
		return err
	}
	if !status.OK() {
		return fmt.Errorf("Status code: %s.\nIncident: %s\nStatus: %s", status.StatusCode, status.Incident, status.Status)
	}
	return nil
}

func (e *encodingComProvider) Capabilities() provider.Capabilities {
	return provider.Capabilities{
		InputFormats:  []string{"prores", "h264"},
		OutputFormats: []string{"mp4", "hls", "webm"},
		Destinations:  []string{"akamai", "s3"},
	}
}

func encodingComFactory(cfg *config.Config) (provider.TranscodingProvider, error) {
	if cfg.EncodingCom.UserID == "" || cfg.EncodingCom.UserKey == "" {
		return nil, errEncodingComInvalidConfig
	}
	client, err := encodingcom.NewClient("https://manage.encoding.com", cfg.EncodingCom.UserID, cfg.EncodingCom.UserKey)
	if err != nil {
		return nil, err
	}
	return &encodingComProvider{client: client, config: cfg}, nil
}
