package service

import (
	"net/http"

	"github.com/nytm/video-transcoding-api/provider"
)

// response for the listProviders operation. Contains the list of providers
// alphabetically ordered.
//
// swagger:response listProviders
type listProvidersResponse struct {
	// in: body
	Providers []string

	baseResponse
}

func newListProvidersResponse(providerNames []string) *listProvidersResponse {
	return &listProvidersResponse{
		baseResponse: baseResponse{payload: providerNames, status: http.StatusOK},
	}
}

// response for the getProvider operation.
//
// swagger:response provider
type getProviderResponse struct {
	// in: body
	Provider *provider.Description

	baseResponse
}

func newGetProviderResponse(p *provider.Description) *getProviderResponse {
	return &getProviderResponse{
		baseResponse: baseResponse{payload: p, status: http.StatusOK},
	}
}

// error returned when the given provider name is not found in the API.
//
// swagger:response providerNotFound
type providerNotFoundResponse struct {
	// in: body
	Error *errorResponse
}

func newProviderNotFoundResponse(err error) *providerNotFoundResponse {
	return &providerNotFoundResponse{Error: newErrorResponse(err).withStatus(http.StatusNotFound)}
}

func (r *providerNotFoundResponse) Result() (int, interface{}, error) {
	return r.Error.Result()
}
