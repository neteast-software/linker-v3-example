package tts

import (
	"context"
	"fmt"
	"strings"

	rpcmeta "github.com/neteast-software/go-module/rpc/meta"
)

type Service struct {
	store  Store
	closed bool
}

func New(store Store) *Service {
	return &Service{store: store}
}

func (p *Service) Transcribe(ctx context.Context, value string) (string, error) {
	if p == nil || !p.store.Ready() || p.closed {
		return "", fmt.Errorf("tts 服务未初始化")
	}
	text := strings.TrimSpace(value)
	if text == "" {
		return "", fmt.Errorf("tts 文本不能为空")
	}
	scope, _ := rpcmeta.Scope(ctx)
	result := "tts:" + text
	if scope != "" {
		result += ":" + scope
	}
	if err := p.store.Save(ctx, text, result, scope); err != nil {
		return "", err
	}
	return result, nil
}

func (p *Service) Close() error {
	if p == nil {
		return nil
	}
	p.closed = true
	return nil
}
