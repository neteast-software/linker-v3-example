package inspection

import (
	"github.com/neteast-software/go-module/db/gorm/model"

	inspectionconstant "linker-v3-example/internal/constant/inspection"
)

type Task struct {
	model.Head
	ApplicationScope string                    `gorm:"type:varchar(64);index;not null" json:"application_scope"`
	Title            string                    `gorm:"type:varchar(128);not null" json:"title"`
	Status           inspectionconstant.Status `gorm:"type:varchar(32);index;not null" json:"status"`
	OwnerID          uint64                    `gorm:"index" json:"owner_id"`
}

func (Task) TableName() string {
	return "task"
}

func (p Task) Validate() error {
	if !p.Status.Valid() {
		return inspectionconstant.ErrStatusInvalid
	}
	return nil
}
