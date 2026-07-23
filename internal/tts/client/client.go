package tts

import (
	grpclinker "github.com/neteast-software/go-module/rpc/grpc/linker"
	linker "github.com/neteast-software/linker/v3"
	"google.golang.org/grpc"

	tts "linker-v3-example/internal/tts"
)

const ID linker.ID = "rpc/client/tts"

type Client = tts.Client

func New(conn grpc.ClientConnInterface) Client {
	return tts.NewClient(conn)
}

func Provider() linker.Component {
	return grpclinker.NewClientProvider[Client](
		ID,
		New,
	)
}
