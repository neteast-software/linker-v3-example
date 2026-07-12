package tts

import (
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

func Provider() linker.Component {
	return grpclinker.NewClientProvider[Client](
		ID,
		New,
	)
}
