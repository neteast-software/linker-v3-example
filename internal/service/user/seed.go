package user

import (
	"context"

	"linker-v3-example/internal/config"
	usermodel "linker-v3-example/internal/model/user"
)

func Seed(ctx context.Context, store Store) error {
	adminHash, err := passwordHash(config.ExampleLoginPassword)
	if err != nil {
		return err
	}
	userHash, err := passwordHash(config.ExampleLoginPassword)
	if err != nil {
		return err
	}
	return store.SaveDefaults(ctx,
		DefaultUser{
			User: usermodel.User{
				Username: "admin",
				Avatar:   "https://static.neteast.cn/avatar/admin.png",
				Email:    "admin@neteast.cn",
				Phone:    "18000000000",
				Role:     "admin",
			},
			Accounts: []usermodel.Account{
				{Provider: "password", Identifier: "admin", SecretHash: adminHash, Salt: passwordSalt},
			},
		},
		DefaultUser{
			User: usermodel.User{
				Username: "example_user",
				Avatar:   "https://static.neteast.cn/avatar/user.png",
				Email:    "example.user@neteast.cn",
				Phone:    config.ExampleUserPhone,
				Role:     "user",
			},
			Accounts: []usermodel.Account{
				{Provider: "phone", Identifier: config.ExampleUserPhone, SecretHash: userHash, Salt: passwordSalt},
			},
		},
	)
}
