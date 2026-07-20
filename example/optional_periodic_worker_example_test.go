package example

import (
	"context"
	"errors"
	"testing"
	"time"

	server "github.com/neteast-software/go-module/linker/server"
	periodic "github.com/neteast-software/go-module/worker/periodic"
	managed "github.com/neteast-software/go-module/worker/periodic/linker"
)

func TestOptionalPeriodicWorkerExample(t *testing.T) {
	failed := make(chan struct{})
	projection := periodic.New(
		"graph-projection",
		periodic.HandlerFunc(func(context.Context) error {
			close(failed)
			return errors.New("图存储暂时不可用")
		}),
		periodic.WithDescription("异步投影关系图"),
		periodic.Every(time.Hour),
		periodic.WithBackoff(periodic.Exponential(time.Second, time.Minute)),
		periodic.UnhealthyAfter(0),
	)
	app := server.New(
		server.WithoutStartupLog(),
		server.WithoutHTTP(),
		server.WithoutEvent(),
		server.WithoutNotice(),
		server.WithoutAudit(),
		server.WithComponents(managed.New(managed.WithWorkers(projection))),
	)
	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("启动可选周期 Worker：%v", err)
	}
	t.Cleanup(func() { _ = app.Stop(context.Background()) })

	select {
	case <-failed:
	case <-time.After(time.Second):
		t.Fatal("可选周期 Worker 未执行")
	}
	waitForWorkerFailure(t, projection)
	if err := app.Health(context.Background()); err != nil {
		t.Fatalf("可选周期 Worker 不应影响主服务健康：%v", err)
	}
	if !planHasAsset(app.Plan(), "worker/periodic", "graph-projection") {
		t.Fatalf("Plan 缺少可选周期 Worker Asset：%#v", app.Plan().Assets)
	}
}

func waitForWorkerFailure(t *testing.T, worker *periodic.Worker) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for worker.Stats().Failed == 0 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if worker.Stats().Failed == 0 {
		t.Fatal("周期 Worker 未记录失败")
	}
}
