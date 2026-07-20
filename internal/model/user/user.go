package user

import "github.com/neteast-software/go-module/db/gorm/model"

type User struct {
	model.Head
	Username string `gorm:"type:varchar(64);uniqueIndex;not null" json:"username"`
	Avatar   string `gorm:"type:varchar(256)" json:"avatar"`
	Email    string `gorm:"type:varchar(128);index" json:"email"`
	Phone    string `gorm:"type:varchar(32);uniqueIndex" json:"phone"`
	Role     string `gorm:"type:varchar(32);index" json:"role"`
}

func (User) TableName() string {
	return "user"
}
