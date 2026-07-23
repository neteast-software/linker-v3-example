package tts

import "github.com/neteast-software/go-module/db/gorm/model"

// Conversion is the persisted TTS conversion asset.
type Conversion struct {
	model.Head
	Text   string `gorm:"comment:输入文本;type:varchar(256)" json:"text"`
	Result string `gorm:"comment:转换结果;type:varchar(256)" json:"result"`
	Scope  string `gorm:"comment:RPC scope;type:varchar(64);index" json:"scope"`
}

func (Conversion) TableName() string {
	return "conversion"
}
