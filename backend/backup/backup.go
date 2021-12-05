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

func New(dir string) *Backup {
	return &Backup{
		backupDir: dir,
	}
}

func (b *Backup) Start() {
	if _, err := os.Stat(b.backupDir); os.IsNotExist(err) {
		err := os.MkdirAll(b.backupDir, os.ModePerm)
		if err != nil {
			log.Println("unable to create backup dir", err)
		}
	}
}

func (b *Backup) List() ([]File, error) {
	files, err := ioutil.ReadDir(b.backupDir)
	if err != nil {
		log.Println("Cannot get list of files in ", b.backupDir, err)
		return nil, err
	}
	var names []File
	for _, x := range files {
		names = append(names, File{b.backupDir, x.Name()})
	}

	return names, nil
}

func (b *Backup) Create(app string) {
	now := time.Now().Format("2006-0102-150405")
	file := fmt.Sprintf("%s/%s-%s.tar.gz", b.backupDir, app, now)
	log.Println("Running backup create", CreateCmd, app, file)
	out, err := exec.Command(CreateCmd, app, file).CombinedOutput()
	log.Printf("Backup create output %s", out)
	if err != nil {
		log.Printf("Backup create failed: %v", err)
	} else {
		log.Printf("Backup create completed")
	}
}

func (b *Backup) Restore(file string) {
	app := strings.Split(file, "-")[0]
	filePath := fmt.Sprintf("%s/%s", b.backupDir, file)
	log.Println("Running backup restore", RestoreCmd, app, filePath)
	out, err := exec.Command(RestoreCmd, app, filePath).CombinedOutput()
	log.Printf("Backup restore output %s", out)
	if err != nil {
		log.Printf("Backup restore failed: %v", err)
	} else {
		log.Printf("Backup restore complete")
	}

}

func (b *Backup) Remove(file string) error {
	filePath := fmt.Sprintf("%s/%s", b.backupDir, file)
	log.Println("Removing backup file", filePath)
	err := os.Remove(filePath)
	if err != nil {
		log.Printf("Backup remove failed: %v", err)
	} else {
		log.Printf("Backup remove completed")
	}
	return err
}
