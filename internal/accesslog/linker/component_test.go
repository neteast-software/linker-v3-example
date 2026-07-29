package accesslog

import (
	"context"
	"testing"
	"time"

	linker "github.com/neteast-software/linker/v3"

	core "linker-v3-example/internal/accesslog"
)

func TestComponentOwnsWorkerLifecycleAndAsset(t *testing.T) {
	log, err := core.New(core.Memory(), 8)
	if err != nil {
		t.Fatal(err)
	}
	app := linker.New(linker.WithComponents(New(log)))
	if err = app.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	if stats := log.Stats(); stats.Capacity != 8 {
		t.Fatalf("stats = %#v", stats)
	}
	found := false
	for _, asset := range app.Plan().Assets {
		if asset.Kind == AssetQueue {
			found = true
		}
	}
	if !found {
		t.Fatal("accesslog queue asset missing")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err = app.Stop(ctx); err != nil {
		t.Fatal(err)
	}
}
