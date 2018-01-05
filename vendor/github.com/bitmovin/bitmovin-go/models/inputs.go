package models

import (
	"github.com/bitmovin/bitmovin-go/bitmovintypes"
)

type AsperaInput struct {
	ID           *string                `json:"id,omitempty"`
	Name         *string                `json:"name,omitempty"`
	Description  *string                `json:"description,omitempty"`
	CustomData   map[string]interface{} `json:"customData,omitempty"`
	Host         *string                `json:"host,omitempty"`
	UserName     *string                `json:"username,omitempty"`
	Password     *string                `json:"password,omitempty"`
	MinBandwidth *string                `json:"minBandwidth,omitempty"`
	MaxBandwidth *string                `json:"maxBandwidth,omitempty"`
}

type AzureInput struct {
	ID          *string                `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	CustomData  map[string]interface{} `json:"customData,omitempty"`
	AccountName *string                `json:"accountName,omitempty"`
	AccountKey  *string                `json:"accountKey,omitempty"`
	Container   *string                `json:"container,omitempty"`
}

type FTPInput struct {
	ID          *string                `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	CustomData  map[string]interface{} `json:"customData,omitempty"`
	Host        *string                `json:"host,omitempty"`
	UserName    *string                `json:"username,omitempty"`
	Password    *string                `json:"password,omitempty"`
	Passive     *bool                  `json:"passive,omitempty"`
}

type GCSInput struct {
	ID          *string                         `json:"id,omitempty"`
	Name        *string                         `json:"name,omitempty"`
	Description *string                         `json:"description,omitempty"`
	CustomData  map[string]interface{}          `json:"customData,omitempty"`
	AccessKey   *string                         `json:"accessKey,omitempty"`
	SecretKey   *string                         `json:"secretKey,omitempty"`
	BucketName  *string                         `json:"bucketName,omitempty"`
	CloudRegion bitmovintypes.GoogleCloudRegion `json:"cloudRegion,omitempty"`
}

type GCSInputItem struct {
	ID          *string                      `json:"id,omitempty"`
	Name        *string                      `json:"name,omitempty"`
	Description *string                      `json:"description,omitempty"`
	BucketName  *string                      `json:"bucketName,omitempty"`
	CloudRegion bitmovintypes.AWSCloudRegion `json:"cloudRegion,omitempty"`
	CreatedAt   *string                      `json:"createdAt,omitempty"`
	UpdatedAt   *string                      `json:"updatedAt,omitempty"`
}

type GCSInputData struct {
	//Success fields
	Result   GCSInputItem `json:"result,omitempty"`
	Messages []Message    `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type GCSInputResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      GCSInputData                 `json:"data,omitempty"`
}

type GCSInputListResult struct {
	TotalCount *int64         `json:"totalCount,omitempty"`
	Previous   *string        `json:"previous,omitempty"`
	Next       *string        `json:"next,omitempty"`
	Items      []GCSInputItem `json:"items,omitempty"`
}

type GCSInputListData struct {
	Result GCSInputListResult `json:"result,omitempty"`
}

type GCSInputListResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      GCSInputListData             `json:"data,omitempty"`
}

type HTTPInput struct {
	ID          *string                `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	CustomData  map[string]interface{} `json:"customData,omitempty"`
	Host        *string                `json:"host,omitempty"`
	UserName    *string                `json:"username,omitempty"`
	Password    *string                `json:"password,omitempty"`
}

type HTTPInputItem struct {
	ID          *string `json:"id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Host        *string `json:"host,omitempty"`
	CreatedAt   *string `json:"createdAt,omitempty"`
	UpdatedAt   *string `json:"updatedAt,omitempty"`
}

type HTTPInputData struct {
	//Success fields
	Result   HTTPInputItem `json:"result,omitempty"`
	Messages []Message     `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type HTTPInputResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      HTTPInputData                `json:"data,omitempty"`
}

type HTTPInputListResult struct {
	TotalCount *int64          `json:"totalCount,omitempty"`
	Previous   *string         `json:"previous,omitempty"`
	Next       *string         `json:"next,omitempty"`
	Items      []HTTPInputItem `json:"items,omitempty"`
}

type HTTPInputListData struct {
	Result HTTPInputListResult `json:"result,omitempty"`
}

type HTTPInputListResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      HTTPInputListData            `json:"data,omitempty"`
}

type HTTPSInput struct {
	ID          *string                `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	CustomData  map[string]interface{} `json:"customData,omitempty"`
	Host        *string                `json:"host,omitempty"`
	UserName    *string                `json:"username,omitempty"`
	Password    *string                `json:"password,omitempty"`
}

type HTTPSInputItem struct {
	ID          *string `json:"id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Host        *string `json:"host,omitempty"`
	CreatedAt   *string `json:"createdAt,omitempty"`
	UpdatedAt   *string `json:"updatedAt,omitempty"`
}

type HTTPSInputData struct {
	//Success fields
	Result   HTTPSInputItem `json:"result,omitempty"`
	Messages []Message      `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type HTTPSInputResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      HTTPSInputData               `json:"data,omitempty"`
}

type HTTPSInputListResult struct {
	TotalCount *int64           `json:"totalCount,omitempty"`
	Previous   *string          `json:"previous,omitempty"`
	Next       *string          `json:"next,omitempty"`
	Items      []HTTPSInputItem `json:"items,omitempty"`
}

type HTTPSInputListData struct {
	Result HTTPSInputListResult `json:"result,omitempty"`
}

type HTTPSInputListResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      HTTPSInputListData           `json:"data,omitempty"`
}

type S3Input struct {
	ID          *string                      `json:"id,omitempty"`
	Name        *string                      `json:"name,omitempty"`
	Description *string                      `json:"description,omitempty"`
	CustomData  map[string]interface{}       `json:"customData,omitempty"`
	AccessKey   *string                      `json:"accessKey,omitempty"`
	SecretKey   *string                      `json:"secretKey,omitempty"`
	BucketName  *string                      `json:"bucketName,omitempty"`
	CloudRegion bitmovintypes.AWSCloudRegion `json:"cloudRegion,omitempty"`
}

type S3InputItem struct {
	ID          *string                      `json:"id,omitempty"`
	Name        *string                      `json:"name,omitempty"`
	Description *string                      `json:"description,omitempty"`
	BucketName  *string                      `json:"bucketName,omitempty"`
	CloudRegion bitmovintypes.AWSCloudRegion `json:"cloudRegion,omitempty"`
	CreatedAt   *string                      `json:"createdAt,omitempty"`
	UpdatedAt   *string                      `json:"updatedAt,omitempty"`
}

type S3InputData struct {
	//Success fields
	Result   S3InputItem `json:"result,omitempty"`
	Messages []Message   `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type S3InputResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      S3InputData                  `json:"data,omitempty"`
}

type S3InputListResult struct {
	TotalCount *int64        `json:"totalCount,omitempty"`
	Previous   *string       `json:"previous,omitempty"`
	Next       *string       `json:"next,omitempty"`
	Items      []S3InputItem `json:"items,omitempty"`
}

type S3InputListData struct {
	Result S3InputListResult `json:"result,omitempty"`
}

type S3InputListResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      S3InputListData              `json:"data,omitempty"`
}

type SFTPInput struct {
	ID          *string                `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	CustomData  map[string]interface{} `json:"customData,omitempty"`
	Host        *string                `json:"host,omitempty"`
	UserName    *string                `json:"username,omitempty"`
	Password    *string                `json:"password,omitempty"`
	Passive     *bool                  `json:"passive,omitempty"`
}

type RTMPInputItem struct {
	ID          *string `json:"id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	CreatedAt   *string `json:"createdAt,omitempty"`
	UpdatedAt   *string `json:"updatedAt,omitempty"`
}

type RTMPInputData struct {
	//Success fields
	Result   RTMPInputItem `json:"result,omitempty"`
	Messages []Message     `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type RTMPInputResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      RTMPInputData                `json:"data,omitempty"`
}

type RTMPInputListResult struct {
	TotalCount *int64          `json:"totalCount,omitempty"`
	Previous   *string         `json:"previous,omitempty"`
	Next       *string         `json:"next,omitempty"`
	Items      []RTMPInputItem `json:"items,omitempty"`
}

type RTMPInputListData struct {
	Result RTMPInputListResult `json:"result,omitempty"`
}

type RTMPInputListResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      RTMPInputListData            `json:"data,omitempty"`
}

type ZixiInput struct {
	ID             *string                `json:"id,omitempty"`
	Name           *string                `json:"name,omitempty"`
	Description    *string                `json:"description,omitempty"`
	CustomData     map[string]interface{} `json:"customData,omitempty"`
	Host           *string                `json:"host,omitempty"`
	Port           *int64                 `json:"port,omitempty"`
	Stream         *string                `json:"stream,omitempty"`
	Password       *string                `json:"password,omitempty"`
	Latency        *int64                 `json:"latency,omitempty"`
	MinBitrate     *int64                 `json:"minBitrate,omitempty"`
	DecryptionType *string                `json:"decryptionType,omitempty"`
	DecryptionKey  *string                `json:"decryptionKey,omitempty"`
}

type ZixiInputItem struct {
	ID             *string                `json:"id,omitempty"`
	Name           *string                `json:"name,omitempty"`
	Description    *string                `json:"description,omitempty"`
	CustomData     map[string]interface{} `json:"customData,omitempty"`
	Host           *string                `json:"host,omitempty"`
	Port           *int64                 `json:"port,omitempty"`
	Stream         *string                `json:"stream,omitempty"`
	Latency        *int64                 `json:"latency,omitempty"`
	MinBitrate     *int64                 `json:"minBitrate,omitempty"`
	DecryptionType *string                `json:"decryptionType,omitempty"`
}

type ZixiInputData struct {
	Result           ZixiInputItem `json:"result,omitempty"`
	Messages         []Message     `json:"messages,omitempty"`
	Code             *int64        `json:"code,omitempty"`
	Message          *string       `json:"message,omitempty"`
	DeveloperMessage *string       `json:"developerMessage,omitempty"`
	Links            []Link        `json:"links,omitempty"`
	Details          []Detail      `json:"details,omitempty"`
}

type ZixiInputResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      ZixiInputData                `json:"data,omitempty"`
}

type ZixiInputListResult struct {
	TotalCount *int64          `json:"totalCount,omitempty"`
	Previous   *string         `json:"previous,omitempty"`
	Next       *string         `json:"next,omitempty"`
	Items      []ZixiInputItem `json:"items,omitempty"`
}

type ZixiInputListData struct {
	Result ZixiInputListResult `json:"result,omitempty"`
}

type ZixiInputListResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      ZixiInputListData            `json:"data,omitempty"`
}
