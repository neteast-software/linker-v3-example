package inspection

import (
	"github.com/neteast-software/go-module/acl"
	"github.com/neteast-software/go-module/db/gorm/query"
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"

	access "linker-v3-example/internal/access"
	inspection "linker-v3-example/internal/inspection"
)

func init() {
	http.RegisterIn("api/v1/app2",
		http.Group("inspection",
			http.Use(access.Application("app2")),
			http.GET("tasks", listTasks).Resource(
				"http.app2.inspection.tasks",
				acl.Scope("app2", 1, "应用二巡检任务"),
			),
		),
	)
}

func listTasks(c *http.Context) {
	app, ok := access.CurrentApplication(c)
	if !ok {
		response.Warning(c, "应用入口未识别")
		return
	}
	svc, err := http.Require(c, inspection.ServiceKey())
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return
	}
	actorID, ok := currentActorID(c)
	if !ok {
		return
	}
	var status inspection.Status
	if raw := c.Query("status"); raw != "" {
		status, err = inspection.ParseStatus(raw)
		if err != nil {
			response.Warning(c, "%s", err.Error())
			return
		}
	}
	rows, req, err := svc.List(c.Request.Context(), app, inspection.ListRequest{
		Page:   query.NewPage(intQuery(c, "page"), intQuery(c, "pageSize")),
		Status: status,
		Access: inspection.NewTaskAccess(app, actorID),
	})
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return
	}
	response.Data(c, map[string]any{
		"items":     taskItems(rows),
		"page":      req.Page.Number,
		"page_size": req.Page.Size,
	})
}
