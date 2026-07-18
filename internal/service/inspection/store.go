package inspection

import (
	"context"
	"slices"

	"github.com/neteast-software/go-module/acl"
	"github.com/neteast-software/go-module/application"
	appstore "github.com/neteast-software/go-module/application/store/gorm"
	"github.com/neteast-software/go-module/db/gorm/query"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	inspectionconstant "linker-v3-example/internal/constant/inspection"
	inspectionmodel "linker-v3-example/internal/model/inspection"
)

type Store struct {
	db *gorm.DB
}

type ListRequest struct {
	Page   query.Page
	Status inspectionconstant.Status
	Access TaskAccess
}

func NewStore(db *gorm.DB) Store {
	return Store{db: db}
}

func (p Store) SaveDefaults(ctx context.Context, tasks ...inspectionmodel.Task) error {
	if len(tasks) == 0 {
		return nil
	}
	for _, task := range tasks {
		if err := task.Validate(); err != nil {
			return err
		}
	}
	records := slices.Clone(tasks)
	return p.db.WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&records).
		Error
}

func (p Store) List(ctx context.Context, app application.Application, req ListRequest) ([]inspectionmodel.Task, query.Page, error) {
	if !req.Access.Can(acl.Read) {
		return nil, query.Page{}, ErrForbidden
	}
	spec := query.Spec{
		Page:   req.Page,
		Orders: []query.Order{query.Desc("id")},
	}
	if req.Status != "" {
		spec.Filters = append(spec.Filters, query.Where("status", req.Status.String()))
	}
	spec.Filters = append(spec.Filters, req.Access.Range.Filters()...)
	db, spec, err := query.Apply(
		p.db.WithContext(ctx).
			Model(&inspectionmodel.Task{}).
			Scopes(appstore.Scope(application.NewDataScope(app))),
		spec,
	)
	if err != nil {
		return nil, query.Page{}, err
	}
	var rows []inspectionmodel.Task
	if err = db.Find(&rows).Error; err != nil {
		return nil, query.Page{}, err
	}
	return rows, spec.Page, nil
}
