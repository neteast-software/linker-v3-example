package inspection

import (
	"strconv"

	http "github.com/neteast-software/go-module/http/gin/linker"
)

func intQuery(c *http.Context, key string) int {
	value := c.Query(key)
	if value == "" {
		return 0
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return n
}
