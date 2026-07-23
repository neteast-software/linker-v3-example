package user

import (
	"context"
	"errors"
	"testing"
)

func TestSaveDefaultsRejectsDuplicateBatchIdentity(t *testing.T) {
	store := Store{}
	err := store.SaveDefaults(context.Background(),
		DefaultUser{User: User{Username: "admin"}},
		DefaultUser{User: User{Username: "admin"}},
	)
	if !errors.Is(err, errDefaultUserDuplicate) {
		t.Fatalf("err=%v", err)
	}

	err = store.SaveDefaults(context.Background(), DefaultUser{
		User: User{Username: "admin"},
		Accounts: []Account{
			{Provider: "password", Identifier: "admin"},
			{Provider: "password", Identifier: "admin"},
		},
	})
	if !errors.Is(err, errDefaultAccountDuplicate) {
		t.Fatalf("err=%v", err)
	}
}
