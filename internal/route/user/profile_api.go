package user

import (
	"github.com/neteast-software/go-module/acl"
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"
)

func init() {
	http.RegisterIn("api",
		http.GET("profile", profileAPI).Resource(
			"http.front.user.profile",
			acl.Scope("front", 1, "前台用户信息"),
		),
	)
}

func profileAPI(c *http.Context) {
	token, err := bearerToken(c)
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return
	}
	svc, ok := service(c)
	if !ok {
		return
	}
	currentUser, err := svc.Profile(c.Request.Context(), token, "front")
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return
	}
	setUserOperator(c, currentUser)
	response.Data(c, newProfile(currentUser))
}
