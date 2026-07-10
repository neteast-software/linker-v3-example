package example

import (
	"context"
	"testing"

	http "github.com/neteast-software/go-module/http/gin/linker"
	linker "github.com/neteast-software/linker/v3"
)

type initProbe struct {
	initialized bool
}

func (p *initProbe) Identity() linker.ID {
	return "example/init-probe"
}

func (p *initProbe) Init(context.Context, linker.Runtime) error {
	p.initialized = true
	return nil
}

func TestCoreBinDoesNotLoadServerDefaults(t *testing.T) {
	probe := &initProbe{}
	app := linker.New(
		linker.WithMode(linker.Bin),
		linker.WithComponents(probe),
	)

	plan := app.Plan()
	for _, component := range plan.Components {
		if component.ID == http.ID {
			t.Fatalf("http component should be disabled: %#v", plan.Components)
		}
	}

	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	if !probe.initialized {
		t.Fatal("bin component was not initialized")
	}
	if _, ok := linker.Resolve(app, linker.NewCapabilityKey[*http.Server](http.ID)); ok {
		t.Fatalf("http capability should be absent")
	}
	if err := app.Stop(context.Background()); err != nil {
		t.Fatalf("stop: %v", err)
	}
}
