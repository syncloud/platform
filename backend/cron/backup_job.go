package cron

import (
	"github.com/syncloud/platform/date"
	"github.com/syncloud/platform/snap/model"
	"go.uber.org/zap"
	"time"
)

type BackupJob struct {
	snapd     Snapd
	config    UserConfig
	backup    Backup
	provider  date.Provider
	scheduler Scheduler
	logger    *zap.Logger
}

type Snapd interface {
	InstalledUserApps() ([]model.SyncloudApp, error)
}

type UserConfig interface {
	GetBackupAuto() string
	GetBackupAutoDay() int
	GetBackupAutoHour() int
	GetBackupAppTime(string, string) time.Time
	SetBackupAppTime(string, string, time.Time)
}

type Backup interface {
	Create(app string) error
	Restore(fileName string) error
}

type Scheduler interface {
	ShouldRun(day int, hour int, now time.Time, last time.Time) bool
}

func NewBackupJob(snapd Snapd, config UserConfig, backup Backup, provider date.Provider, scheduler Scheduler, logger *zap.Logger) *BackupJob {
	return &BackupJob{
		snapd:     snapd,
		config:    config,
		backup:    backup,
		provider:  provider,
		scheduler: scheduler,
		logger:    logger,
	}
}

func (j *BackupJob) Run() error {
	snaps, err := j.snapd.InstalledUserApps()
	if err != nil {
		return err
	}
	auto := j.config.GetBackupAuto()
	if auto == "no" {
		j.logger.Info("auto backup is disabled", zap.String("auto", auto))
		return nil
	}
	day := j.config.GetBackupAutoDay()
	hour := j.config.GetBackupAutoHour()
	now := j.provider.Now()
	for _, snap := range snaps {
		last := j.config.GetBackupAppTime(snap.Name, auto)
		if j.scheduler.ShouldRun(day, hour, now, last) {
			if auto == "backup" {
				err = j.backup.Create(snap.Name)
				if err != nil {
					j.logger.Error("failed", zap.String("app", snap.Name), zap.Error(err))
				}
			} else {
				err = j.backup.Restore("file")
				if err != nil {
					j.logger.Error("failed", zap.String("app", snap.Name), zap.Error(err))
				}
			}
		}
	}
	return nil
}
