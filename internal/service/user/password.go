package user

import useraccount "github.com/neteast-software/go-module/user/account"

const passwordSalt = "linker-v3-example"

var passwordHasher = useraccount.SM3("linker-v3-example")

func passwordHash(password string) (string, error) {
	return passwordHasher.Hash(password, passwordSalt)
}

func verifyPassword(password string, salt string, hash string) (bool, error) {
	return passwordHasher.Verify(password, salt, hash)
}
