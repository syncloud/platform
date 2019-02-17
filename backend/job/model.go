package job

type JobStatus int

const (
	JobStatusIdle JobStatus = iota
	JobStatusWaiting
	JobStatusBusy
)

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
