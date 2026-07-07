package example

import (
	"context"
	"errors"
	"testing"
	"time"

	server "github.com/neteast-software/go-module/linker/server"
	linker "github.com/neteast-software/linker/v3"

	inspectioncomponent "linker-v3-example/internal/component/inspection"
)

type slowStopComponent struct{}

func (slowStopComponent) Identity() linker.ID {
	return "example/slow-stop"
}

func (slowStopComponent) Start(context.Context) error {
	return nil
}

func (slowStopComponent) Stop(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func TestDBBackedComponentReportsMissingStoreDependency(t *testing.T) {
	component := inspectioncomponent.NewComponent()
	runtime := linker.New(linker.WithMode(linker.Bin))

	err := component.Init(context.Background(), runtime)
	if err == nil {
		t.Fatal("expected missing db capability error")
	}
}

func TestServerStopTimeoutExample(t *testing.T) {
	app := server.New(
		server.WithMode(linker.Server),
		server.WithShutdownTimeout(time.Nanosecond),
		server.WithoutHTTP(),
		server.WithoutEvent(),
		server.WithoutNotice(),
		server.WithoutAudit(),
		server.WithComponents(slowStopComponent{}),
	)
	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}

	err := app.Stop(context.Background())
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}
}
