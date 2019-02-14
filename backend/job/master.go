package job

import (
	"sync"
	"fmt"
)

type JobStatus int
		
const (
			JobStatusIdle JobStatus = iota
			JobStatusBusy
)	

type JobBackupCreate struct {
	app string
	file string
}

type JobBackupRestore struct {
	app string
	file string
}

type Master struct {
	mutex  *sync.Mutex
	status JobStatus
	jobQueue chan interface{}
	feedbackQueue chan interface{}
}

func NewMaster() *Master {
	 
	 master := &Master{
		mutex:  &sync.Mutex{},
		status: JobStatusIdle,
		jobQueue: make(chan interface{}),
		feedbackQueue: make(chan interface{}),
	}
		go master.feedback()
		return master
}

func (master *Master) Status() JobStatus {
	master.mutex.Lock()
	defer master.mutex.Unlock()
	return master.status
}

func (master *Master) BackupCreateJob(app string, file string) error {
	return master.newJob(JobBackupCreate{app: app, file: file})
}

func (master *Master) newJob(job interface{}) error {
	master.mutex.Lock()
	defer master.mutex.Unlock()
	if (master.status == JobStatusIdle) {
	 master.status = JobStatusBusy
  go func() { master.jobQueue <- job }()
		return nil
	} else {
		return fmt.Errorf("busy")
	}
}

func (master *Master) JobQueue() chan interface{} {
 return master.jobQueue
}

func (master *Master) FeedbackQueue() chan interface{} {
 return master.feedbackQueue
}
func (master *Master) feedback() {
	for {
		<- master.feedbackQueue
		master.mutex.Lock()
		master.status = JobStatusIdle
		master.mutex.Unlock()
	}

}