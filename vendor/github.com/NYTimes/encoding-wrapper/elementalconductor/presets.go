package elementalconductor

import "encoding/xml"

// GetPresets returns a list of presets
func (c *Client) GetPresets() (*PresetList, error) {
	var result *PresetList
	err := c.do("GET", "/presets", nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetPreset return details of a given presetID
func (c *Client) GetPreset(presetID string) (*Preset, error) {
	var result *Preset
	err := c.do("GET", "/presets/"+presetID, nil, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CreatePreset creates a new preset
func (c *Client) CreatePreset(preset *Preset) (*Preset, error) {
	var result *Preset
	err := c.do("POST", "/presets", preset, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeletePreset removes a preset based on its presetID
func (c *Client) DeletePreset(presetID string) error {
	return c.do("DELETE", "/presets/"+presetID, nil, nil)
}

// PresetList represents the response returned by
// a query for the list of jobs
type PresetList struct {
	Presets []Preset `xml:"preset"`
}

// Preset represents a preset
type Preset struct {
	XMLName       xml.Name `xml:"preset"`
	Name          string   `xml:"name"`
	Href          string   `xml:"href,attr,omitempty"`
	Permalink     string   `xml:"permalink,omitempty"`
	Description   string   `xml:"description,omitempty"`
	Container     string   `xml:"container,omitempty"`
	Width         string   `xml:"video_description>width,omitempty"`
	Height        string   `xml:"video_description>height,omitempty"`
	VideoCodec    string   `xml:"video_description>codec,omitempty"`
	VideoBitrate  string   `xml:"video_description>h264_settings>bitrate,omitempty"`
	GopSize       string   `xml:"video_description>h264_settings>gop_size,omitempty"`
	GopMode       string   `xml:"video_description>h264_settings>gop_mode,omitempty"`
	Profile       string   `xml:"video_description>h264_settings>profile,omitempty"`
	ProfileLevel  string   `xml:"video_description>h264_settings>level,omitempty"`
	RateControl   string   `xml:"video_description>h264_settings>rate_control_mode,omitempty"`
	InterlaceMode string   `xml:"video_description>h264_settings>interlace_mode,omitempty"`
	AudioCodec    string   `xml:"audio_description>codec,omitempty"`
	AudioBitrate  string   `xml:"audio_description>aac_settings>bitrate,omitempty"`
}
