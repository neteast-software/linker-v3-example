package access

import (
	"strings"

	http "github.com/neteast-software/go-module/http/gin/linker"
)

func Bearer(c *http.Context) (string, error) {
	token, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
	if !ok || token == "" {
		return "", ErrBearer
	}
	return token, nil
}
