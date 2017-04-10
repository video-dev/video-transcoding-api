package encodingcom

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/NYTimes/encoding-wrapper/encodingcom"
)

const encodingComDateFormat = "2006-01-02 15:04:05"

var errMediaNotFound = errors.New("media not found")

type request struct {
	Action  string               `json:"action"`
	Name    string               `json:"name"`
	MediaID string               `json:"mediaid"`
	Source  []string             `json:"source"`
	Format  []encodingcom.Format `json:"format"`
}

type errorResponse struct {
	Message string    `json:"message"`
	Errors  errorList `json:"errors"`
}

type errorList struct {
	Error string `json:"error"`
}

type fakePreset struct {
	Name      string
	GivenName string
	Request   request
}

type fakeMedia struct {
	ID       string
	Request  request
	Created  time.Time
	Started  time.Time
	Finished time.Time
	Size     string
	Rotation int
	Status   string
}

// encodingComFakeServer is a fake version of the Encoding.com API.
type encodingComFakeServer struct {
	*httptest.Server
	medias  map[string]*fakeMedia
	presets map[string]*fakePreset
	status  *encodingcom.APIStatusResponse
}

func newEncodingComFakeServer() *encodingComFakeServer {
	server := encodingComFakeServer{
		medias:  make(map[string]*fakeMedia),
		presets: make(map[string]*fakePreset),
		status:  &encodingcom.APIStatusResponse{StatusCode: "ok", Status: "Ok"},
	}
	server.Server = httptest.NewServer(&server)
	return &server
}

func (s *encodingComFakeServer) SetAPIStatus(status *encodingcom.APIStatusResponse) {
	s.status = status
}

func (s *encodingComFakeServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/status.php" {
		s.apiStatus(w, r)
		return
	}
	requestData := r.FormValue("json")
	if requestData == "" {
		s.Error(w, "json is required")
		return
	}
	var m map[string]request
	err := json.Unmarshal([]byte(requestData), &m)
	if err != nil {
		s.Error(w, err.Error())
	}
	req := m["query"]
	switch req.Action {
	case "AddMedia":
		s.addMedia(w, req)
	case "CancelMedia":
		s.cancelMedia(w, req)
	case "GetStatus":
		s.getStatus(w, req)
	case "GetMediaInfo":
		s.getMediaInfo(w, req)
	case "GetPreset":
		s.getPreset(w, req)
	case "SavePreset":
		s.savePreset(w, req)
	case "DeletePreset":
		s.deletePreset(w, req)
	default:
		s.Error(w, "invalid action")
	}
}

func (s *encodingComFakeServer) apiStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.status)
}

func (s *encodingComFakeServer) addMedia(w http.ResponseWriter, req request) {
	id := generateID()
	created := time.Now().UTC()
	s.medias[id] = &fakeMedia{
		ID:       id,
		Request:  req,
		Created:  created,
		Started:  created.Add(time.Second),
		Size:     "1920x1080",
		Rotation: 90,
	}
	resp := map[string]encodingcom.AddMediaResponse{
		"response": {MediaID: id, Message: "it worked"},
	}
	json.NewEncoder(w).Encode(resp)
}

func (s *encodingComFakeServer) cancelMedia(w http.ResponseWriter, req request) {
	media, err := s.getMedia(req.MediaID)
	if err != nil {
		s.Error(w, err.Error())
		return
	}
	media.Status = "Canceled"
	resp := map[string]map[string]interface{}{
		"response": {"message": "Deleted"},
	}
	json.NewEncoder(w).Encode(resp)
}

func (s *encodingComFakeServer) getMediaInfo(w http.ResponseWriter, req request) {
	media, err := s.getMedia(req.MediaID)
	if err != nil {
		s.Error(w, err.Error())
		return
	}
	format := media.Request.Format[0]
	resp := map[string]map[string]interface{}{
		"response": {
			"duration":    "183",
			"size":        media.Size,
			"video_codec": format.VideoCodec,
			"rotation":    strconv.Itoa(media.Rotation),
		},
	}
	json.NewEncoder(w).Encode(resp)
}

func (s *encodingComFakeServer) getStatus(w http.ResponseWriter, req request) {
	media, err := s.getMedia(req.MediaID)
	if err != nil {
		s.Error(w, err.Error())
		return
	}
	now := time.Now().UTC().Truncate(time.Second)
	status := "Saving"
	if media.Status == "Canceled" {
		status = media.Status
	} else if media.Status != "Finished" && now.Sub(media.Started) > time.Second {
		if media.Finished.IsZero() {
			media.Finished = now
		}
		status = "Finished"
		media.Status = status
	} else if media.Status != "" {
		status = media.Status
	}
	resp := map[string]map[string]interface{}{
		"response": {
			"id":         media.ID,
			"sourcefile": "http://some.source.file",
			"userid":     "someuser",
			"status":     status,
			"progress":   "100.0",
			"time_left":  "1",
			"created":    media.Created.Format(encodingComDateFormat),
			"started":    media.Started.Format(encodingComDateFormat),
			"finished":   media.Finished.Format(encodingComDateFormat),
			"output":     hlsOutput,
			"format": map[string]interface{}{
				"destination": []string{
					"https://mybucket.s3.amazonaws.com/dir/job-123/some_hls_preset/video-0.m3u8",
					"https://mybucket.s3.amazonaws.com/dir/job-123/video.m3u8",
				},
				"destination_status": []string{"Saved", "Saved"},
				"convertedsize":      "45674",
				"size":               media.Request.Format[0].Size,
				"bitrate":            media.Request.Format[0].Bitrate,
				"output":             media.Request.Format[0].Output[0],
				"video_codec":        media.Request.Format[0].VideoCodec,
				"stream": []map[string]interface{}{
					{
						"sub_path": "some_hls_preset",
					},
				},
			},
		},
	}
	json.NewEncoder(w).Encode(resp)
}

func (s *encodingComFakeServer) savePreset(w http.ResponseWriter, req request) {
	presetName := req.Name
	if presetName == "" {
		presetName = generateID()
	}
	s.presets[presetName] = &fakePreset{GivenName: req.Name, Request: req}
	resp := map[string]map[string]string{
		"response": {
			"SavedPreset": presetName,
		},
	}
	json.NewEncoder(w).Encode(resp)
}

func (s *encodingComFakeServer) getPreset(w http.ResponseWriter, req request) {
	preset, ok := s.presets[req.Name]
	if !ok {
		s.Error(w, req.Name+" preset not found")
		return
	}
	resp := map[string]*encodingcom.Preset{
		"response": {
			Name:   req.Name,
			Format: convertFormat(preset.Request.Format[0]),
			Output: preset.Request.Format[0].Output[0],
			Type:   encodingcom.UserPresets,
		},
	}
	json.NewEncoder(w).Encode(resp)
}

func (s *encodingComFakeServer) deletePreset(w http.ResponseWriter, req request) {
	if _, ok := s.presets[req.Name]; !ok {
		s.Error(w, "preset not found")
		return
	}
	delete(s.presets, req.Name)
	resp := map[string]*encodingcom.Response{"response": {Message: "Deleted"}}
	json.NewEncoder(w).Encode(resp)
}

func (s *encodingComFakeServer) Error(w http.ResponseWriter, message string) {
	m := map[string]errorResponse{"response": {
		Errors: errorList{Error: message},
	}}
	json.NewEncoder(w).Encode(m)
}

func (s *encodingComFakeServer) getMedia(id string) (*fakeMedia, error) {
	media, ok := s.medias[id]
	if !ok {
		return nil, errMediaNotFound
	}
	return media, nil
}

func generateID() string {
	var id [8]byte
	rand.Read(id[:])
	return fmt.Sprintf("%x", id[:])
}

func convertFormat(format encodingcom.Format) encodingcom.PresetFormat {
	videoCodecParams, err := json.Marshal(format.VideoCodecParameters)
	if err != nil {
		log.Println(err.Error())
		return encodingcom.PresetFormat{}
	}
	keyframe := ""
	if len(format.Keyframe) > 0 {
		keyframe = format.Keyframe[0]
	}
	return encodingcom.PresetFormat{
		NoiseReduction:          format.NoiseReduction,
		Output:                  format.Output[0],
		VideoCodec:              format.VideoCodec,
		AudioCodec:              format.AudioCodec,
		Bitrate:                 format.Bitrate,
		AudioBitrate:            format.AudioBitrate,
		AudioSampleRate:         format.AudioSampleRate,
		AudioChannelsNumber:     format.AudioChannelsNumber,
		AudioVolume:             format.AudioVolume,
		Size:                    format.Size,
		FadeIn:                  format.FadeIn,
		FadeOut:                 format.FadeOut,
		CropLeft:                format.CropLeft,
		CropTop:                 format.CropTop,
		CropRight:               format.CropRight,
		CropBottom:              format.CropBottom,
		KeepAspectRatio:         format.KeepAspectRatio,
		SetAspectRatio:          format.SetAspectRatio,
		AddMeta:                 format.AddMeta,
		Hint:                    format.Hint,
		RcInitOccupancy:         format.RcInitOccupancy,
		MinRate:                 format.MinRate,
		MaxRate:                 format.MaxRate,
		BufSize:                 format.BufSize,
		Keyframe:                keyframe,
		Start:                   format.Start,
		Duration:                format.Duration,
		ForceKeyframes:          format.ForceKeyframes,
		Bframes:                 format.Bframes,
		Gop:                     format.Gop,
		Metadata:                format.Metadata,
		SegmentDuration:         string(format.SegmentDuration),
		Logo:                    format.Logo,
		VideoCodecParameters:    string(videoCodecParams),
		Profile:                 format.Profile,
		TwoPass:                 format.TwoPass,
		Turbo:                   format.Turbo,
		TwinTurbo:               format.TwinTurbo,
		Rotate:                  format.Rotate,
		SetRotate:               format.SetRotate,
		AudioSync:               format.AudioSync,
		VideoSync:               format.VideoSync,
		ForceInterlaced:         format.ForceInterlaced,
		StripChapters:           format.StripChapters,
		Framerate:               format.Framerate,
		FramerateUpperThreshold: format.FramerateUpperThreshold,
	}
}
