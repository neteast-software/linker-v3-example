package inspection

import (
	"context"
	"testing"

	tablecore "github.com/neteast-software/go-module/db/gorm/table"
	postgresql "github.com/neteast-software/go-module/db/postgresql/linker"
	linker "github.com/neteast-software/linker/v3"
)

func TestAssetsDeclareExternalArchiveTable(t *testing.T) {
	assets, err := NewComponent().Assets(context.Background(), nil)
	if err != nil {
		t.Fatalf("assets: %v", err)
	}
	var found bool
	for _, asset := range assets {
		if asset.Kind != linker.AssetTable || asset.Target != postgresql.ID {
			continue
		}
		table, ok := asset.Value.(tablecore.Table)
		if !ok {
			t.Fatalf("asset value = %#v", asset.Value)
		}
		if table.Strategy == tablecore.StrategyExternal {
			found = true
			if table.Enabled() {
				t.Fatal("external archive table should not auto migrate")
			}
		}
	}
	if !found {
		t.Fatalf("missing external archive table asset: %#v", assets)
	}
}
