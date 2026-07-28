// Package client 将 TTS client 装配为 Linker 生命周期组件。
package client

import (
	grpclinker "github.com/neteast-software/go-module/rpc/grpc/linker"
	linker "github.com/neteast-software/linker/v3"

	clientcore "linker-v3-example/internal/tts/client"
)

const ID linker.ID = "example/tts-client"

type Client = clientcore.Client

func Key() linker.CapabilityKey[Client] {
	return grpclinker.ClientKey[Client](ID)
}

func Resolve(runtime linker.Runtime) (Client, bool) {
	return linker.Resolve(runtime, Key())
}

func Require(runtime linker.Runtime) (Client, error) {
	return linker.Require(runtime, Key())
}

func New() linker.Component {
	return grpclinker.NewClientProvider[Client](
		ID,
		clientcore.New,
	)
}
