package user

import (
	"crypto/rand"
	"encoding/hex"

	useraccount "github.com/neteast-software/go-module/user/account"
)

const passwordNamespace = "linker-v3-example"
const passwordSaltBytes = 16

func passwordHash(password string) (string, string, error) {
	salt, err := newPasswordSalt()
	if err != nil {
		return "", "", err
	}
	hash, err := useraccount.SM3(passwordNamespace).Hash(password, salt)
	return hash, salt, err
}

func verifyPassword(password string, salt string, hash string) (bool, error) {
	return useraccount.SM3(passwordNamespace).Verify(password, salt, hash)
}

func newPasswordSalt() (string, error) {
	value := make([]byte, passwordSaltBytes)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return hex.EncodeToString(value), nil
}
