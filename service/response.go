package service

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
