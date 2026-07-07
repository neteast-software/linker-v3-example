package inspection

import (
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"

	routemiddleware "linker-v3-example/internal/route/middleware"
	userservice "linker-v3-example/internal/service/user"
)

func currentActorID(c *http.Context) (uint64, bool) {
	raw, err := routemiddleware.Bearer(c)
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return 0, false
	}
	svc, err := http.Require(c, userservice.ServiceKey())
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return 0, false
	}
	user, err := svc.Profile(c.Request.Context(), raw, "front")
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return 0, false
	}
	return user.ID, true
}
