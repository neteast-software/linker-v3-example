package console

import (
	"testing"

	"github.com/neteast-software/go-module/acl"

	user "linker-v3-example/internal/user"
)

func TestAccessCompilesRoleBoundary(t *testing.T) {
	reader := Access(user.User{Role: "user"}, "2")
	if !reader.Can(
		acl.NewResource("graph.console.menu", acl.Scope("console", 1, "后台菜单", acl.Read)),
		acl.Read,
	) {
		t.Fatal("普通用户应能读取 Console 公共业务资源")
	}
	if reader.Can(
		acl.NewResource("http.app2.order.update", acl.Scope("app2", 2, "订单维护", acl.Update)),
		acl.Update,
	) {
		t.Fatal("普通用户不应获得订单维护权限")
	}

	admin := Access(user.User{Role: "admin"}, "1")
	if !admin.Can(
		acl.NewResource("http.app2.order.update", acl.Scope("app2", 2, "订单维护", acl.Update)),
		acl.Update,
	) {
		t.Fatal("管理员应覆盖后台业务资源")
	}
}
