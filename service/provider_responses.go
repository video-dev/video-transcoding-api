package service

import (
	"net/http"

	"github.com/NYTimes/video-transcoding-api/provider"
	"github.com/NYTimes/video-transcoding-api/swagger"
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
	Error *swagger.ErrorResponse
}

func newProviderNotFoundResponse(err error) *providerNotFoundResponse {
	return &providerNotFoundResponse{Error: swagger.NewErrorResponse(err).WithStatus(http.StatusNotFound)}
}

func (r *providerNotFoundResponse) Result() (int, interface{}, error) {
	return r.Error.Result()
}
