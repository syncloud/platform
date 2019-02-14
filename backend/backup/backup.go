package backup

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type Backup struct {
	backupDir string
}

const (
	BACKUP_DIR = "/data/platform/backup"
	BACKUP_CREATE_CMD = "/snap/platform/current/bin/backup.sh"
)

func NewDefault() *Backup {
	return New(BACKUP_DIR)
}

func New(dir string) *Backup {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	return &Backup{
		backupDir: dir,
	}
}

func (this *Backup) List() ([]string, error) {
	files, err := ioutil.ReadDir(this.backupDir)
	if err != nil {
		log.Println("Cannot get list of files in ", this.backupDir, err)
		return nil, err
	}
	var names []string
	for _, x := range files {
		names = append(names, x.Name())
	}

	return names, nil
}

func (backup *Backup) Create(app string, file string) {
	cmd := exec.Command(BACKUP_CREATE_CMD, app, file)
	log.Println("Running backup", BACKUP_CREATE_CMD, app, file)
	err := cmd.Run()
	log.Printf("Backup finished: %v", err)
}
