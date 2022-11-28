package cron

import (
	"fmt"
	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/date"
	"github.com/syncloud/platform/snap/model"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"time"
)

const (
	AutoNo      = "no"
	AutoBackup  = "backup"
	AutoRestore = "restore"
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
	List() ([]backup.File, error)
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
	apps, err := j.snapd.InstalledUserApps()
	if err != nil {
		return err
	}
	auto := j.config.GetBackupAuto()
	if auto == AutoNo {
		j.logger.Info("auto backup is disabled")
		return nil
	}
	day := j.config.GetBackupAutoDay()
	hour := j.config.GetBackupAutoHour()
	now := j.provider.Now()
	for _, app := range apps {
		if app.Id == "syncthing" {
			j.logger.Info("syncthing app is excluded as it is used as a backup network transport")
			continue
		}
		last := j.config.GetBackupAppTime(app.Id, auto)
		if j.scheduler.ShouldRun(day, hour, now, last) {
			if auto == AutoBackup {
				j.runBackup(app, now)
			} else {
				j.runRestore(app, now)
			}
		}
	}
	return nil
}

func (j *BackupJob) runRestore(app model.SyncloudApp, now time.Time) {
	latestBackup, err := j.LatestBackup(app.Id)
	if err != nil {
		j.logger.Info("no backups to restore yet", zap.String("app", app.Id))
		return
	}
	err = j.backup.Restore(latestBackup)
	if err != nil {
		j.logger.Error("failed", zap.String("app", app.Id), zap.Error(err))
		return
	}
	j.config.SetBackupAppTime(app.Id, AutoRestore, now)
}

func (j *BackupJob) runBackup(app model.SyncloudApp, now time.Time) {
	err := j.backup.Create(app.Id)
	if err != nil {
		j.logger.Error("failed", zap.String("app", app.Id), zap.Error(err))
		return
	}
	j.config.SetBackupAppTime(app.Id, AutoBackup, now)
}

func (j *BackupJob) LatestBackup(app string) (string, error) {
	list, err := j.backup.List()
	if err != nil {
		return "", err
	}
	var appFiles []string
	for _, file := range list {
		if file.App == app {
			appFiles = append(appFiles, file.File)
		}
	}
	if len(appFiles) > 0 {
		slices.Sort(appFiles)
		return appFiles[len(appFiles)-1], nil
	}
	return "", fmt.Errorf("not found")
}
