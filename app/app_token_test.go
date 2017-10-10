package app_test

import (
	"net/http"
	"testing"
)

func TestReturnClientID(t *testing.T) {

	req, _ := http.NewRequest("GET", "/pushes", nil)
	req.Header.Add("authorization", testAuthorization)
	_ = executeRequest(req)

	if len(a.ReturnClientID(req)) == 0 {
		t.Error("Failed to obtain the clientID")
	}
}
