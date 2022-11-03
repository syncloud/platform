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
	/*
		tempDir, err := ioutil.TempDir("", "test")
		if err != nil {
			panic(err)
		}
		appBaseDir := fmt.Sprintf("/var/snap/%s", app)
		APP_CURRENT_DIR=${APP_BASE_DIR}/current
		APP_COMMON_DIR=${APP_BASE_DIR}/common



			APP_CURRENT_SIZE=$(du -s ${APP_CURRENT_DIR} | cut -f1)
			APP_COMMON_SIZE=$(du -s ${APP_COMMON_DIR} | cut -f1)

			TEMP_SPACE_LEFT=$(df --output=avail ${TEMP_DIR} | tail -1)
			TEMP_SPACE_NEEDED=$(( (${APP_CURRENT_SIZE} + ${APP_COMMON_SIZE}) * 2 ))

			echo "temp space left: ${TEMP_SPACE_LEFT}"
			echo "temp space needed: ${TEMP_SPACE_NEEDED}"

			if [[ ${TEMP_SPACE_NEEDED} -gt ${TEMP_SPACE_LEFT} ]]; then
			    echo "not enaugh temp space for the backup"
			    exit 1
			fi

			snap run $APP.backup-create-pre-stop
			snap stop $APP
			snap run $APP.backup-create-post-stop

			mkdir ${TEMP_DIR}/current
			cp -R --preserve ${APP_CURRENT_DIR}/. ${TEMP_DIR}/current

			mkdir ${TEMP_DIR}/common
			cp -R --preserve ${APP_COMMON_DIR}/. ${TEMP_DIR}/common

			snap start $APP
			tar czf ${BACKUP_FILE} -C ${TEMP_DIR} .
			rm -rf ${TEMP_DIR}
	*/
	out, err := exec.Command(CreateCmd, app, file).CombinedOutput()
	b.logger.Info("Backup create output", zap.String("out", string(out)))
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
