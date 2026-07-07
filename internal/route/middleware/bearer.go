package middleware

import (
	"strings"

	http "github.com/neteast-software/go-module/http/gin/linker"

	userconstant "linker-v3-example/internal/constant/user"
)

func Bearer(c *http.Context) (string, error) {
	token, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
	if !ok || token == "" {
		return "", userconstant.ErrLogin
	}
	return token, nil
}
