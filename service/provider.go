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
	Providers []provider.Descriptor

	baseResponse
}

func newListProvidersResponse(providers []provider.Descriptor) *listProvidersResponse {
	return &listProvidersResponse{
		baseResponse: baseResponse{
			payload: providers,
			status:  http.StatusOK,
		},
	}
}

// swagger:route GET /providers providers listProviders
//
// Describe available providers in the API, including their name, capabilities
// and health state.
//
//     Responses:
//       200: listProviders
//       500: genericError
func (s *TranscodingService) listProviders(r *http.Request) gizmoResponse {
	return newListProvidersResponse(provider.DescribeProviders(s.config))
}
