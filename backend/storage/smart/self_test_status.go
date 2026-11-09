package smart

type SelfTestStatus struct {
	Value            int    `json:"value"`
	String           string `json:"string"`
	Passed           *bool  `json:"passed"`
	RemainingPercent *int   `json:"remaining_percent"`
}
