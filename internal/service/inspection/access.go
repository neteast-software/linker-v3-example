package inspection

import (
	"strconv"

	"github.com/neteast-software/go-module/acl"
	"github.com/neteast-software/go-module/application"
)

type TaskAccess struct {
	Actor    acl.Access
	Resource acl.Resource
	Range    RecordRange
}

func NewTaskAccess(app application.Application, actorID uint64) TaskAccess {
	id := strconv.FormatUint(actorID, 10)
	return TaskAccess{
		Actor: acl.Access{
			ID:      id,
			Level:   1,
			Actions: acl.Read,
			Groups: []acl.ResourceGroup{
				acl.Role("inspection-reader", acl.NewResource(taskResourcePattern(app), acl.Scope(app.Scope, 0, taskResourceDescription, acl.Read))),
			},
		},
		Resource: TaskResource(app),
		Range:    OwnerRange(actorID),
	}
}

func (p TaskAccess) Can(action acl.Action) bool {
	return p.Actor.Can(p.Resource, action)
}
