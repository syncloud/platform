package backup

import (
	"fmt"
	cp "github.com/otiai10/copy"
	df "github.com/ricochet2200/go-disk-usage/du"
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/date"
	"github.com/syncloud/platform/du"
	"github.com/syncloud/platform/snap/model"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

type SnapService interface {
	Stop(name string) error
	Start(name string) error
	RunCmdIfExists(snap model.Snap, cmd string) error
}

type SnapInfo interface {
	FindInstalled(name string) (*model.Snap, error)
}

type UserConfig interface {
	GetBackupAuto() string
	SetBackupAuto(auto string)
	GetBackupAutoDay() int
	SetBackupAutoDay(day int)
	GetBackupAutoHour() int
	SetBackupAutoHour(hour int)
}

type Backup struct {
	backupDir    string
	varDir       string
	executor     cli.Executor
	snapCli      SnapService
	snapServer   SnapInfo
	diskusage    du.DiskUsage
	userConfig   UserConfig
	timeProvider date.Provider
	logger       *zap.Logger
}

const (
	Dir              = "/data/platform/backup"
	VarDir           = "/var/snap"
	CreatePreStop    = "backup-pre-stop"
	CreatePostStop   = "backup-post-stop"
	RestorePreStart  = "restore-pre-start"
	RestorePostStart = "restore-post-start"
)

func New(dir string,
	varDir string,
	executor cli.Executor,
	diskusage du.DiskUsage,
	snapCli SnapService,
	snapServer SnapInfo,
	userConfig UserConfig,
	timeProvider date.Provider,
	logger *zap.Logger) *Backup {
	return &Backup{
		backupDir:    dir,
		varDir:       varDir,
		executor:     executor,
		diskusage:    diskusage,
		snapCli:      snapCli,
		snapServer:   snapServer,
		userConfig:   userConfig,
		timeProvider: timeProvider,
		logger:       logger,
	}
}

func (b *Backup) Start() error {
	if _, err := os.Stat(b.backupDir); os.IsNotExist(err) {
		err = os.MkdirAll(b.backupDir, os.ModePerm)
		if err != nil {
			b.logger.Info("unable to create backup dir", zap.Error(err))
			return err
		}
	}
	return nil
}

func (b *Backup) Auto() Auto {
	return Auto{
		Auto: b.userConfig.GetBackupAuto(),
		Day:  b.userConfig.GetBackupAutoDay(),
		Hour: b.userConfig.GetBackupAutoHour(),
	}
}

func (b *Backup) SetAuto(auto Auto) {
	b.userConfig.SetBackupAuto(auto.Auto)
	b.userConfig.SetBackupAutoDay(auto.Day)
	b.userConfig.SetBackupAutoHour(auto.Hour)
}

func (b *Backup) List() ([]File, error) {
	files, err := os.ReadDir(b.backupDir)
	if err != nil {
		b.logger.Error("Cannot get list of files in ", zap.String("backupDir", b.backupDir), zap.Error(err))
		return nil, err
	}
	var names []File
	for _, x := range files {
		file, err := Parse(b.backupDir, x.Name())
		if err != nil {
			b.logger.Error("Cannot parse file name", zap.String("file", x.Name()), zap.Error(err))
		} else {
			names = append(names, file)
		}
	}

	return names, nil
}

func (b *Backup) Create(app string) error {
	now := b.timeProvider.Now().Format("2006-0102-150405")
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
	tempSpaceNeeded := (appCurrentSize + appCommonSize) * 2

	b.logger.Info(fmt.Sprintf("temp space left: %d", tempSpaceLeft))
	b.logger.Info(fmt.Sprintf("temp space needed: %d", tempSpaceNeeded))

	if tempSpaceLeft < tempSpaceNeeded {
		return fmt.Errorf("not enough temp space for the backup")
	}

	snap, err := b.snapServer.FindInstalled(app)
	if err != nil {
		return err
	}
	if snap == nil {
		return fmt.Errorf("app not found: %s", app)
	}

	err = b.snapCli.RunCmdIfExists(*snap, CreatePreStop)
	if err != nil {
		return err
	}

	err = b.snapCli.Stop(app)
	if err != nil {
		return err
	}

	err = b.snapCli.RunCmdIfExists(*snap, CreatePostStop)
	if err != nil {
		return err
	}

	tempCurrentDir := fmt.Sprintf("%s/current", tempDir)
	b.logger.Info(fmt.Sprintf("temp dir %s", tempCurrentDir))
	err = os.Mkdir(tempCurrentDir, 0755)
	if err != nil {
		return err
	}
	versionDir, err := filepath.EvalSymlinks(currentDir)
	if err != nil {
		return err
	}
	b.logger.Info(fmt.Sprintf("copy %s", versionDir))
	err = cp.Copy(versionDir, tempCurrentDir, b.skipUnixSockets())
	if err != nil {
		b.logger.Error("cannot copy", zap.Error(err))
		return err
	}

	tempCommonDir := fmt.Sprintf("%s/common", tempDir)
	err = os.Mkdir(tempCommonDir, 0755)
	if err != nil {
		return err
	}

	b.logger.Info(fmt.Sprintf("copy %s", commonDir))
	err = cp.Copy(commonDir, tempCommonDir, b.skipUnixSockets())
	if err != nil {
		return err
	}

	err = b.snapCli.Start(app)
	if err != nil {
		return err
	}

	out, err := b.executor.CombinedOutput("tar", "czf", file, "-C", tempDir, ".")
	b.logger.Info(fmt.Sprintf("tar output: %s", string(out)))
	if err != nil {
		return err
	}

	b.logger.Info(fmt.Sprintf("cleanup %s", tempDir))
	err = os.RemoveAll(tempDir)
	if err != nil {
		return err
	}

	return nil
}

func (b *Backup) skipUnixSockets() cp.Options {
	return cp.Options{
		Skip: func(src string) (bool, error) {
			info, err := os.Lstat(src)
			if err != nil {
				return true, err
			}
			if info.Mode()&os.ModeSocket != 0 {
				return true, nil
			}
			return false, nil
		},
	}
}

func (b *Backup) Restore(fileName string) error {
	file, err := Parse(b.backupDir, fileName)
	if err != nil {
		return err
	}
	b.logger.Info("Running backup restore", zap.String("app", file.App), zap.String("file", file.FullName))

	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		b.logger.Info("cannot create tmp dir", zap.Error(err))
		return err
	}

	fileStat, err := os.Stat(file.FullName)
	if err != nil {
		return err
	}

	tempSpaceLeft := df.NewDiskUsage(tempDir).Available()
	tempSpaceNeeded := uint64(fileStat.Size()) * 2
	b.logger.Info(fmt.Sprintf("temp space left: %d", tempSpaceLeft))
	b.logger.Info(fmt.Sprintf("temp space needed: %d", tempSpaceNeeded))

	if tempSpaceLeft < tempSpaceNeeded {
		return fmt.Errorf("not enough temp space for the restore")
	}

	out, err := b.executor.CombinedOutput("tar", "-C", tempDir, "-xf", file.FullName)
	b.logger.Info(fmt.Sprintf("tar output: %s", string(out)))
	if err != nil {
		return err
	}

	err = b.snapCli.Stop(file.App)
	if err != nil {
		return err
	}

	appBaseDir := fmt.Sprintf("%s/%s", b.varDir, file.App)

	currentDir := fmt.Sprintf("%s/current", appBaseDir)
	_, err = b.executor.CombinedOutput("rm", "-rf", fmt.Sprintf("%s/*", currentDir))
	if err != nil {
		return err
	}
	tempCurrentDir := fmt.Sprintf("%s/current", tempDir)
	err = cp.Copy(tempCurrentDir, currentDir)
	if err != nil {
		return err
	}

	commonDir := fmt.Sprintf("%s/common", appBaseDir)
	_, err = b.executor.CombinedOutput("rm", "-rf", fmt.Sprintf("%s/*", commonDir))
	if err != nil {
		return err
	}
	tempCommonDir := fmt.Sprintf("%s/common", tempDir)
	err = cp.Copy(tempCommonDir, commonDir)
	if err != nil {
		return err
	}

	snap, err := b.snapServer.FindInstalled(file.App)
	if err != nil {
		return err
	}
	if snap == nil {
		return fmt.Errorf("app not found: %s", file.App)
	}

	err = b.snapCli.RunCmdIfExists(*snap, RestorePreStart)
	if err != nil {
		return err
	}

	err = b.snapCli.Start(file.App)
	if err != nil {
		return err
	}

	err = b.snapCli.RunCmdIfExists(*snap, RestorePostStart)
	if err != nil {
		return err
	}

	err = os.RemoveAll(tempDir)
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
