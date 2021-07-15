package backup

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Backup struct {
	backupDir string
}

const (
	BACKUP_DIR         = "/data/platform/backup"
	BACKUP_CREATE_CMD  = "/snap/platform/current/bin/backup.sh"
	BACKUP_RESTORE_CMD = "/snap/platform/current/bin/restore.sh"
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

func (backup *Backup) List() ([]File, error) {
	files, err := ioutil.ReadDir(backup.backupDir)
	if err != nil {
		log.Println("Cannot get list of files in ", backup.backupDir, err)
		return nil, err
	}
	var names []File
	for _, x := range files {
		names = append(names, File{backup.backupDir, x.Name()})
	}

	return names, nil
}

func (backup *Backup) Create(app string) {
	time := time.Now().Format("2006-0102-150405")
	file := fmt.Sprintf("%s/%s-%s.tar.gz", backup.backupDir, app, time)
	log.Println("Running backup create", BACKUP_CREATE_CMD, app, file)
	out, err := exec.Command(BACKUP_CREATE_CMD, app, file).CombinedOutput()
	log.Printf("Backup create output %s", out)
	if err != nil {
		log.Printf("Backup create failed: %v", err)
	} else {
		log.Printf("Backup create completed")
	}
}

func (backup *Backup) Restore(file string) {
	app := strings.Split(file, "-")[0]
	filePath := fmt.Sprintf("%s/%s", backup.backupDir, file)
	log.Println("Running backup restore", BACKUP_RESTORE_CMD, app, filePath)
	out, err := exec.Command(BACKUP_RESTORE_CMD, app, filePath).CombinedOutput()
	log.Printf("Backup restore output %s", out)
	if err != nil {
		log.Printf("Backup restore failed: %v", err)
	} else {
		log.Printf("Backup restore complete")
	}

}

func (backup *Backup) Remove(file string) error {
	filePath := fmt.Sprintf("%s/%s", backup.backupDir, file)
	log.Println("Removing backup file", filePath)
	err := os.Remove(filePath)
	if err != nil {
		log.Printf("Backup remove failed: %v", err)
	} else {
		log.Printf("Backup remove completed")
	}
	return err
}
