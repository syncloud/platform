package smart

type SelfTestLogTable struct {
	Count int             `json:"count"`
	Table []SelfTestEntry `json:"table"`
}
