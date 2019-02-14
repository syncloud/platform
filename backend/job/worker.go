package job

import 	(
 "log"
 "github.com/syncloud/platform/backup"
)

type Worker struct {
	queue chan interface{}
 backup *backup.Backup
}

func NewWorker(queue chan interface{}, backup *backup.Backup) *Worker {
 return &Worker {
		queue: queue,
 backup: backup,
	}
}

func (worker *Worker) Start() {
	go func() {
		for {
    job := <- worker.queue
    switch jobtype := job.(type) {
							case JobBackupCreate:
       v := job.(JobBackupCreate)
       worker.backup.Create(v.app, v.file)
							default:
       log.Println("not supported job type", jobtype)
						}
			}
	}()
}
