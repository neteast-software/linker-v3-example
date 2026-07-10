package example

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	httpclient "github.com/neteast-software/go-module/http/client"
	clientcomponent "github.com/neteast-software/go-module/http/client/linker"
	linker "github.com/neteast-software/linker/v3"

	"linker-v3-example/internal/client/directory"
)

func TestHTTPClientExample(t *testing.T) {
	var attempts []httpclient.Attempt
	external := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/3/badge" {
			t.Fatalf("unexpected external path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer example-token" {
			t.Fatalf("missing authorization: %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("X-Trace-ID") != "trace-http-client-example" {
			t.Fatalf("missing trace header: %s", r.Header.Get("X-Trace-ID"))
		}
		_ = json.NewEncoder(w).Encode(directory.Badge{
			UserID:      "3",
			DisplayName: "示例用户",
			Avatar:      "https://cdn.example.com/avatar/3.png",
		})
	}))
	defer external.Close()

	component := clientcomponent.New(
		clientcomponent.WithConfig(httpclient.Config{
			BaseURL: external.URL,
			Timeout: 2 * time.Second,
			Retry:   httpclient.Retry(2, time.Millisecond),
		}),
		clientcomponent.WithCredential(httpclient.Bearer(httpclient.StaticToken("example-token"))),
		clientcomponent.WithBefore(httpclient.BeforeFunc(func(_ context.Context, req *http.Request) error {
			req.Header.Set("X-Trace-ID", "trace-http-client-example")
			return nil
		})),
		clientcomponent.WithAfter(httpclient.AfterFunc(func(_ context.Context, attempt httpclient.Attempt) {
			attempts = append(attempts, attempt)
		})),
	)
	app := linker.New(
		linker.WithMode(linker.Bin),
		linker.WithComponents(component),
	)
	plan := preparedPlan(t, app)
	if !planHasComponent(plan, clientcomponent.ID) {
		t.Fatalf("plan missing http client component: %#v", plan.Components)
	}
	if !planHasAsset(plan, "http/client", clientcomponent.ID.String()) {
		t.Fatalf("plan missing http client asset: %#v", plan.Assets)
	}

	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		if err := app.Stop(context.Background()); err != nil {
			t.Fatalf("close: %v", err)
		}
	})

	api, err := linker.RequireCapability(app, clientcomponent.ClientKey())
	if err != nil {
		t.Fatalf("http client capability: %v", err)
	}
	badge, err := directory.New(api).Badge(context.Background(), "3")
	if err != nil {
		t.Fatalf("badge: %v", err)
	}
	if badge.DisplayName != "示例用户" || badge.Avatar == "" {
		t.Fatalf("badge = %#v", badge)
	}
	if len(attempts) != 1 || !attempts[0].Success() {
		t.Fatalf("attempts = %#v", attempts)
	}
}
