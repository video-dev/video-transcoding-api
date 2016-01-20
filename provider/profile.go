package provider

import "fmt"

type rotation *uint

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
	v := new(uint)
	*v = n
	return rotation(v)
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
	Width  uint `json:",omitempty"`
	Height uint `json:",omitempty"`
}

func (s *Size) String() string {
	return fmt.Sprintf("%dx%d", s.Width, s.Height)
}
