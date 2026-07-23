package inspection

import (
	"context"

	"github.com/neteast-software/go-module/application"
	"github.com/neteast-software/go-module/db/gorm/model"
)

type Service struct {
	store Store
}

func NewService(store Store) Service {
	return Service{store: store}
}

func (p Service) List(ctx context.Context, app application.Application, req ListRequest) ([]Task, ListRequest, error) {
	rows, page, err := p.store.List(ctx, app, req)
	req.Page = page
	return rows, req, err
}

func Seed(ctx context.Context, store Store) error {
	return store.SaveDefaults(ctx,
		Task{Head: head(1), ApplicationScope: "app2", Title: "应用二巡检任务", Status: Open, OwnerID: 2},
		Task{Head: head(2), ApplicationScope: "console", Title: "后台巡检任务", Status: Open, OwnerID: 1},
		Task{Head: head(3), ApplicationScope: "app2", Title: "应用二已完成任务", Status: Done, OwnerID: 1},
	)
}

func head(id uint64) model.Head {
	return model.Head{ID: id}
}
