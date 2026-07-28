package rpc

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/wrapperspb"

	tts "linker-v3-example/internal/tts"
)

// Server 把 TTS 业务能力适配为 gRPC endpoint。
type Server struct {
	service *tts.Service
}

func New() *Server {
	return &Server{}
}

// Configure 装配 endpoint 使用的 TTS 业务能力。
func (p *Server) Configure(service *tts.Service) {
	p.service = service
}

func (p *Server) Transcribe(ctx context.Context, req *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	if p.service == nil {
		return nil, fmt.Errorf("tts 服务尚未就绪")
	}
	result, err := p.service.Transcribe(ctx, req.GetValue())
	if err != nil {
		return nil, err
	}
	return wrapperspb.String(result), nil
}
