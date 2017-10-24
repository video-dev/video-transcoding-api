package hybrik

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	failed = "false"
)

type presetResponse struct {
	Success string `json:"success"`
	Msg     string `json:"message"`
}

// ErrCreatePreset occurs when there is a problem creating a preset
type ErrCreatePreset struct {
	Msg string
}

func (e ErrCreatePreset) Error() string {
	return fmt.Sprintf("unable to create preset, error: %s", e.Msg)
}

// ErrGetPreset occurs when there is a problem obtaining a preset
type ErrGetPreset struct {
	Msg string
}

func (e ErrGetPreset) Error() string {
	return fmt.Sprintf("unable to get preset, error: %s", e.Msg)
}

// GetPreset return details of a given presetID
func (c *Client) GetPreset(presetID string) (Preset, error) {

	result, err := c.client.CallAPI("GET", fmt.Sprintf("/presets/%s", presetID), nil, nil)
	if err != nil {
		return Preset{}, err
	}

	var preset Preset
	err = json.Unmarshal([]byte(result), &preset)
	if err != nil {
		return Preset{}, err
	}

	if preset.Name == "" {
		var pr presetResponse
		err = json.Unmarshal([]byte(result), &pr)
		if err != nil {
			return Preset{}, err
		}
		if pr.Success == failed {
			return Preset{}, ErrGetPreset{Msg: pr.Msg}
		}
	}

	return preset, nil
}

// CreatePreset creates a new preset
func (c *Client) CreatePreset(preset Preset) (Preset, error) {
	body, err := json.Marshal(preset)
	if err != nil {
		return Preset{}, err
	}

	resp, err := c.client.CallAPI("POST", "/presets", nil, bytes.NewReader(body))
	if err != nil {
		return Preset{}, err
	}

	var pr presetResponse
	err = json.Unmarshal([]byte(resp), &pr)
	if err != nil {
		return Preset{}, err
	}

	if pr.Success == failed {
		return Preset{}, ErrCreatePreset{Msg: pr.Msg}
	}

	return preset, nil
}

// DeletePreset removes a preset based on its presetID
func (c *Client) DeletePreset(presetID string) error {
	_, err := c.client.CallAPI("DELETE", fmt.Sprintf("/presets/%s", presetID), nil, nil)

	return err
}

// PresetList represents the response returned by
// a query for the list of jobs
type PresetList []Preset

// Preset represents a transcoding preset
type Preset struct {
	Key         string        `json:"key"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	UserData    string        `json:"user_data,omitempty"`
	Kind        string        `json:"kind"`
	Path        string        `json:"path"`
	Payload     PresetPayload `json:"payload"`
}

type PresetPayload struct {
	Targets []PresetTarget `json:"targets"`
}

type PresetTarget struct {
	FilePattern string `json:"file_pattern"`
	Container   struct {
		Kind string `json:"kind"`
	} `json:"container"`
	Video         VideoTarget   `json:"video,omitempty"`
	Audio         []AudioTarget `json:"audio,omitempty"`
	ExistingFiles string        `json:"existing_files,omitempty"`
	UID           string        `json:"uid,omitempty"`
}

type VideoTarget struct {
	Width           *int   `json:"width,omitempty"`
	Height          *int   `json:"height,omitempty"`
	BitrateMode     string `json:"bitrate_mode,omitempty"`
	BitrateKb       int    `json:"bitrate_kb,omitempty"`
	MaxBitrateKb    int    `json:"max_bitrate_kb,omitempty"`
	VbvBufferSizeKb int    `json:"vbv_buffer_size_kb,omitempty"`
	FrameRate       int    `json:"frame_rate,omitempty"`
	Codec           string `json:"codec,omitempty"`
	Profile         string `json:"profile,omitempty"`
	Level           string `json:"level,omitempty"`
	MinGOPFrames    int    `json:"min_gop_frames,omitempty"`
	MaxGOPFrames    int    `json:"max_gop_frames,omitempty"`
	UseClosedGOP    bool   `json:"use_closed_gop,omitempty"`
	InterlaceMode   string `json:"interlace_mode,omitempty"`
}

type AudioTarget struct {
	Codec      string `json:"codec,omitempty"`
	Channels   int    `json:"channels,omitempty"`
	SampleRate int    `json:"sample_rate,omitempty"`
	SampleSize int    `json:"sample_size,omitempty"`
	BitrateKb  int    `json:"bitrate_kb,omitempty"`
}
