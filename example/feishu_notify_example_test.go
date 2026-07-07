package example

import (
	"context"
	"testing"
	"time"

	server "github.com/neteast-software/go-module/linker/server"
	feishu "github.com/neteast-software/go-module/notify/feishu"
	feishucomponent "github.com/neteast-software/go-module/notify/feishu/linker"
	linker "github.com/neteast-software/linker/v3"
)

func TestFeishuNotifyCapabilityExample(t *testing.T) {
	sender := &recordingFeishuSender{}
	app := server.New(
		server.WithMode(linker.Bin),
		server.WithShutdownTimeout(time.Second),
		server.WithComponents(
			feishucomponent.New(feishucomponent.WithSender(sender)),
		),
	)
	if !planHasComponent(app.Plan(), feishucomponent.ID) {
		t.Fatalf("plan missing feishu component: %#v", app.Plan().Components)
	}

	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		if err := app.Close(context.Background()); err != nil {
			t.Fatalf("close: %v", err)
		}
	})

	got, err := linker.RequireCapability(app.App(), feishucomponent.SenderKey())
	if err != nil {
		t.Fatalf("sender capability: %v", err)
	}
	result, err := got.SendMessage(
		context.Background(),
		feishu.Text("服务通知", "linker v3 example"),
		feishu.Target(feishu.ChatID, "oc_example"),
	)
	if err != nil {
		t.Fatalf("send: %v", err)
	}
	if result.Status != feishu.StatusSent {
		t.Fatalf("result = %#v", result)
	}
	if len(sender.messages) != 1 || len(sender.targets) != 1 || sender.targets[0].ReceiveID != "oc_example" {
		t.Fatalf("sender messages=%#v targets=%#v", sender.messages, sender.targets)
	}
}

type recordingFeishuSender struct {
	messages []feishu.Message
	targets  []feishu.Receiver
}

func (p *recordingFeishuSender) SendMessage(_ context.Context, message feishu.Message, targets ...feishu.Receiver) (feishu.Result, error) {
	p.messages = append(p.messages, message)
	p.targets = append(p.targets, targets...)
	return feishu.Result{Status: feishu.StatusSent}, nil
}
