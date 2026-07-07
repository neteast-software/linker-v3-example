package graph

type order struct {
	ID     uint64 `json:"id"`
	Number string `json:"number"`
	Status string `json:"status"`
	Amount int64  `json:"amount"`
}

func sampleOrders() []order {
	return []order{
		{ID: 1, Number: "NO-20260708-001", Status: "open", Amount: 12800},
		{ID: 2, Number: "NO-20260708-002", Status: "closed", Amount: 6800},
	}
}
