package storage

import "github.com/syncloud/platform/job"

type JobsStub struct {
	status job.Status
}

func (j *JobsStub) Status() job.Status {
	return j.status
}
