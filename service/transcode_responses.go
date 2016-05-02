package service

import (
	"net/http"

	"github.com/nytm/video-transcoding-api/provider"
	"github.com/nytm/video-transcoding-api/swagger"
)

// PartialJob is the simple response given to an API
// call that creates a new transcoding job
// swagger:model
type PartialJob struct {
	// unique identifier of the job
	//
	// unique: true
	JobID string `json:"jobId"`
}

// JSON-encoded version of the Job, includes only the id of the job, that can
// be used for querying the current status of the job.
//
// swagger:response job
type jobResponse struct {
	// in: body
	Payload *PartialJob

	baseResponse
}

func newJobResponse(jobID string) *jobResponse {
	return &jobResponse{
		baseResponse: baseResponse{
			payload: &PartialJob{JobID: jobID},
			status:  http.StatusOK,
		},
	}
}

// JSON-encoded JobStatus, containing status information given by the
// underlying provider.
//
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

// error returned when the given job data is not valid.
//
// swagger:response invalidJob
type invalidJobResponse struct {
	// in: body
	Error *swagger.ErrorResponse
}

func newInvalidJobResponse(err error) *invalidJobResponse {
	return &invalidJobResponse{Error: swagger.NewErrorResponse(err).WithStatus(http.StatusBadRequest)}
}

func (r *invalidJobResponse) Result() (int, interface{}, error) {
	return r.Error.Result()
}

// error returned the given job id could not be found on the API.
//
// swagger:response jobNotFound
type jobNotFoundResponse struct {
	// in: body
	Error *swagger.ErrorResponse
}

func newJobNotFoundResponse(err error) *jobNotFoundResponse {
	return &jobNotFoundResponse{Error: swagger.NewErrorResponse(err).WithStatus(http.StatusNotFound)}
}

func (r *jobNotFoundResponse) Result() (int, interface{}, error) {
	return r.Error.Result()
}

// error returned when the given job id could not be found on the underlying
// provider.
//
// swagger:response jobNotFoundInTheProvider
type jobNotFoundProviderResponse struct {
	// in: body
	Error *swagger.ErrorResponse
}

func newJobNotFoundProviderResponse(err error) *jobNotFoundProviderResponse {
	return &jobNotFoundProviderResponse{Error: swagger.NewErrorResponse(err).WithStatus(http.StatusGone)}
}

func (r *jobNotFoundProviderResponse) Result() (int, interface{}, error) {
	return r.Error.Result()
}
