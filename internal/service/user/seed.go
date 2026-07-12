package user

import (
	"context"

	userconstant "linker-v3-example/internal/constant/user"
	usermodel "linker-v3-example/internal/model/user"
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
			User: usermodel.User{
				Username: "admin",
				Avatar:   "https://static.neteast.cn/avatar/admin.png",
				Email:    "admin@neteast.cn",
				Phone:    "18000000000",
				Role:     "admin",
			},
			Accounts: []usermodel.Account{
				{Provider: "password", Identifier: "admin", SecretHash: adminHash, Salt: adminSalt},
			},
		},
		DefaultUser{
			User: usermodel.User{
				Username: "example_user",
				Avatar:   "https://static.neteast.cn/avatar/user.png",
				Email:    "example.user@neteast.cn",
				Phone:    userconstant.ExamplePhone,
				Role:     "user",
			},
			Accounts: []usermodel.Account{
				{Provider: "phone", Identifier: userconstant.ExamplePhone, SecretHash: userHash, Salt: userSalt},
			},
		},
	)
}
