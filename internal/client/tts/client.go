package tts

import (
	tracegrpc "github.com/neteast-software/go-module/observe/tracing/rpc/grpc"
	grpclinker "github.com/neteast-software/go-module/rpc/grpc/linker"
	linker "github.com/neteast-software/linker/v3"
	"google.golang.org/grpc"

	ttsrpc "linker-v3-example/internal/rpc/tts"
)

const ID linker.ID = "rpc/client/tts"

type Client = ttsrpc.Client

func New(conn grpc.ClientConnInterface) Client {
	return ttsrpc.NewClient(conn)
}

func ClientKey() linker.CapabilityKey[Client] {
	return grpclinker.ClientKey[Client](ID)
}

func Provider(config grpclinker.ClientConfig) linker.Component {
	return grpclinker.NewClientProvider[Client](
		ID,
		New,
		grpclinker.WithClientConfig[Client](config),
		grpclinker.WithClientAfter[Client]("rpc/grpc"),
		grpclinker.WithDialOptions[Client](grpc.WithChainUnaryInterceptor(tracegrpc.UnaryClient())),
	)
}
