package gateway

import (
	"encoding/json"
	"net/http"

	gatewaycomponent "github.com/neteast-software/go-module/http/gateway/linker"
)

// Endpoints 返回由当前业务源码拥有、在 Gateway 启动前冻结的本地端点。
func Endpoints() []gatewaycomponent.LocalEndpoint {
	return []gatewaycomponent.LocalEndpoint{
		gatewaycomponent.Handle("gateway/info", "GET /gateway/info", http.HandlerFunc(serveInfo)),
	}
}

func serveInfo(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(writer).Encode(map[string]string{"service": "linker-v3-example"})
}
