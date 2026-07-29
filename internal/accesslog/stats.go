package accesslog

// Stats 是 accesslog 有界投递状态。
type Stats struct {
	Published uint64
	Failed    uint64
	Dropped   uint64
	Queued    int
	Capacity  int
}
