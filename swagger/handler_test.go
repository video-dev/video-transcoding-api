package swagger

import (
	"errors"
	"net/http"
	"testing"
)

func TestHandlerToJSONEndpoint(t *testing.T) {
	errResp := NewErrorResponse(errors.New("something went wrong")).WithStatus(http.StatusGone)
	handler := Handler(func(r *http.Request) GizmoJSONResponse {
		return errResp
	})
	req, _ := http.NewRequest("GET", "/something", nil)
	code, data, err := HandlerToJSONEndpoint(handler)(req)
	if code != http.StatusGone {
		t.Errorf("wrong status code returned, want %d, got %d", http.StatusGone, code)
	}
	if data != nil {
		t.Errorf("unexpected non-nil data: %s", data)
	}
	if err != errResp {
		t.Errorf("wrong error returned\nwant %#v\ngot  %#v", errResp, err)
	}
}
