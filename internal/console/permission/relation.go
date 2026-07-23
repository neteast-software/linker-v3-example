package permission

import (
	"github.com/neteast-software/go-module/graph/console/multilist"
	"github.com/neteast-software/go-module/graph/console/protocol"
	"github.com/neteast-software/go-module/graph/console/viewer"

	order "linker-v3-example/internal/order"
	permissionresource "linker-v3-example/internal/permission"
)

func Relation() *multilist.Multilist[role, resource] {
	roles := viewer.New[role]("角色").
		Identity("permission.role").
		Columns(viewer.Col("id", "ID"), viewer.Col("name", "角色名称")).
		Choose().
		WithData(
			role{ID: 1, Name: "系统管理员"},
			role{ID: 2, Name: "只读用户"},
		)
	resources := viewer.New[resource]("权限资源").
		Identity("permission.resource").
		Key("tag").
		Columns(viewer.Col("tag", "资源"), viewer.Col("description", "说明")).
		Select().
		WithData(
			resource{Tag: order.List, Description: "查看订单列表"},
			resource{Tag: order.Update, Description: "维护订单"},
		)
	target := protocol.Route("/api/v1/app2/permission/role/:id/resource")
	relation := multilist.Relate(
		"id",
		"resources",
		multilist.Load(target, permissionresource.Read),
		multilist.Assign(target, permissionresource.Assign),
		multilist.Remove(target, permissionresource.Remove),
		multilist.WithCurrent(1, order.List, order.Update),
	)
	return multilist.New("角色权限", roles, resources, relation).
		Identity("permission.role-resource").
		Resource(permissionresource.Manage)
}
