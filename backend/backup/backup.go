package backup

import (
	"fmt"
	cp "github.com/otiai10/copy"
	"github.com/ricochet2200/go-disk-usage/du"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Backup struct {
	backupDir string
	varDir    string
	logger    *zap.Logger
}

const (
	Dir        = "/data/platform/backup"
	RestoreCmd = "/snap/platform/current/bin/restore.sh"
	VarDir     = "/var/snap"
)

func New(dir string, varDir string, logger *zap.Logger) *Backup {
	return &Backup{
		backupDir: dir,
		varDir:    varDir,
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
		b.logger.Error("Cannot get list of files in ", zap.String("backupDir", b.backupDir), zap.Error(err))
		return nil, err
	}
	var names []File
	for _, x := range files {
		names = append(names, File{b.backupDir, x.Name()})
	}

	return names, nil
}

func (b *Backup) Create(app string) error {
	now := time.Now().Format("2006-0102-150405")
	file := fmt.Sprintf("%s/%s-%s.tar.gz", b.backupDir, app, now)
	b.logger.Info("Running backup create", zap.String("app", app), zap.String("file", file))

	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		panic(err)
	}
	appBaseDir := fmt.Sprintf("%s/%s", b.varDir, app)
	AppCurrentDir := fmt.Sprintf("%s/current", appBaseDir)
	AppCommonDir := fmt.Sprintf("%s/common", appBaseDir)
	appCurrentSize := du.NewDiskUsage(AppCurrentDir).Used() / 1024 / 1024
	appCommonSize := du.NewDiskUsage(AppCommonDir).Used() / 1024 / 1024

	tempSpaceLeft := du.NewDiskUsage(tempDir).Available() / 1024 / 1024
	TempSpaceNeeded := appCurrentSize + appCommonSize*2

	b.logger.Info(fmt.Sprintf("temp space left: %d", tempSpaceLeft))
	b.logger.Info(fmt.Sprintf("temp space needed: %d", TempSpaceNeeded))

	if tempSpaceLeft < TempSpaceNeeded {
		return fmt.Errorf("not enough temp space for the backup")
	}

	//snap run $APP.backup-create-pre-stop
	out, err := exec.Command("snap", "stop", app).CombinedOutput()
	b.logger.Info(fmt.Sprintf("stop output: %s", string(out)))
	if err != nil {
		return err
	}
	//snap run $APP.backup-create-post-stop

	tempCurrentDir := fmt.Sprintf("%s/current", tempDir)
	err = os.Mkdir(tempCurrentDir, 0755)
	if err != nil {
		return err
	}
	err = cp.Copy(AppCurrentDir, tempCurrentDir)
	if err != nil {
		return err
	}

	tempCommonDir := fmt.Sprintf("%s/common", tempDir)
	err = os.Mkdir(tempCommonDir, 0755)
	if err != nil {
		return err
	}
	err = cp.Copy(AppCommonDir, tempCommonDir)
	if err != nil {
		return err
	}

	out, err = exec.Command("snap", "start", app).CombinedOutput()
	b.logger.Info(fmt.Sprintf("start output: %s", string(out)))
	if err != nil {
		return err
	}
	out, err = exec.Command("tar", "czf", file, "-C", tempDir).CombinedOutput()
	b.logger.Info(fmt.Sprintf("tar output: %s", string(out)))
	if err != nil {
		return err
	}
	err = os.RemoveAll(tempDir)
	if err != nil {
		return err
	}
	return nil
}

func (b *Backup) Restore(fileName string) error {
	app := strings.Split(fileName, "-")[0]
	file := fmt.Sprintf("%s/%s", b.backupDir, fileName)
	b.logger.Info("Running backup restore", zap.String("app", app), zap.String("file", file))
	/*
		EXTRACT_DIR=$(mktemp -d)

		BACKUP_SIZE=$(stat --printf="%s" ${BACKUP_FILE})

		TEMP_SPACE_LEFT=$(df -B 1 --output=avail ${EXTRACT_DIR} | tail -1)
		TEMP_SPACE_NEEDED=$(( ${BACKUP_SIZE} * 10 ))

		echo "temp space left: ${TEMP_SPACE_LEFT}"
		echo "temp space needed: ${TEMP_SPACE_NEEDED}"

		if [[ ${TEMP_SPACE_NEEDED} -gt ${TEMP_SPACE_LEFT} ]]; then
		    echo "not enough temp space for the restore"
		    exit 1
		fi

		tar -C ${EXTRACT_DIR} -xf ${BACKUP_FILE}
		ls -la ${EXTRACT_DIR}
		APP_DIR=/var/snap/$APP

		snap stop $APP

		rm -rf ${APP_DIR}/current/*
		cp -R --preserve ${EXTRACT_DIR}/current/. ${APP_DIR}/current/
		rm -rf ${APP_DIR}/common/*
		cp -R --preserve ${EXTRACT_DIR}/common/. ${APP_DIR}/common/

		snap run $APP.backup-restore-pre-start
		snap start $APP
		snap run $APP.backup-restore-pre-stop

		rm -rf ${EXTRACT_DIR}*/
	out, err := exec.Command(RestoreCmd, app, file).CombinedOutput()
	b.logger.Info("Backup restore output", zap.String("out", string(out)))
	if err != nil {
		return err
	}
	return nil
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
