package inspection

import (
	"strconv"

	"github.com/neteast-software/go-module/acl"
	"github.com/neteast-software/go-module/db/gorm/query"
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"

	inspectionmodel "linker-v3-example/internal/model/inspection"
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
	rows, req, err := svc.List(c.Request.Context(), app, inspectionservice.ListRequest{
		Page:   query.NewPage(intQuery(c, "page"), intQuery(c, "pageSize")),
		Status: c.Query("status"),
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

type taskItem struct {
	ID               uint64 `json:"id"`
	ApplicationScope string `json:"application_scope"`
	Title            string `json:"title"`
	Status           string `json:"status"`
	OwnerID          uint64 `json:"owner_id"`
}

func taskItems(rows []inspectionmodel.Task) []taskItem {
	ret := make([]taskItem, 0, len(rows))
	for _, row := range rows {
		ret = append(ret, taskItem{
			ID:               row.ID,
			ApplicationScope: row.ApplicationScope,
			Title:            row.Title,
			Status:           row.Status,
			OwnerID:          row.OwnerID,
		})
	}
	return ret
}

func intQuery(c *http.Context, key string) int {
	value := c.Query(key)
	if value == "" {
		return 0
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return n
}
