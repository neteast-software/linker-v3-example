package user

import (
	"context"
	"errors"
	"fmt"

	useraccount "github.com/neteast-software/go-module/user/account"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

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

func (p Store) ByAdmin(ctx context.Context, username string) (usermodel.User, usermodel.Account, error) {
	return p.byAccount(ctx, useraccount.ProviderPassword, username, "admin")
}

func (p Store) ByPhone(ctx context.Context, phone string) (usermodel.User, usermodel.Account, error) {
	return p.byAccount(ctx, useraccount.ProviderPhone, phone, "user")
}

func (p Store) ByID(ctx context.Context, id uint64) (usermodel.User, error) {
	var user usermodel.User
	err := p.db.WithContext(ctx).First(&user, id).Error
	return user, err
}

func (p Store) SaveDefaults(ctx context.Context, defaults ...DefaultUser) error {
	if len(defaults) == 0 {
		return nil
	}
	users := make([]usermodel.User, len(defaults))
	usernames := make([]string, len(defaults))
	accountCount := 0
	seenUsers := make(map[string]struct{}, len(defaults))
	type accountIdentity struct {
		provider   string
		identifier string
	}
	seenAccounts := make(map[accountIdentity]struct{})
	for index, item := range defaults {
		if item.User.Username == "" {
			return errors.New("默认用户 username 不能为空")
		}
		if _, exists := seenUsers[item.User.Username]; exists {
			return fmt.Errorf("%w: username=%q", errDefaultUserDuplicate, item.User.Username)
		}
		seenUsers[item.User.Username] = struct{}{}
		users[index] = item.User
		usernames[index] = item.User.Username
		accountCount += len(item.Accounts)
		for _, account := range item.Accounts {
			if account.Provider == "" || account.Identifier == "" {
				return fmt.Errorf("默认用户 %q 的账号标识不完整", item.User.Username)
			}
			identity := accountIdentity{provider: account.Provider, identifier: account.Identifier}
			if _, exists := seenAccounts[identity]; exists {
				return fmt.Errorf("%w: provider=%q identifier=%q", errDefaultAccountDuplicate, account.Provider, account.Identifier)
			}
			seenAccounts[identity] = struct{}{}
		}
	}
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "username"}},
			DoUpdates: clause.AssignmentColumns([]string{"avatar", "email", "phone", "role", "updated_at"}),
		}).Create(&users).Error; err != nil {
			return fmt.Errorf("初始化默认用户失败: %w", err)
		}

		var identities []struct {
			ID       uint64
			Username string
		}
		if err := tx.Model(&usermodel.User{}).
			Select("id", "username").
			Where("username IN ?", usernames).
			Find(&identities).Error; err != nil {
			return fmt.Errorf("读取默认用户标识失败: %w", err)
		}
		ids := make(map[string]uint64, len(identities))
		for _, identity := range identities {
			ids[identity.Username] = identity.ID
		}

		accounts := make([]usermodel.Account, 0, accountCount)
		for _, item := range defaults {
			userID, ok := ids[item.User.Username]
			if !ok {
				return fmt.Errorf("默认用户 %q 未返回标识", item.User.Username)
			}
			for _, account := range item.Accounts {
				account.UserID = userID
				accounts = append(accounts, account)
			}
		}
		if len(accounts) == 0 {
			return nil
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "provider"}, {Name: "identifier"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"user_id", "secret_hash", "salt", "updated_at",
			}),
		}).Create(&accounts).Error; err != nil {
			return fmt.Errorf("初始化默认账号失败: %w", err)
		}
		return nil
	})
}

func (p Store) byAccount(ctx context.Context, provider useraccount.Provider, identifier string, role string) (usermodel.User, usermodel.Account, error) {
	var account usermodel.Account
	err := p.db.WithContext(ctx).
		Where("provider = ? AND identifier = ?", provider.String(), identifier).
		First(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return usermodel.User{}, usermodel.Account{}, ErrLogin
	}
	if err != nil {
		return usermodel.User{}, usermodel.Account{}, err
	}
	user, err := p.first(ctx, "id = ? AND role = ?", account.UserID, role)
	if err != nil {
		return usermodel.User{}, usermodel.Account{}, err
	}
	return user, account, nil
}

func (p Store) first(ctx context.Context, query string, args ...any) (usermodel.User, error) {
	var user usermodel.User
	err := p.db.WithContext(ctx).Where(query, args...).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return usermodel.User{}, ErrLogin
	}
	return user, err
}
