package service

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"

	"github.com/NYTimes/gizmo/web"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/provider"
	"github.com/NYTimes/video-transcoding-api/swagger"
)

// swagger:route POST /jobs jobs newJob
//
// Creates a new transcoding job.
//
//     Responses:
//       200: job
//       400: invalidJob
//       500: genericError
func (s *TranscodingService) newTranscodeJob(r *http.Request) swagger.GizmoJSONResponse {
	defer r.Body.Close()
	var input newTranscodeJobInput
	providerFactory, err := input.ProviderFactory(r.Body)
	if err != nil {
		return newInvalidJobResponse(err)
	}
	providerObj, err := providerFactory(s.config)
	if err != nil {
		formattedErr := fmt.Errorf("Error initializing provider %s for new job: %v %s", input.Payload.Provider, providerObj, err)
		if _, ok := err.(provider.InvalidConfigError); ok {
			return newInvalidJobResponse(formattedErr)
		}
		return swagger.NewErrorResponse(formattedErr)
	}
	job := db.Job{
		SourceMedia:     input.Payload.Source,
		StreamingParams: input.Payload.StreamingParams,
	}
	outputs := make([]db.TranscodeOutput, len(input.Payload.Outputs))
	for i, output := range input.Payload.Outputs {
		presetMap, presetErr := s.db.GetPresetMap(output.Preset)
		if presetErr != nil {
			if presetErr == db.ErrPresetMapNotFound {
				return newInvalidJobResponse(presetErr)
			}
			return swagger.NewErrorResponse(presetErr)
		}
		fileName := output.FileName
		if fileName == "" {
			fileName = s.defaultFileName(input.Payload.Source, presetMap)
		}
		outputs[i] = db.TranscodeOutput{FileName: fileName, Preset: *presetMap}
	}
	job.Outputs = outputs
	job.ID, err = s.genID()
	if err != nil {
		return swagger.NewErrorResponse(err)
	}
	if job.StreamingParams.Protocol == "hls" {
		if job.StreamingParams.PlaylistFileName == "" {
			job.StreamingParams.PlaylistFileName = "hls/index.m3u8"
		}
		if job.StreamingParams.SegmentDuration == 0 {
			job.StreamingParams.SegmentDuration = s.config.DefaultSegmentDuration
		}
	}
	jobStatus, err := providerObj.Transcode(&job)
	if err == provider.ErrPresetMapNotFound {
		return newInvalidJobResponse(err)
	}
	if err != nil {
		providerError := fmt.Errorf("Error with provider %q: %s", input.Payload.Provider, err)
		return swagger.NewErrorResponse(providerError)
	}
	jobStatus.ProviderName = input.Payload.Provider
	job.ProviderName = jobStatus.ProviderName
	job.ProviderJobID = jobStatus.ProviderJobID
	err = s.db.CreateJob(&job)
	if err != nil {
		return swagger.NewErrorResponse(err)
	}
	return newJobResponse(job.ID)
}

func (s *TranscodingService) genID() (string, error) {
	var data [8]byte
	n, err := rand.Read(data[:])
	if err != nil {
		return "", err
	}
	if n != len(data) {
		return "", io.ErrShortWrite
	}
	return fmt.Sprintf("%x", data), nil
}

func (s *TranscodingService) defaultFileName(source string, preset *db.PresetMap) string {
	sourceExtension := filepath.Ext(source)
	_, source = path.Split(source)
	source = source[:len(source)-len(sourceExtension)]
	pattern := "%s_%s.%s"
	if preset.OutputOpts.Extension == "m3u8" {
		pattern = "hls/" + pattern
	}
	return fmt.Sprintf(pattern, source, preset.Name, preset.OutputOpts.Extension)
}

// swagger:route GET /jobs/{jobId} jobs getJob
//
// Finds a trancode job using its ID.
// It also queries the provider to get the status of the job.
//
//     Responses:
//       200: jobStatus
//       404: jobNotFound
//       410: jobNotFoundInTheProvider
//       500: genericError
func (s *TranscodingService) getTranscodeJob(r *http.Request) swagger.GizmoJSONResponse {
	var params getTranscodeJobInput
	params.loadParams(web.Vars(r))
	return s.getJobStatusResponse(s.getTranscodeJobByID(params.JobID))
}

func (s *TranscodingService) getJobStatusResponse(job *db.Job, status *provider.JobStatus, p provider.TranscodingProvider, err error) swagger.GizmoJSONResponse {
	if err != nil {
		if err == db.ErrJobNotFound {
			return newJobNotFoundResponse(err)
		}
		if p != nil {
			providerError := fmt.Errorf("Error with provider %q when trying to retrieve job id %q: %s", job.ProviderName, job.ID, err)
			if _, ok := err.(provider.JobNotFoundError); ok {
				return newJobNotFoundProviderResponse(providerError)
			}
			return swagger.NewErrorResponse(providerError)
		}
		return swagger.NewErrorResponse(err)
	}
	return newJobStatusResponse(status)
}

func (s *TranscodingService) getTranscodeJobByID(jobID string) (*db.Job, *provider.JobStatus, provider.TranscodingProvider, error) {
	job, err := s.db.GetJob(jobID)
	if err != nil {
		if err == db.ErrJobNotFound {
			return nil, nil, nil, err
		}
		return nil, nil, nil, fmt.Errorf("error retrieving job with id %q: %s", jobID, err)
	}
	providerFactory, err := provider.GetProviderFactory(job.ProviderName)
	if err != nil {
		return job, nil, nil, fmt.Errorf("unknown provider %q for job id %q", job.ProviderName, jobID)
	}
	providerObj, err := providerFactory(s.config)
	if err != nil {
		return job, nil, nil, fmt.Errorf("error initializing provider %q on job id %q: %s %s", job.ProviderName, jobID, providerObj, err)
	}
	jobStatus, err := providerObj.JobStatus(job)
	if err != nil {
		return job, nil, providerObj, err
	}
	jobStatus.ProviderName = job.ProviderName
	return job, jobStatus, providerObj, nil
}

// swagger:route POST /jobs/{jobId}/cancel jobs cancelJob
//
// Creates a new transcoding job.
//
//     Responses:
//       200: jobStatus
//       404: jobNotFound
//       410: jobNotFoundInTheProvider
//       500: genericError
func (s *TranscodingService) cancelTranscodeJob(r *http.Request) swagger.GizmoJSONResponse {
	var params cancelTranscodeJobInput
	params.loadParams(web.Vars(r))
	job, _, prov, err := s.getTranscodeJobByID(params.JobID)
	if err != nil {
		if err == db.ErrJobNotFound {
			return newJobNotFoundResponse(err)
		}
		if _, ok := err.(provider.JobNotFoundError); ok {
			return newJobNotFoundProviderResponse(err)
		}
		return swagger.NewErrorResponse(err)
	}
	err = prov.CancelJob(job.ProviderJobID)
	if err != nil {
		return swagger.NewErrorResponse(err)
	}
	status, err := prov.JobStatus(job)
	if err != nil {
		return swagger.NewErrorResponse(err)
	}
	status.ProviderName = job.ProviderName
	return newJobStatusResponse(status)
}
