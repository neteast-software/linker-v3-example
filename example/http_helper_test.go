package example

import (
	"io"
	stdhttp "net/http"
	"testing"
)

func getRaw(t *testing.T, url string) string {
	t.Helper()
	response, err := stdhttp.Get(url)
	if err != nil {
		t.Fatalf("get raw: %v", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read raw: %v", err)
	}
	return string(body)
}
