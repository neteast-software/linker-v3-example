package notification

import (
	"github.com/neteast-software/go-module/acl"
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"
	consumer "github.com/neteast-software/go-module/mq/consumer"
	mq "github.com/neteast-software/go-module/mq/consumer/linker"
	traceconsumer "github.com/neteast-software/go-module/observe/tracing/mq/consumer"
)

type sendRequest struct {
	Key  string `json:"key"`
	Body string `json:"body"`
}

func init() {
	http.RegisterIn("api/v1/app2/notification",
		http.POST("send", sendAPI).Resource(
			"http.app2.notification.send",
			acl.Scope("app2", 1, "通知发送", acl.Write),
		),
	)
}

func sendAPI(c *http.Context) {
	req := sendRequest{Key: "http", Body: "hello notification"}
	if c.Request.Body != nil && c.Request.ContentLength != 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Warning(c, "通知发送参数错误: %s", err.Error())
			return
		}
	}
	executor, err := http.Require(c, mq.ConsumerKey("notification"))
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return
	}
	message := traceconsumer.InjectMessage(c.Request.Context(), consumer.Message{
		Topic: "notification.message",
		Key:   req.Key,
		Body:  []byte(req.Body),
	})
	if err = executor.Submit(c.Request.Context(), message); err != nil {
		response.Warning(c, "%s", err.Error())
		return
	}
	response.Data(c, map[string]any{"status": "queued"})
}
