package storage

import (
	"github.com/syncloud/platform/storage/model"
	"go.uber.org/zap"
)

const LowFreeKB = 2 * 1024 * 1024

type FileSystem interface {
	Stat(path string) (totalKB uint64, freeKB uint64, device uint64, err error)
}

type DiskSpaceConfig interface {
	DiskLink() string
}

type DiskSpace struct {
	fileSystem FileSystem
	config     DiskSpaceConfig
	logger     *zap.Logger
}

func NewDiskSpace(fileSystem FileSystem, config DiskSpaceConfig, logger *zap.Logger) *DiskSpace {
	return &DiskSpace{
		fileSystem: fileSystem,
		config:     config,
		logger:     logger,
	}
}

func (d *DiskSpace) Status() model.DiskSpace {
	status := model.DiskSpace{Mounts: []model.DiskSpaceMount{}}
	devices := make(map[uint64]bool)
	for _, path := range []string{"/", d.config.DiskLink()} {
		total, free, device, err := d.fileSystem.Stat(path)
		if err != nil {
			d.logger.Warn("cannot check free space", zap.String("path", path), zap.Error(err))
			continue
		}
		if devices[device] {
			continue
		}
		devices[device] = true
		low := free < LowFreeKB
		status.Mounts = append(status.Mounts, model.DiskSpaceMount{
			Path:    path,
			TotalKB: total,
			FreeKB:  free,
			Low:     low,
		})
		if low {
			status.Low = true
		}
	}
	return status
}
