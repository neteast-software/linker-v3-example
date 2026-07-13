package tts

import (
	"github.com/neteast-software/go-module/acl"
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"

	ttsclient "linker-v3-example/internal/client/tts"
)

type transcribeRequest struct {
	Text string `json:"text"`
}

func init() {
	http.RegisterIn("api/v1/app2/tts",
		http.POST("transcribe", transcribe).Resource(
			"http.app2.tts.transcribe",
			acl.Scope("app2", 1, "TTS 转写", acl.Write),
		),
	)
}

func transcribe(c *http.Context) {
	var req transcribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Warning(c, "TTS 参数错误: %s", err.Error())
		return
	}
	client, err := http.Require(c, ttsclient.ClientKey())
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return
	}
	result, err := client.Transcribe(c.Request.Context(), req.Text)
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return
	}
	response.Data(c, map[string]any{"result": result})
}
