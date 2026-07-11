package tts

import (
	grpclinker "github.com/neteast-software/go-module/rpc/grpc/linker"
	linker "github.com/neteast-software/linker/v3"
)

func ClientKey() linker.CapabilityKey[Client] {
	return grpclinker.ClientKey[Client](ID)
}

func Resolve(runtime linker.Runtime) (Client, bool) {
	return linker.Resolve(runtime, ClientKey())
}

func Require(runtime linker.Runtime) (Client, error) {
	return linker.RequireCapability(runtime, ClientKey())
}
