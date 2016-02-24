package service

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/provider"
)

// swagger:route POST /jobs jobs newJob
//
// Creates a new transcoding job.
//
//     Responses:
//       200: job
//       400: invalidJob
//       500: genericError
func (s *TranscodingService) newTranscodeJob(r *http.Request) gizmoResponse {
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
		return newErrorResponse(formattedErr)
	}

	var jobStatus *provider.JobStatus
	if len(input.Payload.Profiles) > 0 {
		profileProvider, ok := providerObj.(provider.ProfileTranscodingProvider)
		if !ok {
			return newInvalidJobResponse(fmt.Errorf("Provider %q does not support profile-based transcoding", input.Payload.Provider))
		}
		jobStatus, err = profileProvider.TranscodeWithProfiles(input.Payload.Source, input.Payload.Profiles)
	} else {
		presetProvider, ok := providerObj.(provider.PresetTranscodingProvider)
		if !ok {
			return newInvalidJobResponse(fmt.Errorf("Provider %q does not support preset-based transcoding", input.Payload.Provider))
		}
		presets := make([]db.Preset, len(input.Payload.Presets))
		for i, presetID := range input.Payload.Presets {
			preset, err := s.db.GetPreset(presetID)
			if err != nil {
				if err == db.ErrPresetNotFound {
					return newInvalidJobResponse(err)
				}
				return newErrorResponse(err)
			}
			presets[i] = *preset
		}
		jobStatus, err = presetProvider.TranscodeWithPresets(input.Payload.Source, presets)
	}

	if err != nil {
		providerError := fmt.Errorf("Error with provider %q: %s", input.Payload.Provider, err)
		return newErrorResponse(providerError)
	}
	jobStatus.ProviderName = input.Payload.Provider

	job := db.Job{ProviderName: jobStatus.ProviderName, ProviderJobID: jobStatus.ProviderJobID}
	err = s.db.CreateJob(&job)
	if err != nil {
		return newErrorResponse(err)
	}
	return newJobResponse(job.ID)
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
func (s *TranscodingService) getTranscodeJob(r *http.Request) gizmoResponse {
	var params getTranscodeJobInput
	params.loadParams(mux.Vars(r))
	jobID := params.JobID
	job, err := s.db.GetJob(jobID)
	if err != nil {
		if err == db.ErrJobNotFound {
			return newJobNotFoundResponse(err)
		}
		return newErrorResponse(fmt.Errorf("error retrieving job with id %q: %s", jobID, err))
	}
	providerFactory, err := provider.GetProviderFactory(job.ProviderName)
	if err != nil {
		return newErrorResponse(fmt.Errorf("unknown provider %q for job id %q", job.ProviderName, jobID))
	}
	providerObj, err := providerFactory(s.config)
	if err != nil {
		return newErrorResponse(fmt.Errorf("error initializing provider %q on job id %q: %s %s", job.ProviderName, jobID, providerObj, err))
	}
	jobStatus, err := providerObj.JobStatus(job.ProviderJobID)
	if err != nil {
		providerError := fmt.Errorf("Error with provider %q when trying to retrieve job id %q: %s", job.ProviderName, jobID, err)
		if _, ok := err.(provider.JobNotFoundError); ok {
			return newJobNotFoundProviderResponse(providerError)
		}
		return newErrorResponse(providerError)
	}
	jobStatus.ProviderName = job.ProviderName
	return newJobStatusResponse(jobStatus)
}
