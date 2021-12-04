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
	Dir        = "/data/platform/backup"
	CreateCmd  = "/snap/platform/current/bin/backup.sh"
	RestoreCmd = "/snap/platform/current/bin/restore.sh"
)

func NewDefault() *Backup {
	return New(Dir)
}

func New(dir string) *Backup {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Println("unable to create backup dir", err)
		}
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
	now := time.Now().Format("2006-0102-150405")
	file := fmt.Sprintf("%s/%s-%s.tar.gz", backup.backupDir, app, now)
	log.Println("Running backup create", CreateCmd, app, file)
	out, err := exec.Command(CreateCmd, app, file).CombinedOutput()
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
	log.Println("Running backup restore", RestoreCmd, app, filePath)
	out, err := exec.Command(RestoreCmd, app, filePath).CombinedOutput()
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
