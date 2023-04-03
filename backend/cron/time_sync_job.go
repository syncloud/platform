package cron

import (
	"github.com/syncloud/platform/cli"
	"go.uber.org/zap"
)

type TimeSyncJob struct {
	executor cli.Executor
	ran      bool
	logger   *zap.Logger
}

func NewTimeSyncJob(executor cli.Executor, logger *zap.Logger) *TimeSyncJob {
	return &TimeSyncJob{
		executor: executor,
		logger:   logger,
	}
}

func (t *TimeSyncJob) Run() error {
	if !t.ran {
		output, err := t.CombinedOutput("service", "ntp", "stop")
		t.logger.Info(output)
		if err != nil {
			return err
		}
		output, err = t.CombinedOutput("ntpd", "-gq")
		t.logger.Info(output)
		if err != nil {
			return err
		}
		output, err = t.CombinedOutput("service", "ntp", "start")
		t.logger.Info(output)
		if err != nil {
			return err
		}
		t.ran = true
	}
	return nil
}
}
