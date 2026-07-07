package user

import (
	"context"
	"errors"

	useraccount "github.com/neteast-software/go-module/user/account"
	"gorm.io/gorm"

	userconstant "linker-v3-example/internal/constant/user"
	usermodel "linker-v3-example/internal/model/user"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return Store{db: db}
}

type DefaultUser struct {
	User     usermodel.User
	Accounts []usermodel.Account
}

func (s Store) ByAdmin(ctx context.Context, username string) (usermodel.User, usermodel.Account, error) {
	return s.byAccount(ctx, useraccount.ProviderPassword, username, "admin")
}

func (s Store) ByPhone(ctx context.Context, phone string) (usermodel.User, usermodel.Account, error) {
	return s.byAccount(ctx, useraccount.ProviderPhone, phone, "user")
}

func (s Store) ByID(ctx context.Context, id uint64) (usermodel.User, error) {
	var user usermodel.User
	err := s.db.WithContext(ctx).First(&user, id).Error
	return user, err
}

func (s Store) SaveDefaults(ctx context.Context, users ...DefaultUser) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range users {
			record := item.User
			if err := tx.
				Where("username = ?", item.User.Username).
				Assign(item.User).
				FirstOrCreate(&record).Error; err != nil {
				return err
			}
			for _, account := range item.Accounts {
				account.UserID = record.ID
				if err := tx.
					Where("provider = ? AND identifier = ?", account.Provider, account.Identifier).
					Assign(account).
					FirstOrCreate(&account).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s Store) byAccount(ctx context.Context, provider useraccount.Provider, identifier string, role string) (usermodel.User, usermodel.Account, error) {
	var account usermodel.Account
	err := s.db.WithContext(ctx).
		Where("provider = ? AND identifier = ?", provider.String(), identifier).
		First(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return usermodel.User{}, usermodel.Account{}, userconstant.ErrLogin
	}
	if err != nil {
		return usermodel.User{}, usermodel.Account{}, err
	}
	user, err := s.first(ctx, "id = ? AND role = ?", account.UserID, role)
	if err != nil {
		return usermodel.User{}, usermodel.Account{}, err
	}
	return user, account, nil
}

func (s Store) first(ctx context.Context, query string, args ...any) (usermodel.User, error) {
	var user usermodel.User
	err := s.db.WithContext(ctx).Where(query, args...).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return usermodel.User{}, userconstant.ErrLogin
	}
	return user, err
}
