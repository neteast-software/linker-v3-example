package tts

import "time"

// Conversion is the persisted TTS conversion asset.
type Conversion struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"comment:更新时间" json:"updated_at"`
	Text      string    `gorm:"comment:输入文本;type:varchar(256)" json:"text"`
	Result    string    `gorm:"comment:转换结果;type:varchar(256)" json:"result"`
	Scope     string    `gorm:"comment:RPC scope;type:varchar(64);index" json:"scope"`
}

func (Conversion) TableName() string {
	return "conversion"
}
