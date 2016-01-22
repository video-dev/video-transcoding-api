package provider

import (
	"fmt"
	"strconv"
	"strings"
)

type rotation struct {
	value uint
	set   bool
}

var (
	// Rotate0Degrees represents 0 degrees rotation (no rotation).
	Rotate0Degrees = newRotation(0)

	// Rotate90Degrees represents 90 degrees rotation (rotate right).
	Rotate90Degrees = newRotation(90)

	// Rotate180Degrees represents 170 degrees rotation (two sequential
	// right rotations - or two sequential left rotations).
	Rotate180Degrees = newRotation(180)

	// Rotate270Degrees represents 270 degrees rotation (rotate left).
	Rotate270Degrees = newRotation(270)
)

func newRotation(n uint) rotation {
	return rotation{value: n, set: true}
}

func (r *rotation) UnmarshalJSON(b []byte) error {
	value, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid value for rotation: %s. It must be a positive integer", b)
	}
	switch value {
	case 0, 90, 180, 270:
		r.value = uint(value)
		r.set = true
		return nil
	default:
		r.set = false
		return fmt.Errorf("invalid value for rotation: %d", value)
	}
}

// Profile contains the set of options for transcoding a media.
type Profile struct {
	Output string
	Size   Size

	AudioCodec          string
	AudioBitRate        string
	AudioChannelsNumber string
	AudioSampleRate     uint

	BitRate         string
	FrameRate       string
	KeepAspectRatio bool
	VideoCodec      string

	KeyFrame    string
	AudioVolume uint
	Rotate      rotation
}

// Size represents the size of the media in pixels.
type Size struct {
	Width  uint64 `json:",omitempty"`
	Height uint64 `json:",omitempty"`
}

func (s *Size) String() string {
	return fmt.Sprintf("%dx%d", s.Width, s.Height)
}

// UnmarshalJSON for Size will parse a Size struct out of given slice of bytes
func (s *Size) UnmarshalJSON(b []byte) error {
	sizeString, err := strconv.Unquote(string(b))
	if err != nil {
		return fmt.Errorf("Size expression needs quotes")
	}
	sizeSlice := strings.Split(sizeString, "x")
	if len(sizeSlice) != 2 {
		return fmt.Errorf("Size expression is invalid (needs to be NxN): %s", sizeString)
	}
	width, err := strconv.ParseUint(sizeSlice[0], 10, 64)
	if err != nil {
		return fmt.Errorf("Size expression has invalid width: %s", sizeSlice[0])
	}
	height, err := strconv.ParseUint(sizeSlice[1], 10, 64)
	if err != nil {
		return fmt.Errorf("Size expression has invalid height: %s", sizeSlice[1])
	}
	s.Width = width
	s.Height = height
	return nil
}
