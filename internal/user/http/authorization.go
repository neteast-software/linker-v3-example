package user

import (
	"strings"

	http "github.com/neteast-software/go-module/http/gin/linker"

	user "linker-v3-example/internal/user"
)

func bearerToken(c *http.Context) (string, error) {
	token, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
	if !ok || token == "" {
		return "", user.ErrLogin
	}
	return token, nil
}
