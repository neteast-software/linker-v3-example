package notification

import linker "github.com/neteast-software/linker/v3"

const ProviderID linker.ID = "example/notification/provider"

func ProviderKey() linker.CapabilityKey[*Provider] {
	return linker.NewCapabilityKey[*Provider](ProviderID)
}
