package inspection

import (
	"context"

	"github.com/neteast-software/go-module/application"
	postgresql "github.com/neteast-software/go-module/db/postgresql"

	inspectionconstant "linker-v3-example/internal/constant/inspection"
	inspectionmodel "linker-v3-example/internal/model/inspection"
)

type Service struct {
	store Store
}

func NewService(store Store) Service {
	return Service{store: store}
}

func (p Service) List(ctx context.Context, app application.Application, req ListRequest) ([]inspectionmodel.Task, ListRequest, error) {
	rows, page, err := p.store.List(ctx, app, req)
	req.Page = page
	return rows, req, err
}

func Seed(ctx context.Context, store Store) error {
	return store.SaveDefaults(ctx,
		inspectionmodel.Task{Head: head(1), ApplicationScope: "app2", Title: "应用二巡检任务", Status: inspectionconstant.Open, OwnerID: 2},
		inspectionmodel.Task{Head: head(2), ApplicationScope: "console", Title: "后台巡检任务", Status: inspectionconstant.Open, OwnerID: 1},
		inspectionmodel.Task{Head: head(3), ApplicationScope: "app2", Title: "应用二已完成任务", Status: inspectionconstant.Done, OwnerID: 1},
	)
}

func head(id uint64) postgresql.Head {
	return postgresql.Head{ID: id}
}
