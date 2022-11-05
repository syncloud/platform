package backup

import (
	"fmt"
	cp "github.com/otiai10/copy"
	df "github.com/ricochet2200/go-disk-usage/du"
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/du"
	"github.com/syncloud/platform/snap/model"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type SnapService interface {
	Stop(name string) error
	Start(name string) error
	Run(name string) error
}

type SnapInfo interface {
	Snap(name string) (model.Snap, error)
}

type Backup struct {
	backupDir  string
	varDir     string
	executor   cli.CommandExecutor
	snapCli    SnapService
	snapServer SnapInfo
	diskusage  du.DiskUsage
	logger     *zap.Logger
}

const (
	Dir        = "/data/platform/backup"
	RestoreCmd = "/snap/platform/current/bin/restore.sh"
	VarDir     = "/var/snap"
)

func New(dir string, varDir string, executor cli.CommandExecutor, diskusage du.DiskUsage, snapCli SnapService, snapServer SnapInfo, logger *zap.Logger) *Backup {
	return &Backup{
		backupDir:  dir,
		varDir:     varDir,
		executor:   executor,
		diskusage:  diskusage,
		snapCli:    snapCli,
		snapServer: snapServer,
		logger:     logger,
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

	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		b.logger.Info("cannot create tmp dir", zap.Error(err))
		return err
	}
	appBaseDir := fmt.Sprintf("%s/%s", b.varDir, app)
	currentDir := fmt.Sprintf("%s/current", appBaseDir)
	commonDir := fmt.Sprintf("%s/common", appBaseDir)
	appCurrentSize, err := b.diskusage.Used(currentDir)
	if err != nil {
		b.logger.Info("cannot get current dir size", zap.Error(err))
		return err
	}
	appCommonSize, err := b.diskusage.Used(commonDir)
	if err != nil {
		b.logger.Info("cannot get common dir size", zap.Error(err))
		return err
	}

	tempSpaceLeft := df.NewDiskUsage(tempDir).Available()
	TempSpaceNeeded := (appCurrentSize + appCommonSize) * 2

	b.logger.Info(fmt.Sprintf("temp space left: %d", tempSpaceLeft))
	b.logger.Info(fmt.Sprintf("temp space needed: %d", TempSpaceNeeded))

	if tempSpaceLeft < TempSpaceNeeded {
		return fmt.Errorf("not enough temp space for the backup")
	}

	snap, err := b.snapServer.Snap(app)
	if err != nil {
		return err
	}

	cmd := snap.FindCommand("backup-pre-stop")
	if cmd != nil {
		err = b.snapCli.Run(cmd.FullName())
		if err != nil {
			return err
		}
	}

	err = b.snapCli.Stop(app)
	if err != nil {
		return err
	}
	cmd = snap.FindCommand("backup-post-stop")
	if cmd != nil {
		err = b.snapCli.Run(cmd.FullName())
		if err != nil {
			return err
		}
	}

	tempCurrentDir := fmt.Sprintf("%s/current", tempDir)
	err = os.Mkdir(tempCurrentDir, 0755)
	if err != nil {
		return err
	}
	versionDir, err := filepath.EvalSymlinks(currentDir)
	if err != nil {
		return err
	}
	err = cp.Copy(versionDir, tempCurrentDir)
	if err != nil {
		b.logger.Info("cannot copy", zap.Error(err))
		return err
	}

	tempCommonDir := fmt.Sprintf("%s/common", tempDir)
	err = os.Mkdir(tempCommonDir, 0755)
	if err != nil {
		return err
	}
	err = cp.Copy(commonDir, tempCommonDir)
	if err != nil {
		return err
	}

	err = b.snapCli.Start(app)
	if err != nil {
		return err
	}

	out, err := b.executor.CommandOutput("tar", "czf", file, "-C", tempDir)
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
