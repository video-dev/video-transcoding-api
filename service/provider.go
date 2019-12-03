package service

import (
	"net/http"

	"github.com/NYTimes/gizmo/server"
	"github.com/video-dev/video-transcoding-api/provider"
	"github.com/video-dev/video-transcoding-api/swagger"
)

// swagger:route GET /providers providers listProviders
//
// Describe available providers in the API, including their name, capabilities
// and health state.
//
//     Responses:
//       200: listProviders
//       500: genericError
func (s *TranscodingService) listProviders(*http.Request) swagger.GizmoJSONResponse {
	return newListProvidersResponse(provider.ListProviders(s.config))
}

// swagger:route GET /providers/{name} providers getProvider
//
// Describe available providers in the API, including their name, capabilities
// and health state.
//
//     Responses:
//       200: provider
//       404: providerNotFound
//       500: genericError
func (s *TranscodingService) getProvider(r *http.Request) swagger.GizmoJSONResponse {
	var params getProviderInput
	params.loadParams(server.Vars(r))
	description, err := provider.DescribeProvider(params.Name, s.config)
	switch err {
	case nil:
		return newGetProviderResponse(description)
	case provider.ErrProviderNotFound:
		return newProviderNotFoundResponse(err)
	default:
		return swagger.NewErrorResponse(err)
	}
}
