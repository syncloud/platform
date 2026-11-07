package cron

import (
	"errors"
	"github.com/syncloud/platform/date"
	"github.com/syncloud/platform/storage/model"
	"go.uber.org/zap"
	"time"
)

const (
	StorageCheckHour  = 3
	StorageCheckScrub = "scrub"
	StorageCheckSmart = "smart"
	ScrubInterval     = 30 * 24 * time.Hour
	SmartInterval     = 7 * 24 * time.Hour
)

type StorageChecks interface {
	Schedule() model.CheckSchedule
	ScrubAvailable() bool
	StartScrub() error
	StartSmartTests() error
}

type StorageCheckConfig interface {
	GetStorageCheckTime(check string) time.Time
	SetStorageCheckTime(check string, time time.Time)
}

type StorageCheckJob struct {
	checks   StorageChecks
	config   StorageCheckConfig
	provider date.Provider
	logger   *zap.Logger
}

func NewStorageCheckJob(checks StorageChecks, config StorageCheckConfig, provider date.Provider, logger *zap.Logger) *StorageCheckJob {
	return &StorageCheckJob{
		checks:   checks,
		config:   config,
		provider: provider,
		logger:   logger,
	}
}

func (j *StorageCheckJob) Run() error {
	now := j.provider.Now()
	if now.Hour() != StorageCheckHour {
		return nil
	}
	schedule := j.checks.Schedule()
	var errs []error
	if schedule.Scrub && j.checks.ScrubAvailable() && j.due(StorageCheckScrub, ScrubInterval, now) {
		errs = append(errs, j.run(StorageCheckScrub, now, j.checks.StartScrub))
	}
	if schedule.Smart && j.due(StorageCheckSmart, SmartInterval, now) {
		errs = append(errs, j.run(StorageCheckSmart, now, j.checks.StartSmartTests))
	}
	return errors.Join(errs...)
}

func (j *StorageCheckJob) due(check string, interval time.Duration, now time.Time) bool {
	return now.Sub(j.config.GetStorageCheckTime(check)) >= interval-time.Hour
}

func (j *StorageCheckJob) run(check string, now time.Time, start func() error) error {
	j.logger.Info("scheduled storage check", zap.String("check", check))
	err := start()
	if err != nil {
		return err
	}
	j.config.SetStorageCheckTime(check, now)
	return nil
}
