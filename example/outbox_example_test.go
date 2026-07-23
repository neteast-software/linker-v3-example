package example

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/neteast-software/go-module/outbox"

	notification "linker-v3-example/internal/notification"
)

func TestOutboxDeliveryExample(t *testing.T) {
	ctx := context.Background()
	store := outbox.NewMemoryStore()
	provider := notification.NewProvider()

	message, err := store.Enqueue(ctx, outbox.New(
		"notification.send",
		[]byte("hello outbox"),
		outbox.WithID("notification-1"),
		outbox.WithKey("notification:1"),
	))
	if err != nil {
		t.Fatalf("enqueue: %v", err)
	}

	dispatcher := outbox.NewDispatcher(store, outbox.HandlerFunc(func(ctx context.Context, message outbox.Message) error {
		return provider.Send(ctx, "outbox", string(message.Payload))
	}))
	result, err := dispatcher.Dispatch(ctx)
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	if result.Delivered != 1 || result.Claimed != 1 {
		t.Fatalf("result = %#v", result)
	}
	got, ok := store.Get(message.ID)
	if !ok || got.Status != outbox.StatusDelivered {
		t.Fatalf("message = %#v ok=%v", got, ok)
	}
	if !hasProviderMessage(provider.Messages(), "outbox") {
		t.Fatalf("provider messages = %#v", provider.Messages())
	}
}

func TestOutboxDeadLetterExample(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 7, 8, 10, 0, 0, 0, time.UTC)
	store := outbox.NewMemoryStore(outbox.WithNow(func() time.Time { return now }))
	if _, err := store.Enqueue(ctx, outbox.New(
		"notification.send",
		[]byte("will fail"),
		outbox.WithID("notification-failed"),
		outbox.WithAvailableAt(now),
		outbox.WithMaxAttempts(1),
	)); err != nil {
		t.Fatalf("enqueue: %v", err)
	}

	sendErr := errors.New("provider unavailable")
	dispatcher := outbox.NewDispatcher(
		store,
		outbox.HandlerFunc(func(context.Context, outbox.Message) error { return sendErr }),
		outbox.WithDispatchNow(func() time.Time { return now }),
	)
	result, err := dispatcher.Dispatch(ctx)
	if !errors.Is(err, sendErr) {
		t.Fatalf("dispatch error = %v", err)
	}
	if result.Dead != 1 || result.Failed != 0 {
		t.Fatalf("result = %#v", result)
	}
	got, ok := store.Get("notification-failed")
	if !ok || got.Status != outbox.StatusDead || got.LastError != sendErr.Error() {
		t.Fatalf("message = %#v ok=%v", got, ok)
	}
}
