package service

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/NYTimes/gizmo/web"
	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/provider"
	"github.com/nytm/video-transcoding-api/swagger"
	"golang.org/x/net/context"
)

const maxJobTimeout = 8 * time.Hour

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
	presetsMap := make([]db.PresetMap, len(input.Payload.Presets))
	for i, presetID := range input.Payload.Presets {
		presetMap, err := s.db.GetPresetMap(presetID)
		if err != nil {
			if err == db.ErrPresetMapNotFound {
				return newInvalidJobResponse(err)
			}
			return swagger.NewErrorResponse(err)
		}
		presetsMap[i] = *presetMap
	}
	jobID, err := s.genID()
	if err != nil {
		return swagger.NewErrorResponse(err)
	}
	transcodeProfile := provider.TranscodeProfile{
		SourceMedia:     input.Payload.Source,
		Presets:         presetsMap,
		StreamingParams: input.Payload.StreamingParams,
	}
	job := db.Job{
		ID:                     jobID,
		StatusCallbackURL:      input.Payload.StatusCallbackURL,
		StatusCallbackInterval: input.Payload.StatusCallbackInterval,
		CompletionCallbackURL:  input.Payload.CompletionCallbackURL,
	}
	jobStatus, err := providerObj.Transcode(&job, transcodeProfile)
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
	if transcodeProfile.StreamingParams.Protocol != "" {
		job.StreamingParams = db.StreamingParams{
			SegmentDuration: transcodeProfile.StreamingParams.SegmentDuration,
			Protocol:        transcodeProfile.StreamingParams.Protocol,
		}
	}
	err = s.db.CreateJob(&job)
	if err != nil {
		return swagger.NewErrorResponse(err)
	}
	if job.StatusCallbackURL != "" || job.CompletionCallbackURL != "" {
		ctx, cancel := context.WithTimeout(context.Background(), maxJobTimeout)
		defer cancel()
		go s.statusCallback(ctx, job)
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
	return s.getJobStatusResponse(s.getTranscodeJobByID(params.JobID)).(swagger.GizmoJSONResponse)
}

func (s *TranscodingService) getJobStatusResponse(job *db.Job, jobStatus *provider.JobStatus, providerObj provider.TranscodingProvider, err error) interface{} {
	if err != nil {
		if err == db.ErrJobNotFound {
			return newJobNotFoundResponse(err)
		}
		if providerObj != nil {
			providerError := fmt.Errorf("Error with provider %q when trying to retrieve job id %q: %s", job.ProviderName, job.ID, err)
			if _, ok := err.(provider.JobNotFoundError); ok {
				return newJobNotFoundProviderResponse(providerError)
			}
			return swagger.NewErrorResponse(providerError)
		}
		return swagger.NewErrorResponse(err)
	}
	return newJobStatusResponse(jobStatus)
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
	jobStatus, err := providerObj.JobStatus(job.ProviderJobID)
	if err != nil {
		return job, nil, providerObj, err
	}
	jobStatus.ProviderName = job.ProviderName
	return job, jobStatus, providerObj, nil
}

func (s *TranscodingService) statusCallback(ctx context.Context, job db.Job) error {
	deadline, _ := ctx.Deadline()
	for now := time.Now(); now.Before(deadline); now = time.Now() {
		job, jobStatus, providerObj, err := s.getTranscodeJobByID(job.ID)
		jobStatusResponseObj := s.getJobStatusResponse(job, jobStatus, providerObj, err)
		var callbackPayload interface{}
		if _, ok := jobStatusResponseObj.(*jobStatusResponse); ok {
			callbackPayload = jobStatus
		} else {
			_, _, errorObj := jobStatusResponseObj.(swagger.GizmoJSONResponse).Result()
			callbackPayload = errorObj
		}
		if jobStatus.Status != provider.StatusQueued &&
			jobStatus.Status != provider.StatusStarted {
			if job.CompletionCallbackURL != "" {
				err := s.postStatusToCallback(callbackPayload, job.CompletionCallbackURL)
				if err != nil {
					continue
				}
			}
			break
		}
		if job.StatusCallbackURL != "" {
			err := s.postStatusToCallback(callbackPayload, job.StatusCallbackURL)
			if err != nil {
				continue
			}
		}
		time.Sleep(time.Duration(job.StatusCallbackInterval) * time.Second)
	}
	return nil
}

func (s *TranscodingService) postStatusToCallback(payloadStruct interface{}, callbackURL string) error {
	jsonPayload, err := json.Marshal(payloadStruct)
	if err != nil {
		fmt.Printf("Error generating response for status callback: %v\n", err)
		return err
	}
	req, err := http.NewRequest("POST", callbackURL, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error calling status callback URL %s : %v\n", callbackURL, err)
		return err
	}
	resp.Body.Close()
	return nil
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
	status, err := prov.JobStatus(job.ProviderJobID)
	if err != nil {
		return swagger.NewErrorResponse(err)
	}
	status.ProviderName = job.ProviderName
	return newJobStatusResponse(status)
}
