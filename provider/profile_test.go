package provider

import (
	"testing"
)

func TestSizeString(t *testing.T) {
	var tests = []struct {
		input    Size
		expected string
	}{
		{Size{Height: 360}, "0x360"},
		{Size{Height: 1080}, "0x1080"},
		{Size{Width: 1024}, "1024x0"},
		{Size{Width: 1920, Height: 1080}, "1920x1080"},
		{Size{}, "0x0"},
	}
	for _, test := range tests {
		output := test.input.String()
		if output != test.expected {
			t.Errorf("Size.String(): want %q. Got %q.", test.expected, output)
		}
	}
}
