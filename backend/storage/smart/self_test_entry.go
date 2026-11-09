package smart

type SelfTestEntry struct {
	Type          SelfTestType   `json:"type"`
	Status        SelfTestStatus `json:"status"`
	LifetimeHours int64          `json:"lifetime_hours"`
}
