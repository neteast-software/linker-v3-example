package example

import (
	"context"
	"testing"

	http "github.com/neteast-software/go-module/http/gin/linker"
	linker "github.com/neteast-software/linker/v3"
)

type configReader struct {
	id    linker.ID
	value string
}

func newConfigReader() *configReader {
	return &configReader{id: "example/config-reader"}
}

func (p *configReader) Identity() linker.ID {
	return p.id
}

func (p *configReader) Init(_ context.Context, runtime linker.Runtime) error {
	data, ok := runtime.Setting("example/config")
	if ok {
		p.value = string(data)
	}
	return nil
}

func TestCoreBinDoesNotLoadServerDefaults(t *testing.T) {
	reader := newConfigReader()
	app := linker.New(
		linker.WithMode(linker.Bin),
		linker.WithSource(linker.MapSource{Setting: linker.NewSetting(map[linker.Namespace][]byte{
			"example/config": []byte(`{"name":"bin"}`),
		})}),
		linker.WithComponents(reader),
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
	if reader.value != `{"name":"bin"}` {
		t.Fatalf("reader value = %q", reader.value)
	}
	if _, ok := linker.Resolve(app, linker.NewCapabilityKey[*http.Server](http.ID)); ok {
		t.Fatalf("http capability should be absent")
	}
	if err := app.Stop(context.Background()); err != nil {
		t.Fatalf("stop: %v", err)
	}
}
