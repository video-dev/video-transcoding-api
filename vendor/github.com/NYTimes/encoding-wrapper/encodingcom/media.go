package encodingcom

import "time"

// dateTimeLayout is the time layout used on Media items
const dateTimeLayout = "2006-01-02 15:04:05"

// MediaDateTime is a custom time struct to be used on Media items
type MediaDateTime struct {
	time.Time
}

// UnmarshalJSON implementation on MediaDateTime to use dateTimeLayout
func (mdt *MediaDateTime) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	if string(b) == "0000-00-00 00:00:00" {
		mdt.Time = time.Time{}
		return nil
	}
	mdt.Time, err = time.Parse(dateTimeLayout, string(b))
	return err
}

// AddMediaResponse represents the response returned by the AddMedia action.
//
// See http://goo.gl/Aqg8lc for more details.
type AddMediaResponse struct {
	Message string `json:"message,omitempty"`
	MediaID string `json:"mediaid,omitempty"`
}

// ListMediaResponse represents the response returned by the GetMediaList action.
//
// See http://goo.gl/xhVV6v for more details.
type ListMediaResponse struct {
	Media []ListMediaResponseItem `json:"media,omitempty"`
}

// ListMediaResponseItem represents each individual item returned by the GetMediaList action.
//
// See ListMediaResponse
type ListMediaResponseItem struct {
	MediaFile   string        `json:"mediafile,omitempty"`
	MediaID     string        `json:"mediaid,omitempty"`
	MediaStatus string        `json:"mediastatus,omitempty"`
	CreateDate  MediaDateTime `json:"createdate,string,omitempty"`
	StartDate   MediaDateTime `json:"startdate,string,omitempty"`
	FinishDate  MediaDateTime `json:"finishdate,string,omitempty"`
}

// AddMedia adds a new media to user's queue.
//
// Format specifies details on how the source files are going to be encoded.
//
// See http://goo.gl/whvHwJ for more details on the source file formatting.
func (c *Client) AddMedia(source []string, format []Format, region string) (*AddMediaResponse, error) {
	var result map[string]*AddMediaResponse
	req := request{
		Action: "AddMedia",
		Format: format,
		Source: source,
		Region: region,
	}
	err := c.do(&req, &result)
	if err != nil {
		return nil, err
	}
	return result["response"], nil
}

// StopMedia stops an existing media on user's queue based on the mediaID.
func (c *Client) StopMedia(mediaID string) (*Response, error) {
	return c.doMediaAction(mediaID, "StopMedia")
}

// CancelMedia deletes an existing media on user's queue based on the mediaID.
func (c *Client) CancelMedia(mediaID string) (*Response, error) {
	return c.doMediaAction(mediaID, "CancelMedia")
}

// RestartMedia restart the entire job of an existing media on user's queue based on the mediaID.
// When withErrors enabled it only retry tasks ended with error and not the entire job.
func (c *Client) RestartMedia(mediaID string, withErrors bool) (*Response, error) {
	action := "RestartMedia"
	if withErrors {
		action = "RestartMediaErrors"
	}
	return c.doMediaAction(mediaID, action)
}

// RestartMediaTask restart a specific task on a job.
func (c *Client) RestartMediaTask(mediaID string, taskID string) (*Response, error) {
	var result map[string]*Response
	err := c.do(&request{
		Action:  "RestartMediaTask",
		MediaID: mediaID,
		TaskID:  taskID,
	}, &result)
	if err != nil {
		return nil, err
	}
	return result["response"], nil
}

// ListMedia (GetMediaList action) returns a list of the user's media in the queue.
func (c *Client) ListMedia() (*ListMediaResponse, error) {
	var result map[string]*ListMediaResponse
	err := c.do(&request{Action: "GetMediaList"}, &result)
	if err != nil {
		return nil, err
	}
	return result["response"], nil
}

// MediaInfo is the result of the GetMediaInfo method.
//
// See http://goo.gl/OTX0Ua for more details.
type MediaInfo struct {
	Duration           time.Duration `json:"duration,string"`
	Bitrate            string        `json:"bitrate"`
	VideoCodec         string        `json:"video_codec"`
	VideoBitrate       string        `json:"video_bitrate"`
	Framerate          string        `json:"frame_rate"`
	Size               string        `json:"size"`
	PixelAspectRatio   string        `json:"pixel_aspect_ratio"`
	DisplayAspectRatio string        `json:"display_aspect_ratio"`
	AudioCodec         string        `json:"audio_codec"`
	AudioBitrate       string        `json:"audio_bitrate"`
	AudioSampleRate    uint          `json:"audio_sample_rate,string"`
	AudioChannels      string        `json:"audio_channels"`
	Rotation           uint          `json:"rotation,string"`
}

type mediaInfo struct {
	MediaInfo
	Duration float64 `json:"duration,string"`
}

// GetMediaInfo returns video parameters of the specified media when available.
func (c *Client) GetMediaInfo(mediaID string) (*MediaInfo, error) {
	var result map[string]*mediaInfo
	err := c.do(&request{Action: "GetMediaInfo", MediaID: mediaID}, &result)
	if err != nil {
		return nil, err
	}
	info := result["response"].MediaInfo
	info.Duration = time.Duration(result["response"].Duration * float64(time.Second))
	return &info, nil
}
