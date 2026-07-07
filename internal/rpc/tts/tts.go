package tts

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const FullMethodTranscribe = "/example.tts.TTS/Transcribe"

type Client interface {
	Transcribe(ctx context.Context, text string, opts ...grpc.CallOption) (string, error)
}

type Server interface {
	Transcribe(context.Context, *wrapperspb.StringValue) (*wrapperspb.StringValue, error)
}

type client struct {
	conn grpc.ClientConnInterface
}

func NewClient(conn grpc.ClientConnInterface) Client {
	return client{conn: conn}
}

func (p client) Transcribe(ctx context.Context, text string, opts ...grpc.CallOption) (string, error) {
	resp := new(wrapperspb.StringValue)
	err := p.conn.Invoke(ctx, FullMethodTranscribe, wrapperspb.String(text), resp, opts...)
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

func Register(server grpc.ServiceRegistrar, service Server) {
	server.RegisterService(&grpc.ServiceDesc{
		ServiceName: "example.tts.TTS",
		HandlerType: (*Server)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Transcribe",
				Handler:    transcribeHandler,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "example/tts",
	}, service)
}

func transcribeHandler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	req := new(wrapperspb.StringValue)
	if err := dec(req); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Server).Transcribe(ctx, req)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FullMethodTranscribe,
	}
	handler := func(ctx context.Context, value any) (any, error) {
		return srv.(Server).Transcribe(ctx, value.(*wrapperspb.StringValue))
	}
	return interceptor(ctx, req, info, handler)
}
