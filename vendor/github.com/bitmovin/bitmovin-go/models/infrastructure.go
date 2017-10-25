package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type CreateInfrastructureRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type InfrastructureDetail struct {
	Name                       string `json:"name"`
	Description                string `json:"description"`
	ID                         string `json:"id"`
	Online                     bool   `json:"online"`
	Connected                  bool   `json:"connected"`
	AgentDeploymentDownloadURL string `json:"agentDeploymentDownloadUrl"`
}

type InfrastructureResponseData struct {
	Result InfrastructureDetail `json:"Result"`
}

type InfrastructureResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      InfrastructureResponseData   `json:"data,omitempty"`
}
type InfrastructureListResult struct {
	TotalCount *int64                 `json:"totalCount,omitempty"`
	Previous   *string                `json:"previous,omitempty"`
	Next       *string                `json:"next,omitempty"`
	Items      []InfrastructureDetail `json:"items,omitempty"`
}
type InfrastructureListResponseData struct {
	Result InfrastructureListResult `json:"result,omitempty"`
}
type InfrastructureListResponse struct {
	RequestID *string                        `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus   `json:"status,omitempty"`
	Data      InfrastructureListResponseData `json:"data,omitempty"`
}
