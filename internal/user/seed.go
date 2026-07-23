package user

import (
	"context"
)

func Seed(ctx context.Context, store Store, password string) error {
	adminHash, adminSalt, err := passwordHash(password)
	if err != nil {
		return err
	}
	userHash, userSalt, err := passwordHash(password)
	if err != nil {
		return err
	}
	return store.SaveDefaults(ctx,
		DefaultUser{
			User: User{
				Username: "admin",
				Avatar:   "https://static.neteast.cn/avatar/admin.png",
				Email:    "admin@neteast.cn",
				Phone:    "18000000000",
				Role:     "admin",
			},
			Accounts: []Account{
				{Provider: "password", Identifier: "admin", SecretHash: adminHash, Salt: adminSalt},
			},
		},
		DefaultUser{
			User: User{
				Username: "example_user",
				Avatar:   "https://static.neteast.cn/avatar/user.png",
				Email:    "example.user@neteast.cn",
				Phone:    SeedPhone,
				Role:     "user",
			},
			Accounts: []Account{
				{Provider: "phone", Identifier: SeedPhone, SecretHash: userHash, Salt: userSalt},
			},
		},
	)
}
