package job

import (
	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/installer"
	"log"
	"time"
)

type Worker struct {
	backup backup.AppBackup
 installer installer.AppInstaller
 master JobMaster
}

func NewWorker(master JobMaster, backup backup.AppBackup, installer installer.AppInstaller) *Worker {
	return &Worker{backup, installer, master}
}

func (worker *Worker) Start() {
	for {
		if !worker.Do() {
			time.Sleep(1 * time.Second)
		}
	}
}

func (worker *Worker) Do() bool {
	job, err := worker.master.Take()
	if err != nil {
		return false
	}
	switch jobtype := job.(type) {
	case JobBackupCreate:
		v := job.(JobBackupCreate)
		worker.backup.Create(v.App)
	case JobBackupRestore:
		v := job.(JobBackupRestore)
		worker.backup.Restore(v.File)
 case JobInstallerUpgrade:
		worker.installer.Upgrade()
	default:
		log.Println("not supported job type", jobtype)
	}
	worker.master.Complete()
	return true
}
