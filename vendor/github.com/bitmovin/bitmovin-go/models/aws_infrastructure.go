package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type CreateAWSInfrastructureRequest struct {
	Name          *string `json:"name,omitempty"`
	Description   *string `json:"description,omitempty"`
	AccessKey     string  `json:"accessKey"`
	SecretKey     string  `json:"secretKey"`
	AccountNumber string  `json:"accountNumber"`
}

type AWSInfrastructureDetail struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ID          string  `json:"id"`
	CreatedAt   *string `json:"createdAt,omitempty"`
	UpdatedAt   *string `json:"updatedAt,omitempty"`
}

type AWSInfrastructureResponseData struct {
	Result AWSInfrastructureDetail `json:"result"`
}

type AWSInfrastructureResponse struct {
	RequestID *string                       `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus  `json:"status,omitempty"`
	Data      AWSInfrastructureResponseData `json:"data,omitempty"`
}
type AWSInfrastructureListResult struct {
	TotalCount *int64                    `json:"totalCount,omitempty"`
	Previous   *string                   `json:"previous,omitempty"`
	Next       *string                   `json:"next,omitempty"`
	Items      []AWSInfrastructureDetail `json:"items,omitempty"`
}
type AWSInfrastructureListResponseData struct {
	Result AWSInfrastructureListResult `json:"result,omitempty"`
}
type AWSInfrastructureListResponse struct {
	RequestID *string                           `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus      `json:"status,omitempty"`
	Data      AWSInfrastructureListResponseData `json:"data,omitempty"`
}
