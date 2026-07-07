package inspection

import "github.com/neteast-software/go-module/db/gorm/query"

type RecordRange struct {
	OwnerID uint64
}

func OwnerRange(ownerID uint64) RecordRange {
	return RecordRange{OwnerID: ownerID}
}

func (p RecordRange) Filters() []query.Filter {
	if p.OwnerID == 0 {
		return nil
	}
	return []query.Filter{query.Where("owner_id", p.OwnerID)}
}
