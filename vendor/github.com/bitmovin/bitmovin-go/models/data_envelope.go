package models

type DataEnvelope struct {
	RequestID string `json:"requestId"`
	Status    string `json:"status"`
	Data      struct {
		Code             int    `json:"code"`
		Message          string `json:"message"`
		DeveloperMessage string `json:"developerMessage"`
		Links            []struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"links"`
		Details []struct {
			Type  string `json:"type"`
			Text  string `json:"text"`
			Field string `json:"field"`
		} `json:"details"`
	} `json:"data"`
}
