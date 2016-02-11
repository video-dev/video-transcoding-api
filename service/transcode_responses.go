package service

import (
	"net/http"

	"github.com/nytm/video-transcoding-api/provider"
)

// swagger:model
type partialJob struct {
	// unique identifier of the job
	//
	// unique: true
	JobID string `json:"jobId"`
}

// represents the job after creation, includes the id of the job only, that can
// be used for querying the current status of the job.
//
// swagger:response job
type jobResponse struct {
	// in: body
	Payload *partialJob

	baseResponse
}

func newJobResponse(jobID string) *jobResponse {
	return &jobResponse{
		baseResponse: baseResponse{
			payload: &partialJob{JobID: jobID},
			status:  http.StatusOK,
		},
	}
}

// swagger:response jobStatus
type jobStatusResponse struct {
	// in: body
	Payload *provider.JobStatus

	baseResponse
}

func newJobStatusResponse(jobStatus *provider.JobStatus) *jobStatusResponse {
	return &jobStatusResponse{
		baseResponse: baseResponse{
			payload: jobStatus,
			status:  http.StatusOK,
		},
	}
}

// swagger:response invalidJob
type invalidJobResponse struct {
	// in: body
	Error *errorResponse
}

func newInvalidJobResponse(err error) *invalidJobResponse {
	return &invalidJobResponse{Error: newErrorResponse(err).withStatus(http.StatusBadRequest)}
}

func (r *invalidJobResponse) Result() (int, interface{}, error) {
	return r.Error.Result()
}

// swagger:response jobNotFound
type jobNotFoundResponse struct {
	// in: body
	Error *errorResponse
}

func newJobNotFoundResponse(err error) *jobNotFoundResponse {
	return &jobNotFoundResponse{Error: newErrorResponse(err).withStatus(http.StatusNotFound)}
}

func (r *jobNotFoundResponse) Result() (int, interface{}, error) {
	return r.Error.Result()
}

// swagger:response jobNotFoundInTheProvider
type jobNotFoundProviderResponse struct {
	// in: body
	Error *errorResponse
}

func newJobNotFoundProviderResponse(err error) *jobNotFoundProviderResponse {
	return &jobNotFoundProviderResponse{Error: newErrorResponse(err).withStatus(http.StatusGone)}
}

func (r *jobNotFoundProviderResponse) Result() (int, interface{}, error) {
	return r.Error.Result()
}
