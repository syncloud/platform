package job

type JobStatus int

const (
	JobStatusIdle JobStatus = iota
	JobStatusWaiting
	JobStatusBusy
)

func (status JobStatus) String() string {
	names := []string{
		"JobStatusIdle",
		"JobStatusWaiting",
		"JobStatusBusy",
	}
	return names[status]
}

type JobMaster interface {
	Status() JobStatus
	Offer(job func()) error
	Take() (func(), error)
	Complete() error
}
