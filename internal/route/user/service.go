package user

import (
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"

	userservice "linker-v3-example/internal/service/user"
)

func service(c *http.Context) (userservice.Service, bool) {
	svc, err := http.Require(c, userservice.ServiceKey())
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return userservice.Service{}, false
	}
	return svc, true
}
