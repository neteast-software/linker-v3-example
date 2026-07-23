package user

import linker "github.com/neteast-software/linker/v3"

const ServiceID linker.ID = "example/user/service"

func ServiceKey() linker.CapabilityKey[*Service] {
	return linker.NewCapabilityKey[*Service](ServiceID)
}
