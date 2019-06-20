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

// MOVAssembler is an assembler that creates MOV outputs based on a cfg
type MOVAssembler struct {
	api *bitmovin.BitmovinApi
}

// NewMOVAssembler creates and returns an MOVAssembler
func NewMOVAssembler(api *bitmovin.BitmovinApi) *MOVAssembler {
	return &MOVAssembler{api: api}
}

// Assemble creates MOV outputs
func (a *MOVAssembler) Assemble(cfg AssemblerCfg) error {
	_, err := a.api.Encoding.Encodings.Muxings.ProgressiveMov.Create(cfg.EncID, model.ProgressiveMovMuxing{
		Filename:             path.Base(cfg.OutputFilename),
		Streams:              []model.MuxingStream{cfg.VidMuxingStream, cfg.AudMuxingStream},
		StreamConditionsMode: model.StreamConditionsMode_DROP_STREAM,
		Outputs: []model.EncodingOutput{{
			OutputId:   cfg.OutputID,
			OutputPath: path.Dir(path.Join(cfg.DestPath, cfg.OutputFilename)),
		}},
	})
	if err != nil {
		return errors.Wrap(err, "creating mov muxing")
	}

	return nil
}

// MOVStatusEnricher is responsible for adding MOV output info to a job status
type MOVStatusEnricher struct {
	api *bitmovin.BitmovinApi
}

// NewMOVStatusEnricher creates and returns an MOVStatusEnricher
func NewMOVStatusEnricher(api *bitmovin.BitmovinApi) *MOVStatusEnricher {
	return &MOVStatusEnricher{api: api}
}

// Enrich populates information about MOV outputs if they exist
func (e *MOVStatusEnricher) Enrich(s provider.JobStatus) (provider.JobStatus, error) {
	var totalCount int64 = 1
	var muxings []model.ProgressiveMovMuxing
	for int64(len(muxings)) < totalCount {
		resp, err := e.api.Encoding.Encodings.Muxings.ProgressiveMov.List(s.ProviderJobID, func(params *query.ProgressiveMovMuxingListQueryParams) {
			params.Offset = int32(len(muxings))
			params.Limit = 100
		})
		if err != nil {
			return s, errors.Wrap(err, "retrieving progressive MOV muxings from the Bitmovin API")
		}

		totalCount = int64Value(resp.TotalCount)
		muxings = append(muxings, resp.Items...)
	}

	for _, muxing := range muxings {
		info, err := e.api.Encoding.Encodings.Muxings.ProgressiveMov.Information.Get(s.ProviderJobID, muxing.Id)
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
