package observability

import (
	"github.com/neteast-software/go-module/acl"
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"
)

func init() {
	http.Register(
		http.GET("metrics", metricsAPI).Resource(
			"http.observe.metrics",
			acl.Scope("observe", 0, "Prometheus 指标", acl.Read),
		),
	)
}

func metricsAPI(c *http.Context) {
	recorder, ok := prometheusRecorder(c)
	if !ok {
		response.Warning(c, "Prometheus 指标能力未启用")
		return
	}
	recorder.Handler().ServeHTTP(c.Writer, c.Request)
}
