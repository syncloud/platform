package model

type CheckSchedule struct {
	Scrub bool `json:"scrub"`
	Smart bool `json:"smart"`
}
