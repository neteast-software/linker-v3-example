package example

import (
	"context"
	"sync"
	"testing"
	"time"

	server "github.com/neteast-software/go-module/linker/server"
	feishu "github.com/neteast-software/go-module/notify/feishu"
	feishucomponent "github.com/neteast-software/go-module/notify/feishu/linker"
	linker "github.com/neteast-software/linker/v3"
)

func TestFeishuFrameworkNoticeExample(t *testing.T) {
	sender := newRecordingFeishuSender()
	provider := feishucomponent.New(
		feishucomponent.WithSender(sender),
		feishucomponent.WithTargets(feishu.Target(feishu.ChatID, "oc_example")),
	)
	reporter := &noticeReporter{}
	app := server.New(
		server.WithoutStartupLog(),
		server.WithoutHTTP(),
		server.WithoutAudit(),
		server.WithoutEvent(),
		server.WithShutdownTimeout(time.Second),
		server.WithComponents(provider, reporter),
	)
	if !planHasComponent(app.Plan(), feishucomponent.ID) {
		t.Fatalf("plan missing feishu component: %#v", app.Plan().Components)
	}

	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		if err := app.Stop(context.Background()); err != nil {
			t.Fatalf("close: %v", err)
		}
	})

	if err := linker.Report(context.Background(), reporter.runtime, linker.Detected("provider_unavailable", nil)); err != nil {
		t.Fatalf("report detected: %v", err)
	}
	waitFeishuMessage(t, sender.sent)
	if err := linker.Report(context.Background(), reporter.runtime, linker.Recovered("provider_unavailable")); err != nil {
		t.Fatalf("report recovered: %v", err)
	}
	waitFeishuMessage(t, sender.sent)

	messages, targets := sender.snapshot()
	if len(messages) != 2 || len(targets) != 2 || targets[0].ReceiveID != "oc_example" {
		t.Fatalf("sender messages=%#v targets=%#v", messages, targets)
	}
}

type noticeReporter struct {
	runtime linker.Runtime
}

func (p *noticeReporter) Identity() linker.ID {
	return "example/notice-reporter"
}

func (p *noticeReporter) Init(_ context.Context, runtime linker.Runtime) error {
	p.runtime = runtime
	return nil
}

type recordingFeishuSender struct {
	mu       sync.Mutex
	messages []feishu.Message
	targets  []feishu.Receiver
	sent     chan struct{}
}

func newRecordingFeishuSender() *recordingFeishuSender {
	return &recordingFeishuSender{sent: make(chan struct{}, 2)}
}

func (p *recordingFeishuSender) SendMessage(_ context.Context, message feishu.Message, targets ...feishu.Receiver) (feishu.Result, error) {
	p.mu.Lock()
	p.messages = append(p.messages, message)
	p.targets = append(p.targets, targets...)
	p.mu.Unlock()
	p.sent <- struct{}{}
	return feishu.Result{Status: feishu.StatusSent}, nil
}

func (p *recordingFeishuSender) snapshot() ([]feishu.Message, []feishu.Receiver) {
	p.mu.Lock()
	defer p.mu.Unlock()
	messages := append([]feishu.Message(nil), p.messages...)
	targets := append([]feishu.Receiver(nil), p.targets...)
	return messages, targets
}

func waitFeishuMessage(t *testing.T, sent <-chan struct{}) {
	t.Helper()
	select {
	case <-sent:
	case <-time.After(3 * time.Second):
		t.Fatal("feishu notice timeout")
	}
}
