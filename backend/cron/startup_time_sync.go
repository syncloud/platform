package cron

import (
	"time"

	"go.uber.org/zap"
)

type StartupTimeSync struct {
	job     Job
	timeout time.Duration
	delay   time.Duration
	logger  *zap.Logger
}

func NewStartupTimeSync(job Job, timeout time.Duration, delay time.Duration, logger *zap.Logger) *StartupTimeSync {
	return &StartupTimeSync{job: job, timeout: timeout, delay: delay, logger: logger}
}

func (s *StartupTimeSync) Start() error {
	start := time.Now()
	for {
		err := s.job.Run()
		if err == nil {
			return nil
		}
		if time.Since(start)+s.delay >= s.timeout {
			s.logger.Error("time is not synchronized at startup, cron will retry", zap.Error(err))
			return nil
		}
		s.logger.Info("time sync at startup failed, retrying", zap.Error(err))
		time.Sleep(s.delay)
	}
}
