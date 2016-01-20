package provider

import (
	"encoding/json"
	"errors"
	"fmt"
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

func TestRotationVariables(t *testing.T) {
	var tests = []struct {
		input    rotation
		expected uint
	}{
		{Rotate0Degrees, 0},
		{Rotate90Degrees, 90},
		{Rotate180Degrees, 180},
		{Rotate270Degrees, 270},
	}
	for _, test := range tests {
		if test.input.value != test.expected {
			t.Errorf("Rotation variables: Want %d. Got %d", test.expected, test.input.value)
		}
		if !test.input.set {
			t.Error("Rotation is not set, but it should be")
		}
	}
}

func TestRotationUnmarshalJSON(t *testing.T) {
	var tests = []struct {
		input       rotation
		expectedErr error
	}{
		{Rotate0Degrees, nil},
		{Rotate90Degrees, nil},
		{Rotate180Degrees, nil},
		{Rotate270Degrees, nil},
		{newRotation(500), errors.New("invalid value for rotation: 500")},
		{newRotation(3), errors.New("invalid value for rotation: 3")},
	}
	var m map[string]rotation
	for _, test := range tests {
		data := []byte(fmt.Sprintf(`{"value": %d}`, test.input.value))
		err := json.Unmarshal(data, &m)
		if err != nil && err.Error() != test.expectedErr.Error() {
			t.Fatal(err)
		}
		if test.expectedErr != nil {
			if err == nil {
				t.Errorf("unexpected nil error. Want %#v", test.expectedErr)
				continue
			} else if err.Error() != test.expectedErr.Error() {
				t.Errorf("wrong error: want %#v. Got %#v.", test.expectedErr, err)
				continue
			}
		} else if m["value"].value != test.input.value {
			t.Errorf("wrong value: want %d. Got %d.", test.input.value, m["value"].value)
		}
	}

	data := []byte(`{"value":null}`)
	err := json.Unmarshal(data, &m)
	expectedErrMsg := "invalid value for rotation: null. It must be a positive integer"
	if err == nil {
		t.Error("Unexpected nil error on invalid unmarshalling")
	} else if err.Error() != expectedErrMsg {
		t.Errorf("Wrong error returned. Want %q. Got %q.", expectedErrMsg, err.Error())
	}
}
