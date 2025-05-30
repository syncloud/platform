package storage

import (
	"fmt"
	"github.com/syncloud/platform/cli"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"path"
)

type DiskStorage interface {
	Format(device string)
	BootExtend()
}

type LinkConfig interface {
	DiskLink() string
}

type Storage struct {
	config   LinkConfig
	executor cli.Executor
	logger   *zap.Logger
}

const (
	FormatCmd     = "/snap/platform/current/bin/disk_format.sh"
	BootExtendCmd = "/snap/platform/current/bin/boot_extend.sh"
)

func New(config LinkConfig, executor cli.Executor, logger *zap.Logger) *Storage {
	return &Storage{
		config:   config,
		executor: executor,
		logger:   logger,
	}
}

func (s *Storage) Format(device string) error {
	s.logger.Info("format", zap.String("cmd", FormatCmd), zap.String("device", device))
	out, err := exec.Command(FormatCmd, device).CombinedOutput()
	s.logger.Info("format", zap.String("output", string(out)))
	return err
}

func (s *Storage) BootExtend() error {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		s.logger.Info("boot extend is not supported under docker")
		return nil
	}
	s.logger.Info("boot extend", zap.String("cmd", BootExtendCmd))
	out, err := exec.Command(BootExtendCmd).CombinedOutput()
	s.logger.Info("boot extend", zap.String("output", string(out)))
	return err
}

func (s *Storage) GetAppStorageDir(app string) string {
	return path.Join(s.config.DiskLink(), app)
}

func (s *Storage) InitAppStorageOwner(app, owner string) (string, error) {
	dir, err := s.InitAppStorage(app)
	if err != nil {
		return "", err
	}
	s.logger.Info("fixing permissions", zap.String("dir", dir))
	err = s.ChownRecursive(dir, owner)
	return dir, err

}

func (s *Storage) ChownRecursive(path, user string) error {
	// (-L) is not good as we do not want to traverse symbolic links inside apps
	// (-H) is not good for the same reason (traversal of dead links)
	// (-f) is not good as it hides the real error message
	output, err := s.executor.CombinedOutput("chown", "-R", fmt.Sprintf("%s.%s", user, user), path)
	s.logger.Info("chown", zap.String("output", string(output)))
	return err
}

func (s *Storage) InitAppStorage(app string) (string, error) {
	dir := s.GetAppStorageDir(app)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return "", err
		}
	}
	return dir, nil
}
