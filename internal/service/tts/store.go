package tts

import (
	"context"

	"gorm.io/gorm"

	ttsmodel "linker-v3-example/internal/model/tts"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return Store{db: db}
}

func (p Store) Ready() bool {
	return p.db != nil
}

func (p Store) Save(ctx context.Context, text string, result string, scope string) error {
	return p.db.WithContext(ctx).Create(&ttsmodel.Conversion{
		Text: text, Result: result, Scope: scope,
	}).Error
}
