package user

import (
	"context"
	"errors"
	"testing"

	usermodel "linker-v3-example/internal/model/user"
)

func TestSaveDefaultsRejectsDuplicateBatchIdentity(t *testing.T) {
	store := Store{}
	err := store.SaveDefaults(context.Background(),
		DefaultUser{User: usermodel.User{Username: "admin"}},
		DefaultUser{User: usermodel.User{Username: "admin"}},
	)
	if !errors.Is(err, errDefaultUserDuplicate) {
		t.Fatalf("err=%v", err)
	}

	err = store.SaveDefaults(context.Background(), DefaultUser{
		User: usermodel.User{Username: "admin"},
		Accounts: []usermodel.Account{
			{Provider: "password", Identifier: "admin"},
			{Provider: "password", Identifier: "admin"},
		},
	})
	if !errors.Is(err, errDefaultAccountDuplicate) {
		t.Fatalf("err=%v", err)
	}
}
