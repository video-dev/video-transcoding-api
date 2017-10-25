package encodingcom

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// StatusResponse is the result of the GetStatus method.
//
// See http://goo.gl/NDsN8h for more details.
type StatusResponse struct {
	MediaID             string
	UserID              string
	SourceFile          string
	MediaStatus         string
	PreviousMediaStatus string
	NotifyURL           string
	CreateDate          time.Time
	StartDate           time.Time
	FinishDate          time.Time
	DownloadDate        time.Time
	UploadDate          time.Time
	TimeLeft            string
	Progress            float64
	TimeLeftCurrentJob  string
	ProgressCurrentJob  float64
	Formats             []FormatStatus
}

// FormatStatus is the status of each formatting input for a given MediaID.
//
// It is part of the StatusResponse type.
type FormatStatus struct {
	ID            string
	Status        string
	CreateDate    time.Time
	StartDate     time.Time
	FinishDate    time.Time
	Description   string
	S3Destination string
	CFDestination string
	Size          string
	Bitrate       string
	Output        string
	VideoCodec    string
	AudioCodec    string
	Destinations  []DestinationStatus
	Stream        []Stream
	FileSize      string
}

// DestinationStatus represents the status of a given destination.
type DestinationStatus struct {
	Name   string
	Status string
}

// GetStatus returns the status of the given media ids, it returns an slice of
// StatusResponse, the size of the result slice matches the size of input
// slice.
func (c *Client) GetStatus(mediaIDs []string, extended bool) ([]StatusResponse, error) {
	if len(mediaIDs) == 0 {
		return nil, errors.New("please provide at least one media id")
	}

	var m map[string]map[string]interface{}
	err := c.do(&request{
		Action:   "GetStatus",
		MediaID:  strings.Join(mediaIDs, ","),
		Extended: YesNoBoolean(extended),
	}, &m)
	if err != nil {
		return nil, err
	}

	var apiStatus []statusJSON
	var rawData interface{}
	if extended {
		rawData = m["response"]["job"]
	} else {
		rawData = m["response"]
	}
	rawBytes, _ := json.Marshal(rawData)
	if _, ok := rawData.([]interface{}); ok {
		json.Unmarshal(rawBytes, &apiStatus)
	} else {
		var item statusJSON
		json.Unmarshal(rawBytes, &item)
		apiStatus = append(apiStatus, item)
	}

	statusResponse := make([]StatusResponse, len(apiStatus))
	for i, status := range apiStatus {
		statusResponse[i] = status.toStruct()
	}
	return statusResponse, nil
}

type statusJSON struct {
	MediaID             string        `json:"id"`
	UserID              string        `json:"userid"`
	SourceFile          string        `json:"sourcefile"`
	MediaStatus         string        `json:"status"`
	PreviousMediaStatus string        `json:"prevstatus"`
	NotifyURL           string        `json:"notifyurl"`
	CreateDate          MediaDateTime `json:"created"`
	StartDate           MediaDateTime `json:"started"`
	FinishDate          MediaDateTime `json:"finished"`
	DownloadDate        MediaDateTime `json:"downloaded"`
	UploadDate          MediaDateTime `json:"uploaded"`
	TimeLeft            string        `json:"time_left"`
	Progress            float64       `json:"progress,string"`
	TimeLeftCurrentJob  string        `json:"time_left_current"`
	ProgressCurrentJob  float64       `json:"progress_current,string"`
	Formats             interface{}   `json:"format"`
}

func (s *statusJSON) toStruct() StatusResponse {
	resp := StatusResponse{
		MediaID:             s.MediaID,
		UserID:              s.UserID,
		SourceFile:          s.SourceFile,
		MediaStatus:         s.MediaStatus,
		PreviousMediaStatus: s.PreviousMediaStatus,
		NotifyURL:           s.NotifyURL,
		CreateDate:          s.CreateDate.Time,
		StartDate:           s.StartDate.Time,
		FinishDate:          s.FinishDate.Time,
		DownloadDate:        s.DownloadDate.Time,
		UploadDate:          s.UploadDate.Time,
		TimeLeft:            s.TimeLeft,
		Progress:            s.Progress,
		TimeLeftCurrentJob:  s.TimeLeftCurrentJob,
		ProgressCurrentJob:  s.ProgressCurrentJob,
	}

	// Yes, Encoding.com API is nuts, and when there's a single item in the
	// list, it returns an object instead of a list with a single item, so
	// we marshal it back, and then unmarshal in the proper type. The same
	// happens in the internal list of destinations.
	var formats []formatStatusJSON
	data, _ := json.Marshal(s.Formats)
	if _, ok := s.Formats.([]interface{}); ok {
		json.Unmarshal(data, &formats)
	} else {
		var formatStatus formatStatusJSON
		json.Unmarshal(data, &formatStatus)
		formats = append(formats, formatStatus)
	}

	resp.Formats = make([]FormatStatus, len(formats))
	for i, formatStatus := range formats {
		format := FormatStatus{
			ID:            formatStatus.ID,
			Status:        formatStatus.Status,
			CreateDate:    formatStatus.CreateDate.Time,
			StartDate:     formatStatus.StartDate.Time,
			FinishDate:    formatStatus.FinishDate.Time,
			Description:   formatStatus.Description,
			S3Destination: formatStatus.S3Destination,
			CFDestination: formatStatus.CFDestination,
			Size:          formatStatus.Size,
			Bitrate:       formatStatus.Bitrate,
			AudioCodec:    formatStatus.AudioCodec,
			Output:        formatStatus.Output,
			VideoCodec:    formatStatus.VideoCodec,
			Stream:        formatStatus.Stream,
			FileSize:      formatStatus.FileSize,
		}

		switch dest := formatStatus.Destinations.(type) {
		case string:
			destinationStatus := DestinationStatus{Name: dest}
			if statusStr, ok := formatStatus.DestinationsStatus.(string); ok {
				destinationStatus.Status = statusStr
			}
			format.Destinations = append(format.Destinations, destinationStatus)
		case []interface{}:
			destStats, ok := formatStatus.DestinationsStatus.([]interface{})
			if !ok {
				destStats = make([]interface{}, len(dest))
			}
			format.Destinations = make([]DestinationStatus, len(dest))
			for i, d := range dest {
				format.Destinations[i] = DestinationStatus{}
				if destName, ok := d.(string); ok {
					format.Destinations[i].Name = destName
				}
				if statusStr, ok := destStats[i].(string); ok {
					format.Destinations[i].Status = statusStr
				}
			}
		}

		resp.Formats[i] = format
	}
	return resp
}

type formatStatusJSON struct {
	ID                 string        `json:"id"`
	Status             string        `json:"status"`
	CreateDate         MediaDateTime `json:"created"`
	StartDate          MediaDateTime `json:"started"`
	FinishDate         MediaDateTime `json:"finished"`
	Description        string        `json:"description"`
	S3Destination      string        `json:"s3_destination"`
	CFDestination      string        `json:"cf_destination"`
	Destinations       interface{}   `json:"destination"`
	DestinationsStatus interface{}   `json:"destination_status"`
	Size               string        `json:"size"`
	Bitrate            string        `json:"bitrate"`
	AudioCodec         string        `json:"audio_codec"`
	VideoCodec         string        `json:"video_codec"`
	Output             string        `json:"output"`
	Stream             []Stream      `json:"stream"`
	FileSize           string        `json:"convertedsize"`
}
