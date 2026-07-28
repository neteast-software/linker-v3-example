package inspection

import (
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"

	access "linker-v3-example/internal/access"
	user "linker-v3-example/internal/user"
)

func currentActorID(c *http.Context) (uint64, bool) {
	raw, err := access.Bearer(c)
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return 0, false
	}
	actor, err := http.Require(c, user.ActorKey())
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return 0, false
	}
	id, err := actor.ActorID(c.Request.Context(), raw, "front")
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return 0, false
	}
	return id, true
}
