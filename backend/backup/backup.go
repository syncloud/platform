package backup

import (
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Backup struct {
	backupDir string
	logger    *zap.Logger
}

const (
	Dir        = "/data/platform/backup"
	CreateCmd  = "/snap/platform/current/bin/backup.sh"
	RestoreCmd = "/snap/platform/current/bin/restore.sh"
)

func New(dir string, logger *zap.Logger) *Backup {
	return &Backup{
		backupDir: dir,
		logger:    logger,
	}
}

func (b *Backup) Start() {
	if _, err := os.Stat(b.backupDir); os.IsNotExist(err) {
		err := os.MkdirAll(b.backupDir, os.ModePerm)
		if err != nil {
			b.logger.Info("unable to create backup dir", zap.Error(err))
		}
	}
}

func (b *Backup) List() ([]File, error) {
	files, err := ioutil.ReadDir(b.backupDir)
	if err != nil {
		b.logger.Info("Cannot get list of files in ", zap.String("backupDir", b.backupDir), zap.Error(err))
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
	b.logger.Info("Running backup create", zap.String("CreateCmd", CreateCmd), zap.String("app", app), zap.String("file", file))
	out, err := exec.Command(CreateCmd, app, file).CombinedOutput()
	b.logger.Info("Backup create output", zap.String("out", string(out)))
	if err != nil {
		b.logger.Info("Backup create failed", zap.Error(err))
	} else {
		b.logger.Info("Backup create completed")
	}
}

func (b *Backup) Restore(fileName string) {
	app := strings.Split(fileName, "-")[0]
	file := fmt.Sprintf("%s/%s", b.backupDir, fileName)
	b.logger.Info("Running backup restore", zap.String("RestoreCmd", RestoreCmd), zap.String("app", app), zap.String("file", file))
	out, err := exec.Command(RestoreCmd, app, file).CombinedOutput()
	b.logger.Info("Backup restore output", zap.String("out", string(out)))
	if err != nil {
		b.logger.Info("Backup restore failed", zap.Error(err))
	} else {
		b.logger.Info("Backup restore complete")
	}

}

func (b *Backup) Remove(fileName string) error {
	file := fmt.Sprintf("%s/%s", b.backupDir, fileName)
	b.logger.Info("Removing backup file", zap.String("file", file))
	err := os.Remove(file)
	if err != nil {
		b.logger.Info("Backup remove failed", zap.Error(err))
	} else {
		b.logger.Info("Backup remove completed")
	}
	return err
}
