package storage

import (
	"errors"
	"go.uber.org/zap"
	"os"
)

type Linker struct {
	logger *zap.Logger
}

func NewLinker(logger *zap.Logger) *Linker {
	return &Linker{
		logger: logger,
	}
}
func (d *Linker) RelinkDisk(link string, target string) error {
	d.logger.Info("relink disk")
	err := os.Chmod(target, 0o755)
	if err != nil {
		return err
	}

	fi, err := os.Lstat(link)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			d.logger.Error("stat", zap.Error(err))
			return err
		}
	} else {
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			err = os.Remove(link)
			if err != nil {
				return err
			}
		}
	}

	err = os.Symlink(target, link)
	return err
}
