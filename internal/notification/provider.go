package notification

import (
	"context"
	"sync"
	"time"

	"github.com/neteast-software/go-module/observe/tracing"
)

type Message struct {
	Time      time.Time
	Source    string
	Body      string
	TraceID   string
	RequestID string
}

type Provider struct {
	mu       sync.Mutex
	messages []Message
}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) Send(ctx context.Context, source string, body string) error {
	if p == nil {
		return nil
	}
	trace, _ := tracing.FromContext(ctx)
	p.mu.Lock()
	defer p.mu.Unlock()
	p.messages = append(p.messages, Message{
		Time:      time.Now(),
		Source:    source,
		Body:      body,
		TraceID:   trace.TraceID,
		RequestID: trace.RequestID,
	})
	return nil
}

func (p *Provider) Messages() []Message {
	if p == nil {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	ret := make([]Message, len(p.messages))
	copy(ret, p.messages)
	return ret
}
