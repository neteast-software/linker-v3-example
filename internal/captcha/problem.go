package captcha

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func problem(status int, code, detail string) *http.Response {
	content, _ := json.Marshal(map[string]any{
		"type":   "https://neteast.cn/problems/" + code,
		"title":  http.StatusText(status),
		"status": status,
		"code":   code,
		"detail": detail,
	})
	return &http.Response{
		StatusCode: status,
		Header: http.Header{
			"Content-Type":  {"application/problem+json; charset=utf-8"},
			"Cache-Control": {"no-store"},
		},
		Body:          io.NopCloser(bytes.NewReader(content)),
		ContentLength: int64(len(content)),
	}
}
