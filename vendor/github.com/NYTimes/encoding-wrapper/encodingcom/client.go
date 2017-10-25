// Package encodingcom provides types and methods for interacting with the
// Encoding.com API.
//
// You can get more details on the API at http://api.encoding.com/.
package encodingcom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Client is the basic type for interacting with the API. It provides methods
// matching the available actions in the API.
type Client struct {
	Endpoint string
	UserID   string
	UserKey  string
}

// NewClient creates a instance of the client type.
func NewClient(endpoint, userID, userKey string) (*Client, error) {
	return &Client{Endpoint: endpoint, UserID: userID, UserKey: userKey}, nil
}

// Response represents the generic response in the Encoding.com API. It doesn't
// include error information, as the client will smartly handle errors and
// return an instance of APIError when something goes wrong.
//
// See http://goo.gl/GBEn98 for more details.
type Response struct {
	Message string `json:"message,omitempty"`
}

func (c *Client) doMediaAction(mediaID string, action string) (*Response, error) {
	var result map[string]*Response
	err := c.do(&request{
		Action:  action,
		MediaID: mediaID,
	}, &result)
	if err != nil {
		return nil, err
	}
	return result["response"], nil
}

func (c *Client) do(r *request, out interface{}) error {
	r.UserID = c.UserID
	r.UserKey = c.UserKey
	jsonRequest, err := json.Marshal(r)
	if err != nil {
		return err
	}
	rawMsg := json.RawMessage(jsonRequest)
	m := map[string]interface{}{"query": &rawMsg}
	reqData, err := json.Marshal(m)
	if err != nil {
		return err
	}
	params := url.Values{}
	params.Add("json", string(reqData))
	req, err := http.NewRequest("POST", c.Endpoint, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var errRespWrapper map[string]*errorResponse
	err = json.Unmarshal(respData, &errRespWrapper)
	if err != nil {
		return fmt.Errorf("Error unmarshaling response: %s", err.Error())
	}
	if errResp := errRespWrapper["response"]; errResp.Errors.Error != "" {
		return &APIError{
			Message: errResp.Message,
			Errors:  []string{errResp.Errors.Error},
		}
	}
	return json.Unmarshal(respData, out)
}

// APIError represents an error returned by the Encoding.com API.
//
// See http://goo.gl/BzvXZt for more details.
type APIError struct {
	Message string `json:",omitempty"`
	Errors  []string
}

// Error converts the whole interlying information to a representative string.
//
// It encodes the list of errors in JSON format.
func (apiErr *APIError) Error() string {
	data, _ := json.Marshal(apiErr)
	return fmt.Sprintf("Error returned by the Encoding.com API: %s", data)
}

type errorResponse struct {
	Message string     `json:"message,omitempty"`
	Errors  errorsJSON `json:"errors,omitempty"`
}

type errorsJSON struct {
	Error string `json:"error,omitempty"`
}

type request struct {
	UserID                  string       `json:"userid"`
	UserKey                 string       `json:"userkey"`
	Action                  string       `json:"action"`
	MediaID                 string       `json:"mediaid,omitempty"`
	TaskID                  string       `json:"taskid,omitempty"`
	Source                  []string     `json:"source,omitempty"`
	SplitScreen             *SplitScreen `json:"split_screen,omitempty"`
	Region                  string       `json:"region,omitempty"`
	NotifyFormat            string       `json:"notify_format,omitempty"`
	NotifyURL               string       `json:"notify,omitempty"`
	NotifyEncodingErrorsURL string       `json:"notify_encoding_errors,omitempty"`
	NotifyUploadURL         string       `json:"notify_upload,omitempty"`
	Extended                YesNoBoolean `json:"extended,omitempty"`
	Type                    string       `json:"type,omitempty"`
	Name                    string       `json:"name,omitempty"`
	Format                  []Format     `json:"format,omitempty"`
}

// SplitScreen is the set of options for combining several sources to one split
// screen video.
//
// See http://goo.gl/EolKyv for more details.
type SplitScreen struct {
	Columns       int `json:"columns,string,omitempty"`
	Rows          int `json:"rows,string,omitempty"`
	PaddingLeft   int `json:"padding_left,string,omitempty"`
	PaddingRight  int `json:"padding_right,string,omitempty"`
	PaddingBottom int `json:"padding_bottom,string,omitempty"`
	PaddingTop    int `json:"padding_top,string,omitempty"`
}

// Format is the set of options for defining the output format when encoding
// new media files.
//
// See http://goo.gl/dcE1pF for more details.
type Format struct {
	Output                  []string             `json:"output,omitempty"`
	NoiseReduction          string               `json:"noise_reduction,omitempty"`
	OutputPreset            string               `json:"output_preset,omitempty"`
	VideoCodec              string               `json:"video_codec,omitempty"`
	AudioCodec              string               `json:"audio_codec,omitempty"`
	Bitrate                 string               `json:"bitrate,omitempty"`
	AudioBitrate            string               `json:"audio_bitrate,omitempty"`
	AudioChannelsNumber     string               `json:"audio_channels_number,omitempty"`
	Framerate               string               `json:"framerate,omitempty"`
	FramerateUpperThreshold string               `json:"framerate_upper_threshold,omitempty"`
	Size                    string               `json:"size,omitempty"`
	FadeIn                  string               `json:"fade_in,omitempty"`
	FadeOut                 string               `json:"fade_out,omitempty"`
	AudioSampleRate         uint                 `json:"audio_sample_rate,string,omitempty"`
	AudioVolume             uint                 `json:"audio_volume,string,omitempty"`
	CropLeft                int                  `json:"crop_left,string,omitempty"`
	CropTop                 int                  `json:"crop_top,string,omitempty"`
	CropRight               int                  `json:"crop_right,string,omitempty"`
	CropBottom              int                  `json:"crop_bottom,string,omitempty"`
	SetAspectRatio          string               `json:"set_aspect_ratio,omitempty"`
	RcInitOccupancy         string               `json:"rc_init_occupancy,omitempty"`
	MinRate                 string               `json:"minrate,omitempty"`
	MaxRate                 string               `json:"maxrate,omitempty"`
	BufSize                 string               `json:"bufsize,omitempty"`
	Keyframe                []string             `json:"keyframe,omitempty"`
	Start                   string               `json:"start,omitempty"`
	Duration                string               `json:"duration,omitempty"`
	ForceKeyframes          string               `json:"force_keyframes,omitempty"`
	Bframes                 int                  `json:"bframes,string,omitempty"`
	Gop                     string               `json:"gop,omitempty"`
	Metadata                *Metadata            `json:"metadata,omitempty"`
	Destination             []string             `json:"destination,omitempty"`
	SegmentDuration         uint                 `json:"segment_duration,omitempty"`
	Logo                    *Logo                `json:"logo,omitempty"`
	Overlay                 []Overlay            `json:"overlay,omitempty"`
	TextOverlay             []TextOverlay        `json:"text_overlay,omitempty"`
	VideoCodecParameters    VideoCodecParameters `json:"video_codec_parameters,omitempty"`
	Profile                 string               `json:"profile,omitempty"`
	Rotate                  string               `json:"rotate,omitempty"`
	SetRotate               string               `json:"set_rotate,omitempty"`
	AudioSync               string               `json:"audio_sync,omitempty"`
	VideoSync               string               `json:"video_sync,omitempty"`
	ForceInterlaced         string               `json:"force_interlaced,omitempty"`
	Stream                  []Stream             `json:"stream,omitempty"`
	AddMeta                 YesNoBoolean         `json:"add_meta,omitempty"`
	Hint                    YesNoBoolean         `json:"hint,omitempty"`
	KeepAspectRatio         YesNoBoolean         `json:"keep_aspect_ratio,omitempty"`
	StripChapters           YesNoBoolean         `json:"strip_chapters,omitempty"`
	TwoPass                 YesNoBoolean         `json:"two_pass,omitempty"`
	Turbo                   YesNoBoolean         `json:"turbo,omitempty"`
	TwinTurbo               YesNoBoolean         `json:"twin_turbo,omitempty"`
	PackFiles               *YesNoBoolean        `json:"pack_files,omitempty"`
}

// Stream is the set of options for defining Advanced HLS stream output
// when encoding new media files.
//
// See http://goo.gl/I7qRNo for more details.
type Stream struct {
	AudioBitrate            string       `json:"audio_bitrate,omitempty"`
	AudioChannelsNumber     string       `json:"audio_channels_number,omitempty"`
	AudioCodec              string       `json:"audio_codec,omitempty"`
	AudioSampleRate         uint         `json:"audio_sample_rate,string,omitempty"`
	AudioVolume             uint         `json:"audio_volume,string,omitempty"`
	Bitrate                 string       `json:"bitrate,omitempty"`
	Deinterlacing           string       `json:"deinterlacing,omitempty"`
	DownmixMode             string       `json:"downmix_mode,omitempty"`
	DurationPrecision       uint         `json:"duration_precision,string,omitempty"`
	Encoder                 string       `json:"encoder,omitempty"`
	EncryptionMethod        string       `json:"encryption_method,omitempty"`
	Framerate               uint         `json:"framerate,string,omitempty"`
	Keyframe                string       `json:"keyframe,omitempty"`
	MediaPath               string       `json:"media_path,omitempty"`
	PixFormat               string       `json:"pix_format,omitempty"`
	Profile                 string       `json:"profile,omitempty"`
	Rotate                  string       `json:"rotate,omitempty"`
	SetRotate               string       `json:"set_rotate,omitempty"`
	Size                    string       `json:"size,omitempty"`
	StillImageSize          string       `json:"still_image_size,omitempty"`
	StillImageTime          string       `json:"still_image_time,omitempty"`
	SubPath                 string       `json:"sub_path,omitempty"`
	VeryFast                string       `json:"veryfast,omitempty"`
	VideoCodec              string       `json:"video_codec,omitempty"`
	VideoSync               string       `json:"video_sync,omitempty"`
	VideoCodecParametersRaw interface{}  `json:"video_codec_parameters,omitempty"`
	AudioOnly               YesNoBoolean `json:"audio_only,omitempty"`
	AddIframeStream         YesNoBoolean `json:"add_iframe_stream,omitempty"`
	ByteRange               YesNoBoolean `json:"byte_range,omitempty"`
	Cbr                     YesNoBoolean `json:"cbr,omitempty"`
	CopyNielsenMetadata     YesNoBoolean `json:"copy_nielsen_metadata,omitempty"`
	CopyTimestamps          YesNoBoolean `json:"copy_timestamps,omitempty"`
	Encryption              YesNoBoolean `json:"encryption,omitempty"`
	HardCbr                 YesNoBoolean `json:"hard_cbr,omitempty"`
	Hint                    YesNoBoolean `json:"hint,omitempty"`
	KeepAspectRatio         YesNoBoolean `json:"keep_aspect_ratio,omitempty"`
	MetadataCopy            YesNoBoolean `json:"metadata_copy,omitempty"`
	StillImage              YesNoBoolean `json:"still_image,omitempty"`
	StripChapters           YesNoBoolean `json:"strip_chapters,omitempty"`
	TwoPass                 YesNoBoolean `json:"two_pass,omitempty"`
	VideoOnly               YesNoBoolean `json:"video_only,omitempty"`
}

// VideoCodecParameters function returns settings for H.264 video codec.
func (s Stream) VideoCodecParameters() VideoCodecParameters {
	var params VideoCodecParameters
	rawParameters, _ := json.Marshal(s.VideoCodecParametersRaw)
	json.Unmarshal(rawParameters, &params)
	return params
}

// VideoCodecParameters are settings for H.264 video codec.
//
// See http://goo.gl/8y7VSU for more details.
type VideoCodecParameters struct {
	Coder       string `json:"coder,omitempty"`
	Flags       string `json:"flags,omitempty"`
	Flags2      string `json:"flags2,omitempty"`
	Cmp         string `json:"cmp,omitempty"`
	Partitions  string `json:"partitions,omitempty"`
	MeMethod    string `json:"me_method,omitempty"`
	Subq        string `json:"subq,omitempty"`
	MeRange     string `json:"me_range,omitempty"`
	KeyIntMin   string `json:"keyint_min,omitempty"`
	ScThreshold string `json:"sc_threshold,omitempty"`
	Iqfactor    string `json:"i_qfactor,omitempty"`
	Bstrategy   string `json:"b_strategy,omitempty"`
	Qcomp       string `json:"qcomp,omitempty"`
	Qmin        string `json:"qmin,omitempty"`
	Qmax        string `json:"qmax,omitempty"`
	Qdiff       string `json:"qdiff,omitempty"`
	DirectPred  string `json:"directpred,omitempty"`
	Level       string `json:"level,omitempty"`
	Vprofile    string `json:"vprofile,omitempty"`
}

// Logo is the set of options for watermarking media during encoding, allowing
// users to add a image to the final media.
//
// See http://goo.gl/4z2Q5S for more details.
type Logo struct {
	LogoSourceURL string `json:"logo_source,omitempty"`
	LogoX         int    `json:"logo_x,string,omitempty"`
	LogoY         int    `json:"logo_y,string,omitempty"`
	LogoMode      int    `json:"logo_mode,string,omitempty"`
	LogoThreshold string `json:"logo_threshold,omitempty"`
}

// Overlay is the set of options for adding a video overlay in the media being
// encoded.
//
// See http://goo.gl/Q6sjkR for more details.
type Overlay struct {
	OverlaySource   string  `json:"overlay_source,omitempty"`
	OverlayLeft     string  `json:"overlay_left,omitempty"`
	OverlayRight    string  `json:"overlay_right,omitempty"`
	OverlayTop      string  `json:"overlay_top,omitempty"`
	OverlayBottom   string  `json:"overlay_bottom"`
	Size            string  `json:"size,omitempty"`
	OverlayStart    float64 `json:"overlay_start,string,omitempty"`
	OverlayDuration float64 `json:"overlay_duration,string,omitempty"`
}

// TextOverlay is the set of options for adding a text overlay in the media
// being encoded.
//
// See http://goo.gl/gUKi5t for more details.
type TextOverlay struct {
	Text            []string       `json:"text,omitempty"`
	FontSourceURL   string         `json:"font_source,omitempty"`
	FontSize        uint           `json:"font_size,string,omitempty"`
	FontRotate      int            `json:"font_rotate,string,omitempty"`
	FontColor       string         `json:"font_color,omitempty"`
	AlignCenter     ZeroOneBoolean `json:"align_center,omitempty"`
	OverlayX        int            `json:"overlay_x,string,omitempty"`
	OverlayY        int            `json:"overlay_y,string,omitempty"`
	Size            string         `json:"size,omitempty"`
	OverlayStart    float64        `json:"overlay_start,string,omitempty"`
	OverlayDuration float64        `json:"overlay_duration,string,omitempty"`
}

// Metadata represents media metadata, as provided in the Format struct when
// encoding new media.
//
// See http://goo.gl/jNSio9 for more details.
type Metadata struct {
	Title       string `json:"title,omitempty"`
	Copyright   string `json:"copyright,omitempty"`
	Author      string `json:"author,omitempty"`
	Description string `json:"description,omitempty"`
	Album       string `json:"album,omitempty"`
}

// YesNoBoolean is a boolean that turns true into "yes" and false into "no"
// when encoded as JSON.
type YesNoBoolean bool

// MarshalJSON is the method that ensures that YesNoBoolean satisfies the
// json.Marshaler interface.
func (b YesNoBoolean) MarshalJSON() ([]byte, error) {
	return boolToBytes(bool(b), "yes", "no"), nil
}

// UnmarshalJSON is the method that ensure that YesNoBoolean can be converted
// back from 1 or 0 to a boolean value.
func (b *YesNoBoolean) UnmarshalJSON(data []byte) error {
	v, err := bytesToBool(data, "yes", "no")
	if err != nil {
		return err
	}
	*b = YesNoBoolean(v)
	return nil
}

// ZeroOneBoolean is a boolean that turns true into "1" and false into "0" when
// encoded as JSON.
type ZeroOneBoolean bool

// MarshalJSON is the method that ensures that ZeroOneBoolean satisfies the
// json.Marshaler interface.
func (b ZeroOneBoolean) MarshalJSON() ([]byte, error) {
	return boolToBytes(bool(b), "1", "0"), nil
}

// UnmarshalJSON is the method that ensure that ZeroOneBoolean can be converted
// back from 1 or 0 to a boolean value.
func (b *ZeroOneBoolean) UnmarshalJSON(data []byte) error {
	v, err := bytesToBool(data, "1", "0")
	if err != nil {
		return err
	}
	*b = ZeroOneBoolean(v)
	return nil
}

func boolToBytes(b bool, t, f string) []byte {
	if b {
		return []byte(`"` + t + `"`)
	}
	return []byte(`"` + f + `"`)
}

func bytesToBool(data []byte, t, f string) (bool, error) {
	switch string(data) {
	case `"` + t + `"`:
		return true, nil
	case `"` + f + `"`:
		return false, nil
	default:
		return false, fmt.Errorf("invalid value: %s", data)
	}
}
