package zencoder

import (
	"fmt"
)

// Get Input Details
func (z *Zencoder) GetInputDetails(id int32) (*InputMediaFile, error) {
	var details InputMediaFile

	if err := z.getBody(fmt.Sprintf("inputs/%d.json", id), &details); err != nil {
		return nil, err
	}

	return &details, nil
}

// Input Progress
func (z *Zencoder) GetInputProgress(id int32) (*FileProgress, error) {
	var details FileProgress

	if err := z.getBody(fmt.Sprintf("inputs/%d/progress.json", id), &details); err != nil {
		return nil, err
	}

	return &details, nil
}
