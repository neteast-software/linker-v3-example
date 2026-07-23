package inspection

import (
	"github.com/neteast-software/go-module/db/gorm/model"
)

type Task struct {
	model.Head
	ApplicationScope string `gorm:"type:varchar(64);index;not null" json:"application_scope"`
	Title            string `gorm:"type:varchar(128);not null" json:"title"`
	Status           Status `gorm:"type:varchar(32);index;not null" json:"status"`
	OwnerID          uint64 `gorm:"index" json:"owner_id"`
}

func (Task) TableName() string {
	return "task"
}

func (p Task) Validate() error {
	if !p.Status.Valid() {
		return ErrStatusInvalid
	}
	return nil
}
