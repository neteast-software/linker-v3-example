package user

import (
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"
)

type loginRequest struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func bindLoginRequest(c *http.Context) (loginRequest, bool) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Warning(c, "登录参数错误: %s", err.Error())
		return loginRequest{}, false
	}
	return req, true
}
