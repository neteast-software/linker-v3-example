package gateway

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInfo(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/gateway/info", nil)
	response := httptest.NewRecorder()
	serveInfo(response, request)
	if response.Code != http.StatusOK || response.Body.String() != "{\"service\":\"linker-v3-example\"}\n" {
		t.Fatalf("status=%d body=%q", response.Code, response.Body.String())
	}
	if contentType := response.Header().Get("Content-Type"); contentType != "application/json; charset=utf-8" {
		t.Fatalf("content-type=%q", contentType)
	}
}
