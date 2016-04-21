package swagger

import "net/http"

// ErrorResponse represents the basic error returned by the API on operation
// failures.
//
// swagger:response genericError
type ErrorResponse struct {
	// the error message
	//
	// in: body
	Message string `json:"error"`

	status int
}

// NewErrorResponse creates a new ErrorResponse with the given error. The
// default status code for error responses is 500 (InternalServerError). Use
// the method WithError to customize it.
func NewErrorResponse(err error) *ErrorResponse {
	errResp := &ErrorResponse{
		Message: err.Error(),
		status:  http.StatusInternalServerError,
	}
	return errResp
}

// WithStatus creates a new copy of ErrorResponse using the given status.
func (r *ErrorResponse) WithStatus(status int) *ErrorResponse {
	return &ErrorResponse{Message: r.Message, status: status}
}

// Error returns the underlying error message.
func (r *ErrorResponse) Error() string {
	return r.Message
}

// Result ensures that ErrorResponse implements the interface Handler.
func (r *ErrorResponse) Result() (int, interface{}, error) {
	return r.status, nil, r
}
