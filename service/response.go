package service

import "net/http"

type baseResponse struct {
	payload interface{}
	status  int
	err     error
}

func (r *baseResponse) Result() (int, interface{}, error) {
	return r.status, r.payload, r.err
}

// emptyResponse represents an empty response returned by the API, it's
// composed only by the HTTP status code.
//
// swagger:response
type emptyResponse int

func (r emptyResponse) Result() (int, interface{}, error) {
	return int(r), nil, nil
}

// errorReponse represents the basic error returned by the API on operation
// failures.
//
// swagger:response genericError
type errorResponse struct {
	// the error message
	//
	// in: body
	Message string `json:"error"`

	baseResponse
}

func newErrorResponse(err error) *errorResponse {
	errResp := &errorResponse{
		Message: err.Error(),
		baseResponse: baseResponse{
			status: http.StatusInternalServerError,
		},
	}
	errResp.err = errResp
	return errResp
}

func (r *errorResponse) withStatus(status int) *errorResponse {
	errResp := &errorResponse{Message: r.Message}
	baseResp := r.baseResponse
	baseResp.status = status
	baseResp.err = errResp
	errResp.baseResponse = baseResp
	return errResp
}

func (r *errorResponse) Error() string {
	return r.Message
}
