package example

import (
	"context"
	"testing"

	linker "github.com/neteast-software/linker/v3"
)

func preparedPlan(t *testing.T, app *linker.App) linker.Plan {
	t.Helper()
	if err := app.Prepare(context.Background()); err != nil {
		t.Fatalf("prepare plan: %v", err)
	}
	return app.Plan()
}
