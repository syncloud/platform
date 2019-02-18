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
	Offer(job interface{}) error
	Take() (interface{}, error)
	Complete() error
}

type JobBackupCreate struct {
	App  string
	File string
}

type JobBackupRestore struct {
	App  string
	File string
}
