package db

import (
	"errors"
	"testing"
)

func TestOutputOptionsValidation(t *testing.T) {
	var tests = []struct {
		testCase string
		opts     OutputOptions
		errMsg   string
	}{
		{
			"valid options",
			OutputOptions{Extension: "mp4"},
			"",
		},
		{
			"missing extension",
			OutputOptions{Extension: ""},
			"extension is required",
		},
	}
	for _, test := range tests {
		err := test.opts.Validate()
		if err == nil {
			err = errors.New("")
		}
		if err.Error() != test.errMsg {
			t.Errorf("wrong error message\nWant %q\nGot  %q", test.errMsg, err.Error())
		}
	}
}
