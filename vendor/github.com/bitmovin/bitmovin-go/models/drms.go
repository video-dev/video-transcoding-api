package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type PlayReadyDrm struct {
	ID          *string                        `json:"id,omitempty"`
	Name        *string                        `json:"name,omitempty"`
	Description *string                        `json:"description,omitempty"`
	CustomData  map[string]interface{}         `json:"customData,omitempty"`
	Key         *string                        `json:"key,omitempty"`
	KID         *string                        `json:"kid,omitempty"`
	KeySeed     *string                        `json:"keySeed,omitempty"`
	LaUrl       *string                        `json:"laUrl,omitempty"`
	Method      bitmovintypes.EncryptionMethod `json:"method,omitempty"`
	Outputs     []Output                       `json:"outputs,omitempty"`
}

func (p *PlayReadyDrm) AddOutput(output *Output) {
	p.Outputs = append(p.Outputs, *output)
}

type FairPlayDrm struct {
	ID          *string                `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	CustomData  map[string]interface{} `json:"customData,omitempty"`
	Key         *string                `json:"key,omitempty"`
	IV          *string                `json:"iv,omitempty"`
	URI         *string                `json:"uri,omitempty"`
	Outputs     []Output               `json:"outputs,omitempty"`
}

func (p *FairPlayDrm) AddOutput(output *Output) {
	p.Outputs = append(p.Outputs, *output)
}

type WidevineDrm struct {
	ID          *string                `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	CustomData  map[string]interface{} `json:"customData,omitempty"`
	Key         *string                `json:"key,omitempty"`
	KID         *string                `json:"kid,omitempty"`
	PSSH        *string                `json:"pssh,omitempty"`
	Outputs     []Output               `json:"outputs,omitempty"`
}

func (p *WidevineDrm) AddOutput(output *Output) {
	p.Outputs = append(p.Outputs, *output)
}

type WidevineCencDrm struct {
	PSSH *string `json:"pssh,omitempty"`
}

type PlayReadyCencDrm struct {
	LaURL *string `json:"laUrl,omitEmpty"`
	PSSH  *string `json:"pssh,omitEmpty"`
}

type CencDrm struct {
	ID          *string                `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	CustomData  map[string]interface{} `json:"customData,omitempty"`
	Key         *string                `json:"key,omitempty"`
	KID         *string                `json:"kid,omitempty"`
	Outputs     []Output               `json:"outputs,omitempty"`
	Widevine    WidevineCencDrm        `json:"widevine,omitempty"`
	PlayReady   PlayReadyCencDrm       `json:"playReady,omitEmpty"`
}

func (p *CencDrm) AddOutput(output *Output) {
	p.Outputs = append(p.Outputs, *output)
}

type DrmResponseData struct {
	Messages []Message `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type WidevineDrmData struct {
	DrmResponseData
	Result WidevineDrm `json:"result,omitempty"`
}

type FairPlayDrmData struct {
	DrmResponseData
	Result FairPlayDrm `json:"result,omitempty"`
}

type PlayReadyDrmData struct {
	DrmResponseData
	Result PlayReadyDrm `json:"result,omitempty"`
}

type CencDrmData struct {
	DrmResponseData
	Result CencDrm `json:"result,omitempty"`
}

type DrmResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
}

type WidevineDrmResponse struct {
	DrmResponse
	Data WidevineDrmData `json:"data,omitempty"`
}

type FairPlayDrmResponse struct {
	DrmResponse
	Data FairPlayDrmData `json:"data,omitempty"`
}

type PlayReadyDrmResponse struct {
	DrmResponse
	Data PlayReadyDrmData `json:"data,omitempty"`
}

type CencDrmResponse struct {
	DrmResponse
	Data CencDrmData `json:"data,omitempty"`
}
