package inspection

import (
	"fmt"

	"github.com/neteast-software/go-module/acl"
	"github.com/neteast-software/go-module/application"
)

const taskResourceDescription = "巡检任务记录"

func TaskResource(app application.Application) acl.Resource {
	return acl.NewResource(taskResourceTag(app), acl.Scope(app.Scope, 1, taskResourceDescription, acl.Read))
}

func taskResourceTag(app application.Application) string {
	return fmt.Sprintf("db.%s.inspection.task", app.Scope)
}

func taskResourcePattern(app application.Application) string {
	return fmt.Sprintf("db.%s.inspection.*", app.Scope)
}
