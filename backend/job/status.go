package job

const (
	Idle int = iota
	Waiting
	Busy
)

type Status struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func NewStatus(name string, status int) Status {
	return Status{
		Name:   name,
		Status: []string{"Idle", "Waiting", "Busy"}[status],
	}
}
