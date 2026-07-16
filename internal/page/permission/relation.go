package permission

import (
	"github.com/neteast-software/go-module/graph/console/multilist"
	"github.com/neteast-software/go-module/graph/console/protocol"
	"github.com/neteast-software/go-module/graph/console/viewer"

	orderconstant "linker-v3-example/internal/constant/order"
	permissionconstant "linker-v3-example/internal/constant/permission"
)

func Relation() *multilist.Multilist[map[string]any, map[string]any] {
	roles := viewer.New[map[string]any]("角色").
		Identity("permission.role").
		Columns(viewer.Col("id", "ID"), viewer.Col("name", "角色名称")).
		Choose().
		WithData(
			map[string]any{"id": 1, "name": "系统管理员"},
			map[string]any{"id": 2, "name": "只读用户"},
		)
	resources := viewer.New[map[string]any]("权限资源").
		Identity("permission.resource").
		Key("tag").
		Columns(viewer.Col("tag", "资源"), viewer.Col("description", "说明")).
		Select().
		WithData(
			map[string]any{"tag": orderconstant.List, "description": "查看订单列表"},
			map[string]any{"tag": orderconstant.Update, "description": "维护订单"},
		)
	target := protocol.Route("/api/v1/app2/permission/role/:id/resource")
	relation := multilist.Relate(
		"id",
		"resources",
		multilist.Load(target, permissionconstant.Read),
		multilist.Assign(target, permissionconstant.Assign),
		multilist.Remove(target, permissionconstant.Remove),
		multilist.WithCurrent(1, orderconstant.List, orderconstant.Update),
	)
	return multilist.New("角色权限", roles, resources, relation).
		Identity("permission.role-resource").
		Resource(permissionconstant.Manage)
}
