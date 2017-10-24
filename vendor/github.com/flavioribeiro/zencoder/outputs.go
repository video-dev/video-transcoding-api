package zencoder

import (
	"fmt"
)

// Get Output Details
func (z *Zencoder) GetOutputDetails(id int64) (*OutputMediaFile, error) {
	var details OutputMediaFile

	if err := z.getBody(fmt.Sprintf("outputs/%d.json", id), &details); err != nil {
		return nil, err
	}

	return &details, nil
}

// Output Progress
func (z *Zencoder) GetOutputProgress(id int64) (*FileProgress, error) {
	var details FileProgress

	if err := z.getBody(fmt.Sprintf("outputs/%d/progress.json", id), &details); err != nil {
		return nil, err
	}

	return &details, nil
}
