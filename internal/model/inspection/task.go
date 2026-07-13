package inspection

import postgresql "github.com/neteast-software/go-module/db/postgresql"

type Task struct {
	postgresql.Head
	ApplicationScope string `gorm:"type:varchar(64);index;not null" json:"application_scope"`
	Title            string `gorm:"type:varchar(128);not null" json:"title"`
	Status           string `gorm:"type:varchar(32);index;not null" json:"status"`
	OwnerID          uint64 `gorm:"index" json:"owner_id"`
}

func (Task) TableName() string {
	return "task"
}
