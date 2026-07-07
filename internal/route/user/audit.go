package user

import (
	"strconv"

	audithttp "github.com/neteast-software/go-module/audit/operate/http/gin"
	http "github.com/neteast-software/go-module/http/gin/linker"

	usermodel "linker-v3-example/internal/model/user"
)

func setUserOperator(c *http.Context, user usermodel.User) {
	if user.ID == 0 {
		return
	}
	audithttp.SetOperator(c, strconv.FormatUint(user.ID, 10), user.Username)
}
