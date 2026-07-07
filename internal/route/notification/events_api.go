package notification

import (
	"net/http"

	"github.com/neteast-software/go-module/acl"
	ginhttp "github.com/neteast-software/go-module/http/gin/linker"
)

func init() {
	ginhttp.RegisterIn("api/v1/app2/notification",
		ginhttp.GET("events", eventsAPI).
			With(ginhttp.SSEHeader).
			Resource("http.app2.notification.events", acl.Scope("app2", 1, "通知事件流", acl.Read)),
	)
}

func eventsAPI(c *ginhttp.Context) {
	c.SSEvent("ready", map[string]string{"status": "connected"})
	c.Status(http.StatusOK)
	if flusher, ok := c.Writer.(http.Flusher); ok {
		flusher.Flush()
	}
}
