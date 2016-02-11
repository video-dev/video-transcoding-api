package service

import (
	"errors"
	"net/http"
	"testing"
)

func TestErrorResponse(t *testing.T) {
	internalError := errors.New("something went wrong")
	errResp := newErrorResponse(internalError)
	code, data, err := errResp.Result()
	if code != http.StatusInternalServerError {
		t.Errorf("Wrong error code. Want %d. Got %d", http.StatusInternalServerError, code)
	}
	if data != nil {
		t.Errorf("Unexpected non-nil data: %#v", data)
	}
	if err != errResp {
		t.Errorf("Wrong error returned. Want %#v. Got %#v", errResp, err)
	}
}

func TestErrorResponseCustomStatus(t *testing.T) {
	internalError := errors.New("something went wrong")
	errResp := newErrorResponse(internalError).withStatus(http.StatusBadRequest)
	code, data, err := errResp.Result()
	if code != http.StatusBadRequest {
		t.Errorf("Wrong error code. Want %d. Got %d", http.StatusBadRequest, code)
	}
	if data != nil {
		t.Errorf("Unexpected non-nil data: %#v", data)
	}
	if err != errResp {
		t.Errorf("Wrong error returned. Want %#v. Got %#v", errResp, err)
	}
}

func TestErrorResponsErrorInterface(t *testing.T) {
	var err error
	msg := "something went wrong"
	err = &errorResponse{Message: msg}
	if err.Error() != msg {
		t.Errorf("Got wrong error message. Want %q. Got %q", msg, err.Error())
	}
}
