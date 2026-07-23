package notification

import (
	stdhttp "net/http"

	"github.com/neteast-software/go-module/acl"
	http "github.com/neteast-software/go-module/http/gin/linker"
)

func init() {
	http.RegisterIn("api/v1/app2/notification",
		http.GET("events", events).
			With(http.SSEHeader).
			Resource("http.app2.notification.events", acl.Scope("app2", 1, "通知事件流", acl.Read)),
	)
}

func events(c *http.Context) {
	c.SSEvent("ready", map[string]string{"status": "connected"})
	c.Status(stdhttp.StatusOK)
	if flusher, ok := c.Writer.(stdhttp.Flusher); ok {
		flusher.Flush()
	}
}
