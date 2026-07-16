package inspection

import (
	"github.com/neteast-software/go-module/acl"
	"github.com/neteast-software/go-module/db/gorm/query"
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"

	inspectionconstant "linker-v3-example/internal/constant/inspection"
	routemiddleware "linker-v3-example/internal/route/middleware"
	inspectionservice "linker-v3-example/internal/service/inspection"
)

func init() {
	http.RegisterIn("api/v1/app2",
		http.Group("inspection",
			http.Use(routemiddleware.Application("app2")),
			http.GET("tasks", listTasks).Resource(
				"http.app2.inspection.tasks",
				acl.Scope("app2", 1, "应用二巡检任务"),
			),
		),
	)
}

func listTasks(c *http.Context) {
	app, ok := routemiddleware.CurrentApplication(c)
	if !ok {
		response.Warning(c, "应用入口未识别")
		return
	}
	svc, err := http.Require(c, inspectionservice.ServiceKey())
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return
	}
	actorID, ok := currentActorID(c)
	if !ok {
		return
	}
	var status inspectionconstant.Status
	if raw := c.Query("status"); raw != "" {
		status, err = inspectionconstant.ParseStatus(raw)
		if err != nil {
			response.Warning(c, "%s", err.Error())
			return
		}
	}
	rows, req, err := svc.List(c.Request.Context(), app, inspectionservice.ListRequest{
		Page:   query.NewPage(intQuery(c, "page"), intQuery(c, "pageSize")),
		Status: status,
		Access: inspectionservice.NewTaskAccess(app, actorID),
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
