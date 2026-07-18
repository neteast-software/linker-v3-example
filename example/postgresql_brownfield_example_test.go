package example

import (
	"context"
	"database/sql"
	"testing"

	table "github.com/neteast-software/go-module/db/gorm/table"
	postgresql "github.com/neteast-software/go-module/db/postgresql/linker"
	transition "github.com/neteast-software/go-module/db/postgresql/transition"
	linker "github.com/neteast-software/linker/v3"
)

type legacyAccount struct {
	ID     uint64
	Status string
}

func (p legacyAccount) TableName() string {
	return "legacy_account"
}

type accountSchema struct{}

func (p accountSchema) Identity() linker.ID {
	return "example/account-schema"
}

// Assets 把表所有权和一次性变化集中在拥有它们的 component，不散落到 store。
func (p accountSchema) Assets(context.Context, linker.Runtime) ([]linker.Asset, error) {
	return []linker.Asset{
		postgresql.Table(&legacyAccount{}, postgresql.External(), postgresql.Comment("既有账号表")),
		postgresql.ExternalTable("legacy_audit", postgresql.Comment("外部迁移维护的审计表")),
		postgresql.ReadOnlyTable("public.schema_migrations", postgresql.Comment("外部迁移账本")),
		postgresql.Transition(transition.BeforeSQL(
			"account-status-backfill",
			`UPDATE legacy_account SET status = 'enabled' WHERE status IS NULL`,
		)),
	}, nil
}

// sqlcDBTX 模拟历史 sqlc 生成代码常见的小接口；新项目不以此作为推荐访问方案。
type sqlcDBTX interface {
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type accountQueries struct {
	db sqlcDBTX
}

func newAccountQueries(db sqlcDBTX) *accountQueries {
	return &accountQueries{db: db}
}

func (p *accountQueries) Enable(ctx context.Context, id uint64) error {
	_, err := p.db.ExecContext(ctx, "UPDATE legacy_account SET status = $1 WHERE id = $2", "enabled", id)
	return err
}

// updateAccount 演示棕地项目保留历史 sqlc 查询层，并与 GORM 使用同一个事务。
func updateAccount(ctx context.Context, runtime linker.Runtime, account legacyAccount) error {
	db, err := postgresql.Require(runtime)
	if err != nil {
		return err
	}
	return postgresql.Transaction(ctx, db, func(tx postgresql.Tx) error {
		if err := tx.GORM().Save(&account).Error; err != nil {
			return err
		}
		return newAccountQueries(tx.SQL()).Enable(ctx, account.ID)
	})
}

// readAccounts 演示历史 sqlc 查询直接借用 component 拥有的同一个连接池。
func readAccounts(runtime linker.Runtime) (*accountQueries, error) {
	db, err := postgresql.RequireSQL(runtime)
	if err != nil {
		return nil, err
	}
	return newAccountQueries(db), nil
}

func TestBrownfieldPostgreSQLDeclarations(t *testing.T) {
	assets, err := (accountSchema{}).Assets(context.Background(), nil)
	if err != nil || len(assets) != 4 {
		t.Fatalf("brownfield assets = %#v, %v", assets, err)
	}
	tableAsset, ok := assets[0].Value.(table.Table)
	if !ok || tableAsset.Strategy.Normalize() != table.StrategyExternal {
		t.Fatalf("既有表未声明 External: %#v", assets[0].Value)
	}
	external, ok := assets[1].Value.(table.Table)
	if !ok || external.Name != "legacy_audit" || external.Model != nil ||
		external.Strategy.Normalize() != table.StrategyExternal {
		t.Fatalf("SQL-only 表未声明 External: %#v", assets[1].Value)
	}
	readOnly, ok := assets[2].Value.(table.Table)
	if !ok || readOnly.Name != "public.schema_migrations" || readOnly.Model != nil ||
		readOnly.Strategy.Normalize() != table.StrategyReadOnly {
		t.Fatalf("迁移账本未声明 ReadOnly: %#v", assets[2].Value)
	}
	descriptions := assets[3].Descriptions()
	if len(descriptions) != 1 || descriptions[0].Name != "account-status-backfill" ||
		descriptions[0].Detail["checksum"] == "" {
		t.Fatalf("transition Plan 投影不完整: %#v", descriptions)
	}
}

var _ sqlcDBTX = (postgresql.Executor)(nil)
var _ linker.AssetProvider = accountSchema{}
