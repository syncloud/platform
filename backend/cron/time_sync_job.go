package cron

import (
	"github.com/syncloud/platform/cli"
	"github.com/syncloud/platform/date"
	"go.uber.org/zap"
	"time"
)

type TimeSyncJob struct {
	executor     cli.Executor
	lastRun      time.Time
	dateProvider date.Provider
	logger       *zap.Logger
}

func NewTimeSyncJob(executor cli.Executor, dateProvider date.Provider, logger *zap.Logger) *TimeSyncJob {
	return &TimeSyncJob{
		executor:     executor,
		dateProvider: dateProvider,
		logger:       logger,
	}
}

func (t *TimeSyncJob) Run() error {
	now := t.dateProvider.Now()
	if t.lastRun.Add(time.Hour * 24).Before(now) {
		output, err := t.executor.CombinedOutput("service", "ntp", "stop")
		t.logger.Info(string(output))
		if err != nil {
			return err
		}
		output, err = t.executor.CombinedOutput("ntpd", "-gq")
		t.logger.Info(string(output))
		if err != nil {
			return err
		}
		output, err = t.executor.CombinedOutput("service", "ntp", "start")
		t.logger.Info(string(output))
		if err != nil {
			return err
		}
		t.lastRun = now
	}
	return nil
}
