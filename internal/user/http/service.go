package user

import (
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"

	user "linker-v3-example/internal/user"
)

func service(c *http.Context) (*user.Service, bool) {
	svc, err := http.Require(c, user.ServiceKey())
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return nil, false
	}
	return svc, true
}
