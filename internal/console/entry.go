package console

import (
	configprotocol "github.com/neteast-software/go-module/graph/console/config"
	"github.com/neteast-software/go-module/graph/console/field"
	"github.com/neteast-software/go-module/graph/console/login"
	"github.com/neteast-software/go-module/graph/console/protocol"
)

func Entry() configprotocol.Entry {
	value := configprotocol.New(
		"Linker v3 Example",
		protocol.Graph("/console/page/dashboard"),
		configprotocol.Pages(configprotocol.Multiple),
		configprotocol.IframeOrigins("https://docs.neteast.cn"),
		configprotocol.Themes("default", Theme()),
	)
	value.Login = []login.Method{
		login.PasswordMethod(
			field.Text("username", "用户名").Required("请输入用户名"),
			field.Password("password", "密码").Required("请输入密码"),
		),
	}
	value.Features["notification"] = true
	value.Features["theme"] = true
	return value
}
