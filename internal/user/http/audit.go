package user

import (
	"strconv"

	"github.com/neteast-software/go-module/audit/operate/http/gin"
	http "github.com/neteast-software/go-module/http/gin/linker"

	user "linker-v3-example/internal/user"
)

func setUserOperator(c *http.Context, current user.User) {
	if current.ID == 0 {
		return
	}
	audit.SetOperator(c, strconv.FormatUint(current.ID, 10), current.Username)
}
