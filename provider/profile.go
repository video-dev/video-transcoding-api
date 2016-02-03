package provider

import (
	"fmt"
	"strconv"
)

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

// Rotation presents rotation configuration in the provider. It's composed by
// two interval values: the rotation in degrees and whether it has been
// defined.
type Rotation struct {
	value uint
	set   bool
}

// Value returns the underlying rotation data.
func (r Rotation) Value() (uint, bool) {
	return r.value, r.set
}

func newRotation(n uint) Rotation {
	return Rotation{value: n, set: true}
}

// UnmarshalJSON is the method used for unmarshaling a rotation instance from
// JSON format. It validates whether this is a valid rotation and  then set the
// values.
func (r *Rotation) UnmarshalJSON(b []byte) error {
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
	Output []string
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
	Rotate      Rotation

	TwoPassEncoding bool
}

// Size represents the size of the media in pixels.
type Size struct {
	Width  uint64 `json:",omitempty"`
	Height uint64 `json:",omitempty"`
}

func (s *Size) String() string {
	return fmt.Sprintf("%dx%d", s.Width, s.Height)
}
