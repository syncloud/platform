package job

import (
	"github.com/syncloud/platform/backup"
	"log"
	"time"
)

type Worker struct {
	backup backup.AppBackup
	master JobMaster
}

func NewWorker(master JobMaster, backup backup.AppBackup) *Worker {
	return &Worker{
		backup: backup,
		master: master,
	}
}

func (worker *Worker) Start() {
	for {
		if !worker.Do() {
			time.Sleep(100)
		}
	}
}

func (worker *Worker) Do() bool {
	job, err := worker.master.Take()
	if err != nil {
		log.Println(err)
		return false
	}
	switch jobtype := job.(type) {
	case JobBackupCreate:
		v := job.(JobBackupCreate)
		worker.backup.Create(v.App, v.File)
	default:
		log.Println("not supported job type", jobtype)
	}
	worker.master.Complete()
	return true
}
