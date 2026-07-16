package user

import (
	"strings"

	http "github.com/neteast-software/go-module/http/gin/linker"

	userservice "linker-v3-example/internal/service/user"
)

func bearerToken(c *http.Context) (string, error) {
	token, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
	if !ok || token == "" {
		return "", userservice.ErrLogin
	}
	return token, nil
}
