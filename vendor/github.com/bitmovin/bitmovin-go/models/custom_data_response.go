package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type Data struct {
	Result Result `json:"result"`
}

type Result struct {
	CreatedAt  *string                `json:"createdAt"`
	ModifiedAt *string                `json:"modifiedAt"`
	CustomData map[string]interface{} `json:"customData"`
}

type CustomDataResponse struct {
	RequestID *string                      `json:"requestId"`
	Status    bitmovintypes.ResponseStatus `json:"status"`
	Data      Data                         `json:"data"`
}
