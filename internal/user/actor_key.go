package user

import linker "github.com/neteast-software/linker/v3"

const ActorID linker.ID = "example/user/actor"

func ActorKey() linker.CapabilityKey[Actor] {
	return linker.NewCapabilityKey[Actor](ActorID)
}
