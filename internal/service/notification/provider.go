package notification

import (
	"context"
	"sync"
	"time"
)

type Message struct {
	Time   time.Time
	Source string
	Body   string
}

type Provider struct {
	mu       sync.Mutex
	messages []Message
}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) Send(_ context.Context, source string, body string) error {
	if p == nil {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.messages = append(p.messages, Message{
		Time:   time.Now(),
		Source: source,
		Body:   body,
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
