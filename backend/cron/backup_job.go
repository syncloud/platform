package cron

import (
	"github.com/syncloud/platform/snap/model"

	"go.uber.org/zap"
)

type BackupJob struct {
	snapd  Snapd
 userConfig UserConfig
	logger *zap.Logger
}

type Snapd interface {
	InstalledUserApps() ([]model.SyncloudApp, error)
}

type UserConfig interface {
	GetBackupAuto() string
	GetBackupAutoDay() int
	GetBackupAutoHour() int
 GetBackupAppTime(string, string) *int64
	SetBackupAppTime(string, string, int64)
}

func NewBackupJob(snapd Snapd, userConfig UserConfig, logger *zap.Logger) *BackupJob {
	return &BackupJob{
		snapd:  snapd,
  userConfig: userConfig,
		logger: logger,
	}
}

func (j *BackupJob) Run() error {
	snaps, err := j.snapd.InstalledUserApps()
	if err != nil {
		return err
	}
 auto := j.userConfig.GetBackupAuto()
 if auto == "no" {
   j.logger.Info("auto backup is disabled", zap.String("auto", auto))
   return nil
 }
	for _, snap := range snaps {
		j.logger.Info("auto backup", zap.String("app", snap.Name))
	}
	return nil
}


