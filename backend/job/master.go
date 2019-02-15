package job

import (
	"fmt"
	"sync"
)

type JobStatus int

const (
	JobStatusIdle JobStatus = iota
 JobStatusWaiting
	JobStatusBusy
)

type JobBackupCreate struct {
	app  string
	file string
}

type JobBackupRestore struct {
	app  string
	file string
}

type Master struct {
	mutex         *sync.Mutex
	status        JobStatus
	job     interface{}

}

func NewMaster() *Master {

	master := &Master{
		mutex:         &sync.Mutex{},
		status:        JobStatusIdle,
		job:      nil,
	}
	return master
}

func (master *Master) Status() JobStatus {
	master.mutex.Lock()
	defer master.mutex.Unlock()
	return master.status
}

func (master *Master) BackupCreateJob(app string, file string) error {
	return master.offer(JobBackupCreate{app: app, file: file})
}

func (master *Master) offer(job interface{}) error {
	master.mutex.Lock()
	defer master.mutex.Unlock()
	if master.status == JobStatusIdle {
		master.status = JobStatusWaiting
		master.job = job
		return nil
	} else {
		return fmt.Errorf("busy")
	}
}

func (master *Master) Take() (interface{}, error) {
 
 master.mutex.Lock()
 defer master.mutex.Unlock()
 if master.status == JobStatusWaiting {
   master.status = JobStatusBusy
   return master.job, nil
 } else {
  return nil, fmt.Errorf("busy")
 }
}

func (master *Master) Complete() error {
 master.mutex.Lock()
 defer master.mutex.Unlock()
 if master.status == JobStatusBusy {
	  	master.status = JobStatusIdle
  master.job = nil
   return nil
	 } else {
   return fmt.Errorf("nothing to complete")
 }
  
}
