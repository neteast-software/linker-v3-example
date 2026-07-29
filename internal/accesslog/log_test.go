package accesslog

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/neteast-software/go-module/http/gateway"
)

func TestAccessLogPublishesOnlyStableRedactedFacts(t *testing.T) {
	publisher := Memory()
	log, err := New(publisher, 4)
	if err != nil {
		t.Fatal(err)
	}
	if err = log.Open(); err != nil {
		t.Fatal(err)
	}
	_, current := log.Start(context.Background(), gateway.RequestStart{
		Route:  "equipment/profile",
		Method: "GET",
	})
	current.Finish(gateway.RequestResult{
		Route: "equipment/profile", Method: "GET", Status: 200, Duration: 25 * time.Millisecond,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err = log.Close(ctx); err != nil {
		t.Fatal(err)
	}

	records := publisher.Records()
	if len(records) != 1 {
		t.Fatalf("records = %d", len(records))
	}
	record := records[0]
	if record.Topic != "operateLog" || record.RequestURL != "equipment/profile" ||
		record.RequestMethod != "GET" || !record.Success || record.ExecutionTime != 25 {
		t.Fatalf("record = %#v", record)
	}
	if record.RequestParams != "" || record.ResponseData != "" || record.UserAgent != "" ||
		strings.Contains(record.SimpleErrorMsg, "token") {
		t.Fatalf("record contains request content: %#v", record)
	}
}

func TestAccessLogQueueIsBoundedAndShutdownDrains(t *testing.T) {
	publisher := &blockingPublisher{
		entered: make(chan struct{}, 1),
		release: make(chan struct{}),
	}
	log, err := New(publisher, 1)
	if err != nil {
		t.Fatal(err)
	}
	if err = log.Open(); err != nil {
		t.Fatal(err)
	}
	finish(log, 200)
	<-publisher.entered
	finish(log, 200)
	finish(log, 503)
	if log.Stats().Dropped != 1 {
		t.Fatalf("stats before drain = %#v", log.Stats())
	}
	close(publisher.release)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err = log.Close(ctx); err != nil {
		t.Fatal(err)
	}
	if stats := log.Stats(); stats.Published != 2 || stats.Dropped != 1 || stats.Queued != 0 {
		t.Fatalf("stats after drain = %#v", stats)
	}
}

func TestAccessLogPublishFailureIsObservable(t *testing.T) {
	log, err := New(PublisherFunc(func(context.Context, Record) error {
		return errors.New("broker unavailable")
	}), 1)
	if err != nil {
		t.Fatal(err)
	}
	if err = log.Open(); err != nil {
		t.Fatal(err)
	}
	finish(log, 502)
	deadline := time.Now().Add(time.Second)
	for log.Stats().Failed == 0 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if log.Health() == nil || log.Stats().Failed != 1 {
		t.Fatalf("health=%v stats=%#v", log.Health(), log.Stats())
	}
	if err = log.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func finish(log *Log, status int) {
	_, current := log.Start(context.Background(), gateway.RequestStart{Route: "route-one", Method: "GET"})
	current.Finish(gateway.RequestResult{
		Route: "route-one", Method: "GET", Status: status, Duration: time.Millisecond,
	})
}

type blockingPublisher struct {
	entered chan struct{}
	release chan struct{}
}

func (p *blockingPublisher) Publish(ctx context.Context, _ Record) error {
	select {
	case p.entered <- struct{}{}:
	default:
	}
	select {
	case <-p.release:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
