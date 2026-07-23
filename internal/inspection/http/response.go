package inspection

import inspection "linker-v3-example/internal/inspection"

type taskItem struct {
	ID               uint64            `json:"id"`
	ApplicationScope string            `json:"application_scope"`
	Title            string            `json:"title"`
	Status           inspection.Status `json:"status"`
	OwnerID          uint64            `json:"owner_id"`
}

func taskItems(rows []inspection.Task) []taskItem {
	ret := make([]taskItem, 0, len(rows))
	for _, row := range rows {
		ret = append(ret, taskItem{
			ID:               row.ID,
			ApplicationScope: row.ApplicationScope,
			Title:            row.Title,
			Status:           row.Status,
			OwnerID:          row.OwnerID,
		})
	}
	return ret
}
