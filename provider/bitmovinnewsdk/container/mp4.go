package container

import (
	"fmt"
	"path"

	"github.com/NYTimes/video-transcoding-api/provider"
	"github.com/bitmovin/bitmovin-api-sdk-go"
	"github.com/bitmovin/bitmovin-api-sdk-go/model"
	"github.com/bitmovin/bitmovin-api-sdk-go/query"
	"github.com/pkg/errors"
)

// MP4Assembler is an assembler that creates MP4 outputs based on a cfg
type MP4Assembler struct {
	api *bitmovin.BitmovinApi
}

// NewMP4Assembler creates and returns an MP4Assembler
func NewMP4Assembler(api *bitmovin.BitmovinApi) *MP4Assembler {
	return &MP4Assembler{api: api}
}

// Assemble creates MP4 outputs
func (a *MP4Assembler) Assemble(cfg AssemblerCfg) error {
	_, err := a.api.Encoding.Encodings.Muxings.Mp4.Create(cfg.EncID, model.Mp4Muxing{
		Filename:             path.Base(cfg.OutputFilename),
		Streams:              []model.MuxingStream{cfg.VidMuxingStream, cfg.AudMuxingStream},
		StreamConditionsMode: model.StreamConditionsMode_DROP_STREAM,
		Outputs: []model.EncodingOutput{{
			OutputId:   cfg.OutputID,
			OutputPath: path.Dir(path.Join(cfg.DestPath, cfg.OutputFilename)),
		}},
	})
	if err != nil {
		return errors.Wrap(err, "creating mp4 muxing")
	}

	return nil
}

// MP4StatusEnricher is responsible for adding MP4 output info to a job status
type MP4StatusEnricher struct {
	api *bitmovin.BitmovinApi
}

// NewMP4StatusEnricher creates and returns an MP4StatusEnricher
func NewMP4StatusEnricher(api *bitmovin.BitmovinApi) *MP4StatusEnricher {
	return &MP4StatusEnricher{api: api}
}

// Enrich populates information about MP4 outputs if they exist
func (e *MP4StatusEnricher) Enrich(s provider.JobStatus) (provider.JobStatus, error) {
	var totalCount int64 = 1
	var muxings []model.Mp4Muxing
	for int64(len(muxings)) < totalCount {
		resp, err := e.api.Encoding.Encodings.Muxings.Mp4.List(s.ProviderJobID, func(params *query.Mp4MuxingListQueryParams) {
			params.Offset = int32(len(muxings))
			params.Limit = 100
		})
		if err != nil {
			return s, errors.Wrap(err, "retrieving MP4 muxings from the Bitmovin API")
		}

		totalCount = int64Value(resp.TotalCount)
		muxings = append(muxings, resp.Items...)
	}

	for _, muxing := range muxings {
		info, err := e.api.Encoding.Encodings.Muxings.Mp4.Information.Get(s.ProviderJobID, muxing.Id)
		if err != nil {
			return s, errors.Wrapf(err, "retrieving muxing information with ID %q", muxing.Id)
		}

		if len(info.VideoTracks) == 0 {
			return s, fmt.Errorf("no video track found for encodingID %s muxingID %s", s.ProviderJobID, muxing.Id)
		}

		firstVidTrack := info.VideoTracks[0]
		s.Output.Files = append(s.Output.Files, provider.OutputFile{
			Path:       s.Output.Destination + muxing.Filename,
			Container:  info.ContainerFormat,
			FileSize:   int64Value(info.FileSize),
			VideoCodec: firstVidTrack.Codec,
			Width:      int64(int32Value(firstVidTrack.FrameWidth)),
			Height:     int64(int32Value(firstVidTrack.FrameHeight)),
		})
	}

	return s, nil
}
