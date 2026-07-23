package inspection

import (
	"time"

	"github.com/neteast-software/go-module/db/gorm/model"
)

type Archive struct {
	model.Head
	ApplicationScope string    `gorm:"type:varchar(64);index;not null" json:"application_scope"`
	TaskID           uint64    `gorm:"index;not null" json:"task_id"`
	ArchivedAt       time.Time `gorm:"index;not null" json:"archived_at"`
}

func (Archive) TableName() string {
	return "archive"
}
