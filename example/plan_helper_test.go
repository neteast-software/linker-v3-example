package example

import (
	"context"
	"testing"

	server "github.com/neteast-software/go-module/linker/server"
)

func preparedPlan(t *testing.T, app *server.App) server.Plan {
	t.Helper()
	if err := app.Prepare(context.Background()); err != nil {
		t.Fatalf("prepare plan: %v", err)
	}
	return app.Plan()
}
