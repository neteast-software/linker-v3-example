package user

import linker "github.com/neteast-software/linker/v3"

const AuthID linker.ID = "example/user/auth"

func AuthKey() linker.CapabilityKey[Auth] {
	return linker.NewCapabilityKey[Auth](AuthID)
}
