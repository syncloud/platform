package backup

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	cp "github.com/otiai10/copy"
	df "github.com/ricochet2200/go-disk-usage/du"
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/date"
	"github.com/syncloud/platform/du"
	"github.com/syncloud/platform/snap/model"
	"go.uber.org/zap"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
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
	RestorePreStop   = "restore-pre-stop"
	RestorePostStop  = "restore-post-stop"
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
	err = cp.Copy(versionDir, tempCurrentDir, b.options())
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
	err = cp.Copy(commonDir, tempCommonDir, b.options())
	if err != nil {
		return err
	}

	err = b.snapCli.Start(app)
	if err != nil {
		return err
	}

	err = createTarGz(file, tempDir)
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

func (b *Backup) options() cp.Options {
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
		PreserveOwner: true,
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

	err = extractTarGz(file.FullName, tempDir)
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

	err = b.snapCli.RunCmdIfExists(*snap, RestorePreStop)
	if err != nil {
		return err
	}

	err = b.snapCli.Stop(file.App)
	if err != nil {
		return err
	}

	err = b.snapCli.RunCmdIfExists(*snap, RestorePostStop)
	if err != nil {
		return err
	}

	appBaseDir := fmt.Sprintf("%s/%s", b.varDir, file.App)

	currentDir := fmt.Sprintf("%s/current", appBaseDir)
	targetCurrentDir, err := filepath.EvalSymlinks(currentDir)
	if err != nil {
		return err
	}
	err = b.recreateDir(targetCurrentDir, file.App)
	if err != nil {
		return err
	}
	tempCurrentDir := fmt.Sprintf("%s/current", tempDir)
	b.logger.Info(fmt.Sprintf("copy %s to %s", tempCurrentDir, targetCurrentDir))
	err = cp.Copy(tempCurrentDir, targetCurrentDir, b.options())
	if err != nil {
		return err
	}
	err = b.chown(targetCurrentDir, file.App)
	if err != nil {
		return err
	}

	commonDir := fmt.Sprintf("%s/common", appBaseDir)
	err = b.recreateDir(commonDir, file.App)
	if err != nil {
		return err
	}
	tempCommonDir := fmt.Sprintf("%s/common", tempDir)
	b.logger.Info(fmt.Sprintf("copy %s to %s", tempCommonDir, commonDir))
	err = cp.Copy(tempCommonDir, commonDir, b.options())
	if err != nil {
		return err
	}
	err = b.chown(commonDir, file.App)
	if err != nil {
		return err
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

func (b *Backup) recreateDir(dir, app string) error {
	b.logger.Info(fmt.Sprintf("recreate dir %s", dir))
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	return os.MkdirAll(dir, 0755)
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

func (b *Backup) chown(dir, app string) error {
	u, err := user.Lookup(app)
	if err != nil {
		b.logger.Error("looking up user", zap.String("user", app), zap.Error(err))
		return err
	}
	g, err := user.LookupGroup(app)
	if err != nil {
		b.logger.Error("looking up group", zap.String("user", app), zap.Error(err))
		return err
	}
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		b.logger.Error("converting UID", zap.String("uid", u.Uid), zap.Error(err))
		return err
	}
	gid, err := strconv.Atoi(g.Gid)
	if err != nil {
		b.logger.Error("converting GID", zap.String("gid", u.Gid), zap.Error(err))
		return err
	}
	err = os.Chown(dir, uid, gid)
	if err != nil {
		b.logger.Error("changing ownership", zap.String("dir", dir), zap.Error(err))
	}
	return err

}

func createTarGz(outputFile, sourceDir string) error {
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	gw := gzip.NewWriter(f)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = "./" + rel
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			header.Uid = int(stat.Uid)
			header.Gid = int(stat.Gid)
		}
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(tw, file)
		return err
	})
}

func extractTarGz(archiveFile, destDir string) error {
	f, err := os.Open(archiveFile)
	if err != nil {
		return err
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(destDir, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
			os.Chown(target, header.Uid, header.Gid)
		case tar.TypeReg:
			dir := filepath.Dir(target)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
			file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(file, tr); err != nil {
				file.Close()
				return err
			}
			file.Close()
			os.Chown(target, header.Uid, header.Gid)
		}
	}
	return nil
}
