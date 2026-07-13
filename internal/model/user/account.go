package user

import postgresql "github.com/neteast-software/go-module/db/postgresql"

type Account struct {
	postgresql.Head
	UserID     uint64 `gorm:"not null;index" json:"user_id"`
	Provider   string `gorm:"type:varchar(32);uniqueIndex:idx_example_account_identity;not null" json:"provider"`
	Identifier string `gorm:"type:varchar(128);uniqueIndex:idx_example_account_identity;not null" json:"identifier"`
	SecretHash string `gorm:"type:char(64);not null" json:"-"`
	Salt       string `gorm:"type:varchar(64);not null" json:"-"`
}

func (Account) TableName() string {
	return "account"
}
