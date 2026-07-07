package tts

import (
	"context"
	"fmt"
	"strings"

	rpcmeta "github.com/neteast-software/go-module/rpc/meta"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"

	ttsmodel "linker-v3-example/internal/model/tts"
)

type Service struct {
	db     *gorm.DB
	closed bool
}

func New(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (p *Service) Transcribe(ctx context.Context, req *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	if p == nil || p.db == nil || p.closed {
		return nil, fmt.Errorf("tts 服务未初始化")
	}
	text := strings.TrimSpace(req.Value)
	if text == "" {
		return nil, fmt.Errorf("tts 文本不能为空")
	}
	scope, _ := rpcmeta.Scope(ctx)
	result := "tts:" + text
	if scope != "" {
		result += ":" + scope
	}
	record := ttsmodel.Record{
		Text:   text,
		Result: result,
		Scope:  scope,
	}
	if err := p.db.WithContext(ctx).Create(&record).Error; err != nil {
		return nil, err
	}
	return wrapperspb.String(result), nil
}

func (p *Service) Close() error {
	if p == nil {
		return nil
	}
	p.closed = true
	return nil
}
