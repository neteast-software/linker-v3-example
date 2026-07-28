package client

import (
	"google.golang.org/grpc"

	tts "linker-v3-example/internal/tts"
)

type Client = tts.Client

func New(conn grpc.ClientConnInterface) Client {
	return tts.NewClient(conn)
}
