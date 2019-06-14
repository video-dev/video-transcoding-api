package container

import (
	"path"
	"path/filepath"

	"github.com/NYTimes/video-transcoding-api/provider/bitmovinnewsdk/types"

	"github.com/NYTimes/video-transcoding-api/provider"

	"github.com/bitmovin/bitmovin-api-sdk-go"
	"github.com/bitmovin/bitmovin-api-sdk-go/model"
	"github.com/pkg/errors"
)

const (
	// CustomDataKeyManifest is used as the base key to store the manifestID in an encoding
	CustomDataKeyManifest = "manifest"
	// CustomDataKeyManifestID is the key used to store the manifestID in an encoding
	CustomDataKeyManifestID = "id"
)

// HLSAssembler is an assembler that creates HLS outputs based on a cfg
type HLSAssembler struct {
	api *bitmovin.BitmovinApi
}

// NewHLSAssembler creates and returns an HLSAssembler
func NewHLSAssembler(api *bitmovin.BitmovinApi) *HLSAssembler {
	return &HLSAssembler{api: api}
}

// Assemble creates HLS outputs
func (a *HLSAssembler) Assemble(cfg AssemblerCfg) error {
	if !cfg.SkipAudioCreation {
		audTSMuxing, err := a.api.Encoding.Encodings.Muxings.Ts.Create(cfg.EncID, model.TsMuxing{
			SegmentLength: floatToPtr(float64(cfg.SegDuration)),
			SegmentNaming: "seg_%number%.ts",
			Streams:       []model.MuxingStream{cfg.AudMuxingStream},
			Outputs: []model.EncodingOutput{{
				OutputId:   cfg.OutputID,
				OutputPath: path.Join(cfg.ManifestMasterPath, cfg.AudCfgID),
			}},
		})
		if err != nil {
			return errors.Wrap(err, "creating audio ts muxing")
		}

		_, err = a.api.Encoding.Manifests.Hls.Media.Audio.Create(cfg.ManifestID, model.AudioMediaInfo{
			Uri:             cfg.AudCfgID + ".m3u8",
			GroupId:         cfg.AudCfgID,
			Language:        "en",
			Name:            cfg.AudCfgID,
			IsDefault:       boolToPtr(false),
			Autoselect:      boolToPtr(false),
			Forced:          boolToPtr(false),
			SegmentPath:     cfg.AudCfgID,
			Characteristics: []string{"public.accessibility.describes-video"},
			EncodingId:      cfg.EncID,
			StreamId:        cfg.AudStreamID,
			MuxingId:        audTSMuxing.Id,
		})
		if err != nil {
			return errors.Wrap(err, "creating audio media")
		}
	}

	vidTSMuxing, err := a.api.Encoding.Encodings.Muxings.Ts.Create(cfg.EncID, model.TsMuxing{
		SegmentLength: floatToPtr(float64(cfg.SegDuration)),
		SegmentNaming: "seg_%number%.ts",
		Streams:       []model.MuxingStream{cfg.VidMuxingStream},
		Outputs: []model.EncodingOutput{{
			OutputId:   cfg.OutputID,
			OutputPath: path.Join(cfg.ManifestMasterPath, cfg.VidCfgID),
		}},
	})
	if err != nil {
		return errors.Wrap(err, "creating video ts muxing")
	}

	vidManifestLoc, err := filepath.Rel(cfg.ManifestMasterPath, path.Join(cfg.DestPath, cfg.OutputFilename))
	if err != nil {
		return errors.Wrap(err, "constructing video manifest location")
	}

	vidSegLoc, err := filepath.Rel(path.Dir(path.Join(cfg.DestPath, cfg.OutputFilename)), path.Join(cfg.ManifestMasterPath, cfg.VidCfgID))
	if err != nil {
		return errors.Wrap(err, "constructing video segment location")
	}

	_, err = a.api.Encoding.Manifests.Hls.Streams.Create(cfg.ManifestID, model.StreamInfo{
		Audio:       cfg.AudCfgID,
		Uri:         vidManifestLoc,
		SegmentPath: vidSegLoc,
		EncodingId:  cfg.EncID,
		StreamId:    cfg.VidStreamID,
		MuxingId:    vidTSMuxing.Id,
	})
	if err != nil {
		return errors.Wrap(err, "creating video stream info")
	}

	return nil
}

// HLSStatusEnricher is responsible for adding output HLS info to a job status
type HLSStatusEnricher struct {
	api *bitmovin.BitmovinApi
}

// NewHLSStatusEnricher creates and returns an HLSStatusEnricher
func NewHLSStatusEnricher(api *bitmovin.BitmovinApi) *HLSStatusEnricher {
	return &HLSStatusEnricher{api: api}
}

// Enrich populates information about the HLS output if it exists
func (e *HLSStatusEnricher) Enrich(s provider.JobStatus) (provider.JobStatus, error) {
	data, err := e.api.Encoding.Encodings.Customdata.Get(s.ProviderJobID)
	if err != nil {
		return s, errors.Wrap(err, "retrieving the encoding from the Bitmovin API")
	}

	manifestID, err := types.CustomDataStringValAtKeys(data.CustomData, CustomDataKeyManifest, CustomDataKeyManifestID)
	if err == nil && manifestID != "" {
		s.ProviderStatus["manifestStatus"] = model.Status_FINISHED
	}

	return s, nil
}

func floatToPtr(f float64) *float64 {
	return &f
}

func boolToPtr(b bool) *bool {
	return &b
}
