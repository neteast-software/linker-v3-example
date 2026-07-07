package example

import (
	"context"
	"testing"

	rpcgrpc "github.com/neteast-software/go-module/rpc/grpc"
	rpcmeta "github.com/neteast-software/go-module/rpc/meta"
	stdgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type exampleMetaServer interface {
	Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
}

type exampleMetaService struct {
	scope string
	user  rpcmeta.User
}

func (p *exampleMetaService) Ping(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	p.scope, _ = rpcmeta.Scope(ctx)
	p.user, _ = rpcmeta.UserFrom(ctx)
	return &emptypb.Empty{}, nil
}

func TestGRPCMetadataExample(t *testing.T) {
	service := &exampleMetaService{}
	server := rpcgrpc.NewServerWithOptions(
		rpcgrpc.ServerConfig{Addr: "127.0.0.1:0"},
		[]stdgrpc.ServerOption{stdgrpc.UnaryInterceptor(rpcgrpc.UnaryServerMeta())},
		func(server *stdgrpc.Server) {
			registerExampleMeta(server, service)
		},
	)
	if err := server.Start(context.Background()); err != nil {
		t.Fatalf("start grpc server: %v", err)
	}
	t.Cleanup(func() {
		_ = server.Stop(context.Background())
	})

	conn, err := rpcgrpc.NewClientConn(
		context.Background(),
		rpcgrpc.ClientConfig{Target: server.Addr()},
		stdgrpc.WithTransportCredentials(insecure.NewCredentials()),
		stdgrpc.WithChainUnaryInterceptor(rpcgrpc.UnaryClientMeta(rpcmeta.System("front"))),
	)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	t.Cleanup(func() {
		_ = conn.Close()
	})

	if err = conn.Invoke(context.Background(), "/example.Meta/Ping", &emptypb.Empty{}, &emptypb.Empty{}); err != nil {
		t.Fatalf("invoke: %v", err)
	}
	if service.scope != "front" || service.user.ID != "0" || service.user.Name != "system" {
		t.Fatalf("metadata not propagated: scope=%s user=%+v", service.scope, service.user)
	}
}

func registerExampleMeta(server *stdgrpc.Server, service exampleMetaServer) {
	server.RegisterService(&stdgrpc.ServiceDesc{
		ServiceName: "example.Meta",
		HandlerType: (*exampleMetaServer)(nil),
		Methods: []stdgrpc.MethodDesc{{
			MethodName: "Ping",
			Handler:    exampleMetaPingHandler,
		}},
	}, service)
}

func exampleMetaPingHandler(server any, ctx context.Context, decode func(any) error, interceptor stdgrpc.UnaryServerInterceptor) (any, error) {
	request := &emptypb.Empty{}
	if err := decode(request); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return server.(exampleMetaServer).Ping(ctx, request)
	}
	info := &stdgrpc.UnaryServerInfo{
		Server:     server,
		FullMethod: "/example.Meta/Ping",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return server.(exampleMetaServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, request, info, handler)
}
