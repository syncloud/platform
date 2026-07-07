package cron

import (
	"strings"
	"time"

	"github.com/syncloud/platform/date"
	"github.com/syncloud/platform/job"
	"github.com/syncloud/platform/snap/model"
	"go.uber.org/zap"
)

const (
	SnapdUpgradeEveryDay = 0
	SnapdUpgradeHour     = 11
)

type SnapdInstallerInfo interface {
	Installer() (*model.InstallerInfo, error)
}

type SnapdUpgrader interface {
	Upgrade() error
}

type JobMaster interface {
	Offer(name string, j job.Job) error
}

type SnapdUpgradeJob struct {
	installerInfo SnapdInstallerInfo
	upgrader      SnapdUpgrader
	jobMaster     JobMaster
	scheduler     Scheduler
	provider      date.Provider
	lastRun       time.Time
	logger        *zap.Logger
}

func NewSnapdUpgradeJob(installerInfo SnapdInstallerInfo, upgrader SnapdUpgrader, jobMaster JobMaster, scheduler Scheduler, provider date.Provider, logger *zap.Logger) *SnapdUpgradeJob {
	return &SnapdUpgradeJob{
		installerInfo: installerInfo,
		upgrader:      upgrader,
		jobMaster:     jobMaster,
		scheduler:     scheduler,
		provider:      provider,
		logger:        logger,
	}
}

func (j *SnapdUpgradeJob) Run() error {
	now := j.provider.Now()
	if !j.scheduler.ShouldRun(SnapdUpgradeEveryDay, SnapdUpgradeHour, now, j.lastRun) {
		return nil
	}
	info, err := j.installerInfo.Installer()
	if err != nil {
		return err
	}
	if strings.TrimSpace(info.StoreVersion) == strings.TrimSpace(info.InstalledVersion) {
		j.logger.Info("snapd is up to date", zap.String("version", strings.TrimSpace(info.InstalledVersion)))
		j.lastRun = now
		return nil
	}
	j.logger.Info("upgrading snapd",
		zap.String("from", strings.TrimSpace(info.InstalledVersion)),
		zap.String("to", strings.TrimSpace(info.StoreVersion)))
	err = j.jobMaster.Offer("installer.upgrade", func() error { return j.upgrader.Upgrade() })
	if err != nil {
		return err
	}
	j.lastRun = now
	return nil
}
