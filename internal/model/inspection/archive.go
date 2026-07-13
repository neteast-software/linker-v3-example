package inspection

import (
	"time"

	postgresql "github.com/neteast-software/go-module/db/postgresql"
)

type Archive struct {
	postgresql.Head
	ApplicationScope string    `gorm:"type:varchar(64);index;not null" json:"application_scope"`
	TaskID           uint64    `gorm:"index;not null" json:"task_id"`
	ArchivedAt       time.Time `gorm:"index;not null" json:"archived_at"`
}

func (Archive) TableName() string {
	return "archive"
}
