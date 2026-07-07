package user

import (
	"github.com/neteast-software/go-module/acl"
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"
)

func init() {
	http.RegisterIn("system",
		http.GET("profile", systemProfile).Resource(
			"http.console.user.profile",
			acl.Scope("console", 1, "后台用户信息"),
		),
	)
}

func systemProfile(c *http.Context) {
	token, err := bearerToken(c)
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return
	}
	svc, ok := service(c)
	if !ok {
		return
	}
	currentUser, err := svc.Profile(c.Request.Context(), token, "console")
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return
	}
	setUserOperator(c, currentUser)
	response.Data(c, newProfile(currentUser))
}
