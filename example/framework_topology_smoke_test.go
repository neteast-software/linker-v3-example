package example

import (
	"context"
	"errors"
	"slices"
	"testing"

	linker "github.com/neteast-software/linker/v3"
)

var smokeServiceKey = linker.NewCapabilityKey[string]("example/smoke-service")

type smokeManager struct {
	events *[]string
}

func (p *smokeManager) Identity() linker.ID {
	return "example/smoke-manager"
}

func (p *smokeManager) AssetPolicies() linker.AssetPolicies {
	return linker.AssetPolicies{linker.SourceFirst(linker.AssetRoute)}
}

func (p *smokeManager) Init(context.Context, linker.Runtime) error {
	*p.events = append(*p.events, "manager:init")
	return nil
}

func (p *smokeManager) Stop(context.Context) error {
	*p.events = append(*p.events, "manager:stop")
	return nil
}

type smokeProvider struct {
	events *[]string
}

func (p *smokeProvider) Identity() linker.ID {
	return "example/smoke-provider"
}

func (p *smokeProvider) Capabilities() linker.Capabilities {
	return linker.Capabilities{
		linker.Offer(smokeServiceKey, func() string {
			return "ready"
		}),
	}
}

func (p *smokeProvider) Assets(context.Context, linker.Runtime) ([]linker.Asset, error) {
	return []linker.Asset{
		linker.RouteAsset("example/smoke-manager", "GET /smoke"),
	}, nil
}

func (p *smokeProvider) Init(context.Context, linker.Runtime) error {
	*p.events = append(*p.events, "provider:init")
	return nil
}

func (p *smokeProvider) Start(context.Context) error {
	*p.events = append(*p.events, "provider:start")
	return nil
}

func (p *smokeProvider) OnMounted(context.Context, linker.Runtime) error {
	*p.events = append(*p.events, "provider:mounted")
	return nil
}

func (p *smokeProvider) Stop(context.Context) error {
	*p.events = append(*p.events, "provider:stop")
	return nil
}

func (p *smokeProvider) Close(context.Context) error {
	*p.events = append(*p.events, "provider:close")
	return nil
}

type smokeConsumer struct {
	events   *[]string
	failInit bool
}

func (p *smokeConsumer) Identity() linker.ID {
	return "example/smoke-consumer"
}

func (p *smokeConsumer) Capabilities() linker.Capabilities {
	return linker.Capabilities{linker.Need(smokeServiceKey)}
}

func (p *smokeConsumer) Init(context.Context, linker.Runtime) error {
	*p.events = append(*p.events, "consumer:init")
	if p.failInit {
		return errors.New("模拟 consumer 初始化失败")
	}
	return nil
}

func TestFrameworkTopologySurvivesComponentReordering(t *testing.T) {
	permutations := [][]string{
		{"manager", "consumer", "provider"},
		{"consumer", "provider", "manager"},
		{"provider", "manager", "consumer"},
	}
	for _, order := range permutations {
		t.Run(order[0]+"-"+order[1]+"-"+order[2], func(t *testing.T) {
			var events []string
			components := map[string]linker.Component{
				"manager":  &smokeManager{events: &events},
				"provider": &smokeProvider{events: &events},
				"consumer": &smokeConsumer{events: &events},
			}
			ordered := make([]linker.Component, 0, len(order))
			for _, name := range order {
				ordered = append(ordered, components[name])
			}
			app := linker.New(linker.WithMode(linker.Server), linker.WithComponents(ordered...))
			if err := app.Prepare(context.Background()); err != nil {
				t.Fatalf("prepare: %v", err)
			}
			ids := componentIDs(app.Plan())
			provider := slices.Index(ids, linker.ID("example/smoke-provider"))
			manager := slices.Index(ids, linker.ID("example/smoke-manager"))
			consumer := slices.Index(ids, linker.ID("example/smoke-consumer"))
			if provider < 0 || manager < 0 || consumer < 0 || provider > manager || provider > consumer {
				t.Fatalf("topology does not honor autonomous declarations: %#v", ids)
			}
		})
	}
}

func TestFrameworkTopologyFailsClosedOnMissingDependency(t *testing.T) {
	var events []string
	t.Run("capability", func(t *testing.T) {
		app := linker.New(
			linker.WithMode(linker.Server),
			linker.WithComponents(&smokeConsumer{events: &events}),
		)
		if err := app.Prepare(context.Background()); !errors.Is(err, linker.ErrMissingCapability) {
			t.Fatalf("missing capability: %v", err)
		}
	})
	t.Run("asset target", func(t *testing.T) {
		app := linker.New(
			linker.WithMode(linker.Server),
			linker.WithComponents(&smokeProvider{events: &events}),
		)
		if err := app.Prepare(context.Background()); !errors.Is(err, linker.ErrAssetTarget) {
			t.Fatalf("missing asset target: %v", err)
		}
	})
}

func TestFrameworkStartupRollbackAndGracefulShutdown(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		var events []string
		app := linker.New(
			linker.WithMode(linker.Server),
			linker.WithComponents(
				&smokeProvider{events: &events},
				&smokeConsumer{events: &events, failInit: true},
				&smokeManager{events: &events},
			),
		)
		if err := app.Start(context.Background()); err == nil {
			t.Fatal("expected startup failure")
		}
		if !containsInOrder(events, "provider:mounted", "consumer:init", "provider:stop", "provider:close") {
			t.Fatalf("rollback lifecycle = %#v", events)
		}
	})

	t.Run("shutdown", func(t *testing.T) {
		var events []string
		app := linker.New(
			linker.WithMode(linker.Server),
			linker.WithComponents(
				&smokeManager{events: &events},
				&smokeProvider{events: &events},
			),
		)
		if err := app.Start(context.Background()); err != nil {
			t.Fatalf("start: %v", err)
		}
		if err := app.Stop(context.Background()); err != nil {
			t.Fatalf("stop: %v", err)
		}
		if !containsInOrder(events, "provider:mounted", "manager:init", "manager:stop", "provider:stop", "provider:close") {
			t.Fatalf("graceful lifecycle = %#v", events)
		}
	})
}

func componentIDs(plan linker.Plan) []linker.ID {
	ids := make([]linker.ID, 0, len(plan.Components))
	for _, component := range plan.Components {
		ids = append(ids, component.ID)
	}
	return ids
}

func containsInOrder(values []string, expected ...string) bool {
	index := 0
	for _, value := range values {
		if index < len(expected) && value == expected[index] {
			index++
		}
	}
	return index == len(expected)
}
