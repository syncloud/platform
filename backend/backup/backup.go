package backup

import (
		"log"
		"io/ioutil"
		"os"
)

type Backup struct {
	backupDir string
}

const BACKUP_DIR = "/data/platform/backup"

func NewDefault() *Backup {
	return New(BACKUP_DIR)
}

func New(dir string) *Backup {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm);
	}
	return &Backup {
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
