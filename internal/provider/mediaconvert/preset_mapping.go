package mediaconvert

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/mediaconvert/types"
	"github.com/pkg/errors"
	"github.com/video-dev/video-transcoding-api/v2/db"
	"github.com/video-dev/video-transcoding-api/v2/internal/provider"
)

func providerStatusFrom(status types.JobStatus) provider.Status {
	switch status {
	case types.JobStatusSubmitted:
		return provider.StatusQueued
	case types.JobStatusProgressing:
		return provider.StatusStarted
	case types.JobStatusComplete:
		return provider.StatusFinished
	case types.JobStatusCanceled:
		return provider.StatusCanceled
	case types.JobStatusError:
		return provider.StatusFailed
	default:
		return provider.StatusUnknown
	}
}

func containerFrom(container string) (types.ContainerType, error) {
	container = strings.ToLower(container)
	switch container {
	case "m3u8":
		return types.ContainerTypeM3u8, nil
	case "mp4":
		return types.ContainerTypeMp4, nil
	case "none":
		return types.ContainerTypeRaw, nil
	default:
		return "", fmt.Errorf("container %q not supported with mediaconvert", container)
	}
}

func h264RateControlModeFrom(rateControl string) (types.H264RateControlMode, error) {
	rateControl = strings.ToLower(rateControl)
	switch rateControl {
	case "vbr":
		return types.H264RateControlModeVbr, nil
	case "", "cbr":
		return types.H264RateControlModeCbr, nil
	case "qvbr":
		return types.H264RateControlModeQvbr, nil
	default:
		return "", fmt.Errorf("rate control mode %q is not supported with mediaconvert", rateControl)
	}
}

func h264CodecProfileFrom(profile string) (types.H264CodecProfile, error) {
	profile = strings.ToLower(profile)
	switch profile {
	case "baseline":
		return types.H264CodecProfileBaseline, nil
	case "main":
		return types.H264CodecProfileMain, nil
	case "", "high":
		return types.H264CodecProfileHigh, nil
	default:
		return "", fmt.Errorf("h264 profile %q is not supported with mediaconvert", profile)
	}
}

func h264CodecLevelFrom(level string) (types.H264CodecLevel, error) {
	switch level {
	case "":
		return types.H264CodecLevelAuto, nil
	case "1", "1.0":
		return types.H264CodecLevelLevel1, nil
	case "1.1":
		return types.H264CodecLevelLevel11, nil
	case "1.2":
		return types.H264CodecLevelLevel12, nil
	case "1.3":
		return types.H264CodecLevelLevel13, nil
	case "2", "2.0":
		return types.H264CodecLevelLevel2, nil
	case "2.1":
		return types.H264CodecLevelLevel21, nil
	case "2.2":
		return types.H264CodecLevelLevel22, nil
	case "3", "3.0":
		return types.H264CodecLevelLevel3, nil
	case "3.1":
		return types.H264CodecLevelLevel31, nil
	case "3.2":
		return types.H264CodecLevelLevel32, nil
	case "4", "4.0":
		return types.H264CodecLevelLevel4, nil
	case "4.1":
		return types.H264CodecLevelLevel41, nil
	case "4.2":
		return types.H264CodecLevelLevel42, nil
	case "5", "5.0":
		return types.H264CodecLevelLevel5, nil
	case "5.1":
		return types.H264CodecLevelLevel51, nil
	case "5.2":
		return types.H264CodecLevelLevel52, nil
	default:
		return "", fmt.Errorf("h264 level %q is not supported with mediaconvert", level)
	}
}

func h264InterlaceModeFrom(mode string) (types.H264InterlaceMode, error) {
	mode = strings.ToLower(mode)
	switch mode {
	case "", "progressive":
		return types.H264InterlaceModeProgressive, nil
	default:
		return "", fmt.Errorf("h264 interlace mode %q is not supported with mediaconvert", mode)
	}
}

func thumbnailPresetFrom(preset db.Preset) (*types.VideoDescription, error) {

	thumbnailPreset := types.VideoDescription{}
	
	if preset.Thumbnail.Width != "" {
		width, err := strconv.ParseInt(preset.Thumbnail.Width, 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing thumbnail width %q to int32", preset.Thumbnail.Width)
		}
		thumbnailPreset.Width = int32(width)
	}

	if preset.Thumbnail.Height != "" {
		height, err := strconv.ParseInt(preset.Thumbnail.Height, 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing thumbnail width %q to int32", preset.Thumbnail.Height)
		}
		thumbnailPreset.Height = int32(height)
	}

	codec := strings.ToLower(preset.Thumbnail.Codec)
	switch(codec){
		case "frame_capture":
			numerator, err := strconv.ParseInt(preset.Thumbnail.FrameCaptureNumerator, 10, 32)
			if err != nil {
				return nil, errors.Wrapf(err, "parsing video numerator %q to int32", preset.Thumbnail.FrameCaptureNumerator)
			}

			denominator, err := strconv.ParseInt(preset.Thumbnail.FrameCaptureDenominator, 10, 32)
			if err != nil {
				return nil, errors.Wrapf(err, "parsing video denominator %q to int32", preset.Thumbnail.FrameCaptureDenominator)
			}

			maxCaptures, err := strconv.ParseInt(preset.Thumbnail.MaxCaptures, 10, 32)
			if err != nil {
				return nil, errors.Wrapf(err, "parsing video maxCaptures %q to int32", preset.Thumbnail.MaxCaptures)
			}

			quality, err := strconv.ParseInt(preset.Thumbnail.Quality, 10, 32)
			if err != nil {
				return nil, errors.Wrapf(err, "parsing video quality %q to int32", preset.Thumbnail.Quality)
			}

			thumbnailPreset.CodecSettings = &types.VideoCodecSettings{
				Codec: types.VideoCodecFrameCapture,
				FrameCaptureSettings: &types.FrameCaptureSettings{
					FramerateNumerator: int32(numerator),
					FramerateDenominator: int32(denominator),
					MaxCaptures: int32(maxCaptures),
					Quality: int32(quality),
				},
			}
		default:
			return nil, fmt.Errorf("thumbnail codec %q is not yet supported with mediaconvert", codec)
	}
	
	return &thumbnailPreset, nil
}

func videoPresetFrom(preset db.Preset) (*types.VideoDescription, error) {
	videoPreset := types.VideoDescription{
		ScalingBehavior:   types.ScalingBehaviorDefault,
		TimecodeInsertion: types.VideoTimecodeInsertionDisabled,
		AntiAlias:         types.AntiAliasEnabled,
		RespondToAfd:      types.RespondToAfdNone,
	}

	if preset.Video.Width != "" {
		width, err := strconv.ParseInt(preset.Video.Width, 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing video width %q to int32", preset.Video.Width)
		}
		videoPreset.Width = int32(width)
	}

	if preset.Video.Height != "" {
		height, err := strconv.ParseInt(preset.Video.Height, 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing video height %q to int32", preset.Video.Height)
		}
		videoPreset.Height = int32(height)
	}

	codec := strings.ToLower(preset.Video.Codec)
	switch codec {
	case "h264":
		bitrate, err := strconv.ParseInt(preset.Video.Bitrate, 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing video bitrate %q to int32", preset.Video.Bitrate)
		}

		gopSize, err := strconv.ParseFloat(preset.Video.GopSize, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing gop size %q to float64", preset.Video.GopSize)
		}

		rateControl, err := h264RateControlModeFrom(preset.RateControl)
		if err != nil {
			return nil, err
		}

		profile, err := h264CodecProfileFrom(preset.Video.Profile)
		if err != nil {
			return nil, err
		}

		level, err := h264CodecLevelFrom(preset.Video.ProfileLevel)
		if err != nil {
			return nil, err
		}

		interlaceMode, err := h264InterlaceModeFrom(preset.Video.InterlaceMode)
		if err != nil {
			return nil, err
		}

		tuning := types.H264QualityTuningLevelSinglePassHq
		if preset.TwoPass {
			tuning = types.H264QualityTuningLevelMultiPassHq
		}

		var bframes int64
		if preset.Video.BFrames != "" {
			bframes, err = strconv.ParseInt(preset.Video.BFrames, 10, 32)
			if err != nil {
				return nil, errors.Wrapf(err, "parsing bframes %q to int32", preset.Video.BFrames)
			}
		}

		videoPreset.CodecSettings = &types.VideoCodecSettings{
			Codec: types.VideoCodecH264,
			H264Settings: &types.H264Settings{
				Bitrate:                             int32(bitrate),
				GopSize:                             gopSize,
				RateControlMode:                     rateControl,
				CodecProfile:                        profile,
				CodecLevel:                          level,
				InterlaceMode:                       interlaceMode,
				QualityTuningLevel:                  tuning,
				NumberBFramesBetweenReferenceFrames: int32(bframes),
			},
		}
	default:
		return nil, fmt.Errorf("video codec %q is not yet supported with mediaconvert", codec)
	}

	return &videoPreset, nil
}

func audioPresetFrom(preset db.Preset) (*types.AudioDescription, error) {
	audioPreset := types.AudioDescription{}

	codec := strings.ToLower(preset.Audio.Codec)
	switch codec {
	case "aac":
		bitrate, err := strconv.ParseInt(preset.Audio.Bitrate, 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing audio bitrate %q to int32", preset.Audio.Bitrate)
		}

		audioPreset.CodecSettings = &types.AudioCodecSettings{
			Codec: types.AudioCodecAac,
			AacSettings: &types.AacSettings{
				SampleRate:      defaultAudioSampleRate,
				Bitrate:         int32(bitrate),
				CodecProfile:    types.AacCodecProfileLc,
				CodingMode:      types.AacCodingModeCodingMode20,
				RateControlMode: types.AacRateControlModeCbr,
			},
		}
	default:
		return nil, fmt.Errorf("audio codec %q is not yet supported with mediaconvert", codec)
	}

	return &audioPreset, nil
}
